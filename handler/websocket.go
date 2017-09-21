package handler

import (
  //"fmt"
  "log"
  "bytes"
	"net/http"
  "github.com/gorilla/rpc"
  "github.com/gorilla/rpc/json"
  "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/util"
)

var(
  codec *json.Codec = json.NewCodec()
  connections []*WebsocketConnection
  upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024}
)

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
  println("Writing response: " + string(p))
  if err := writer.Socket.Conn.WriteMessage(writer.MessageType, p); err != nil {
    return 0, err
  }
  return len(p), nil
}

/*
func (writer *WebsocketWriter) WriteJson(args interface{}) error {
  return writer.Socket.Conn.WriteJSON(args)
}
*/

func (writer *WebsocketWriter) WriteError(err *StatusError) error {
  return writer.Request.WriteResponse(writer, nil, err)
}

/*
func (w *WebsocketConnection) WriteError(req *RpcRequest, err *StatusError) *StatusError {
  res:= &RpcResponse{Id: req.Id, Status: err.Status, Error: err}
  return w.WriteResponse(res)
}

func (w *WebsocketConnection) WriteResponse(res *RpcResponse) *StatusError {
  if err := w.Conn.WriteJSON(res); err != nil {
    return CommandError(http.StatusInternalServerError, err.Error())
  }
  return nil
}
*/

func (w *WebsocketConnection) ReadRequest() {
  for {
    // Read in the message
    messageType, p, err := w.Conn.ReadMessage()
    if err != nil {
      return
    }

    // Treat text messages as JSON-RPC
    if messageType == websocket.TextMessage {
      println("request bytes: " + string(p))
      r := bytes.NewBuffer(p)
      if fake, err := http.NewRequest(http.MethodPost, "/ws/", r); err != nil {
        log.Println(err)
      } else {
        req := codec.NewRequest(fake)
        writer := &WebsocketWriter{Socket: w, MessageType: messageType, Request: req}
        if method, err := req.Method(); err != nil {
          log.Println(err)
          req.WriteResponse(writer, nil, err)
        } else {
          println("method: " + method)
          println(string(p))

          hasServiceMethod := w.Handler.Services.HasMethod(method)
          // Check if the service method is available
          if !hasServiceMethod {
            writer.WriteError(CommandError(http.StatusNotFound, "Service %s does not exist", method))
            return
          }

          // Get a service method call request
          if rpcreq, err := w.Handler.Services.Request(method, 0); err != nil {
            writer.WriteError(CommandError(http.StatusInternalServerError, err.Error()))
            return
          } else {
            // TODO: read params into correct type

            if reply, err := w.Handler.Services.Call(rpcreq); err != nil {
              writer.WriteError(CommandError(http.StatusInternalServerError, err.Error()))
              return
            } else {
              // Reply with error when available
              if reply.Error != nil {
                // Send status error if we can
                if err, ok := reply.Error.(*StatusError); ok {
                  writer.WriteError(err)
                  return
                // Otherwise handle as plain error
                } else {
                  // TODO: wrap error
                  req.WriteResponse(writer, reply, err)
                  return
                }
              // Success send the response to the client
              } else {
                println("Sending reply!!!")
                req.WriteResponse(writer, reply, nil)
              }
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
