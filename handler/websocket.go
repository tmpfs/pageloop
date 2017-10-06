package handler

import (
  //"fmt"
  "log"
  "bytes"
	"net/http"
  "github.com/gorilla/rpc/v2"
  "github.com/gorilla/rpc/v2/json"
  "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/service"
  . "github.com/tmpfs/pageloop/util"
)

var(
  ping = []byte("{}")
  codec *json.Codec = json.NewCodec()
  connections []*WebsocketConnection
  upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024}
)

// Wrapped result object for JSON-RPC messages so the client
// can test on status code
type RpcWebsocketReply struct {
  Document interface{} `json:"document"`
  Status int `json:"status"`
}

type WebsocketConnection struct {
  Handler WebsocketHandler
  Conn *websocket.Conn
}

// Implements http.ResponseWriter for JSON-RPC responses
type WebsocketWriter struct {
  MessageType int
  Socket *WebsocketConnection
  Request rpc.CodecRequest
}

func (writer *WebsocketWriter) WriteHeader(int) {}

func (writer *WebsocketWriter) Header() http.Header {
  return http.Header{}
}

func (writer *WebsocketWriter) Write(p []byte) (int, error) {
  if err := writer.Socket.Conn.WriteMessage(writer.MessageType, p); err != nil {
    return 0, err
  }
  return len(p), nil
}

// Write an error response to the result object so that clients
// can inspect the status code for the error.
func (writer *WebsocketWriter) WriteError(err *StatusError) {
  m := make(map[string]*StatusError)
  m["error"] = err
  // TODO: use rpc v2 style error handling?
  writer.Request.WriteResponse(writer, m)
}

// Determine the type of argument for the given service method and
// call ReadRequest() on the CodecRequest to parse the input params
// into the correct type.
func (w *WebsocketConnection) RequestArgv(req rpc.CodecRequest, method string) (argv interface{}, err error) {
  argv = &VoidArgs{}
  switch(method) {
    case "Archive.Export":
      argv = &ArchiveRequest{}
    case "Container.Read":
      argv = &ContainerRequest{}
    case "Container.CreateApp":
      fallthrough
    case "Application.ReadFiles":
      fallthrough
    case "Application.ReadPages":
      fallthrough
    case "Application.DeleteFiles":
      fallthrough
    case "Application.Delete":
      fallthrough
    case "Application.RunTask":
      fallthrough
    case "Application.Read":
      argv = &ApplicationRequest{}
    case "File.Move":
      argv = &FileMoveRequest{}
    case "File.ReadPage":
      fallthrough
    case "File.ReadSource":
      fallthrough
    case "File.ReadSourceRaw":
      fallthrough
    case "File.Read":
      argv = &FileReferenceRequest{}
    case "File.Create":
      fallthrough
    case "File.CreateTemplate":
      fallthrough
    case "File.Save":
      argv = &FileRequest{}
    case "Service.Read":
      argv = &ServiceRequest{}
    case "Service.ReadMethodCalls":
      fallthrough
    case "Service.ReadMethod":
      argv = &ServiceMethodRequest{}
    case "Job.Delete":
      fallthrough
    case "Job.Read":
      argv = &JobRequest{}

  }
  if argv != nil {
    // Read in the request params to the type we expect
    err = req.ReadRequest(argv)
  }
  return
}

func (w *WebsocketConnection) ReadRequest() {
  for {
    // Read in the message
    messageType, p, err := w.Conn.ReadMessage()
    if err != nil {
      log.Println(err.Error())
      // Cannot re-read now, we need to stop reading
      // close the socket connection on unrecoverable error
      w.Conn.Close()
      return
    }

    // Treat text messages as JSON-RPC
    if messageType == websocket.TextMessage {
      // Drop ping requests
      if bytes.Equal(p, ping) {
        continue
      }

      // TODO: restore create app validation!

      r := bytes.NewBuffer(p)
      if fake, err := http.NewRequest(http.MethodPost, "/ws/", r); err != nil {
        log.Println(err.Error())
        continue
      } else {
        req := codec.NewRequest(fake)
        writer := &WebsocketWriter{Socket: w, MessageType: messageType, Request: req}
        if method, err := req.Method(); err != nil {
          //req.WriteResponse(writer, nil, err)
          writer.WriteError(CommandError(http.StatusInternalServerError, err.Error()))
          continue
        } else {
          hasServiceMethod := w.Handler.Services.HasMethod(method)
          // Check if the service method is available
          if !hasServiceMethod {
            writer.WriteError(
              CommandError(http.StatusNotFound, "Service %s does not exist", method))
            continue
          }

          // Get a service method call request
          if rpcreq, err := w.Handler.Services.Request(method, 0); err != nil {
            writer.WriteError(
              CommandError(http.StatusInternalServerError, err.Error()))
            continue
          } else {
            if argv, err := w.RequestArgv(req, method); err != nil {
              // If we had an error while reading the request
              // if is likely a JSON unmarshal error so treat as
              // a bad request
              writer.WriteError(
                CommandError(http.StatusBadRequest, err.Error()))
              continue
            } else {
              if argv != nil {
                rpcreq.Argv(argv)
              }
            }

            // Call the service function
            Stats.Rpc.Add("calls", 1)
            if reply, err := w.Handler.Services.Call(rpcreq); err != nil {
              Stats.Rpc.Add("errors", 1)
              if ex, ok := err.(*StatusError); ok {
                writer.WriteError(ex)
              } else {
                writer.WriteError(CommandError(http.StatusInternalServerError, err.Error()))
              }
              continue
            } else {
              // NOTE: we don't need to test reply.Error as the error is always returned

              // Success send the response to the client
              status := http.StatusOK
              replyData := reply.Reply

              if result, ok := replyData.(*ServiceReply); ok {
                replyData = result.Reply
                if result.Status != 0 {
                  status = result.Status
                }
              }

              if method == "Container.CreateApp" {
                // Mount the application, needs to be done here due to some funky
                // package cyclic references
                if app, ok := replyData.(*Application); ok {
                  MountApplication(w.Handler.Mountpoints.MountpointMap, w.Handler.Host, app)
                }
              }

              // Wrap the result object so we can extract
              // status code client side
              replyData = &RpcWebsocketReply{Document: replyData, Status: status}

              req.WriteResponse(writer, replyData)
            }
          }
        }
      }
    }
  }
}

// Handles requests for application data.
type WebsocketHandler struct {
  Services *ServiceMap
  Host *Host
  Mountpoints *MountpointManager
}

// Configure the service. Adds a handler for the websocket URL to
// the passed servemux.
func WebsocketService(mux *http.ServeMux, services *ServiceMap, host *Host, mountpoints *MountpointManager) http.Handler {
  handler := WebsocketHandler{Services: services, Host: host, Mountpoints: mountpoints}
  mux.Handle(WEBSOCKET_URL, http.StripPrefix(WEBSOCKET_URL, handler))
	return handler
}

// Handle websocket endpoint requests.
func (h WebsocketHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  conn, err := upgrader.Upgrade(res, req, nil)
  if err != nil {
    log.Println(err)
    return
  }

  ws := &WebsocketConnection{Conn: conn, Handler: h}
  connections = append(connections, ws)
  Stats.Websocket.Add("connections", 1)

  conn.SetCloseHandler(func(code int, text string) error {
    for i, ws := range connections {
      if ws.Conn == conn {
        before := connections[0:i]
        after := connections[i+1:]
        connections = append(before, after...)
        Stats.Websocket.Add("connections", -1)
      }
    }
    return nil
  })

  // Start reading messages from socket
  go ws.ReadRequest()
}
