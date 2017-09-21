package handler

import (
  //"fmt"
  "log"
  "bytes"
	"net/http"
  "github.com/gorilla/rpc/json"
  "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
  //. "github.com/tmpfs/pageloop/model"
  //. "github.com/tmpfs/pageloop/util"
)

var(
  codec *json.Codec = json.NewCodec()
  connections []*WebsocketConnection
  upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024}
)

type WebsocketConnection struct {
  Adapter *CommandAdapter
  Conn *websocket.Conn
}

/*
type RpcRequest struct {
  Id uint `json:"id"`
  Method string `json:"method"`
  Params []*ActionParameters `json:"params"`
  Arguments []interface{} `json:"args"`
}

type RpcResponse struct {
  Id uint `json:"id"`
  Status int `json:"status"`
  Error *StatusError `json:"error,omitempty"`
  Result interface{} `json:"result,omitempty"`
}
*/

// Implements http.ResponseWriter for JSON-RPC responses
type WebsocketWriter struct {
  MessageType int
  Socket *WebsocketConnection
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
      writer := &WebsocketWriter{Socket: w, MessageType: messageType}
      r := bytes.NewBuffer(p)
      if fake, err := http.NewRequest(http.MethodPost, "/ws/", r); err != nil {
        log.Println(err)
      } else {
        req := codec.NewRequest(fake)

        // TODO: get response writer

        if method, err := req.Method(); err != nil {
          log.Println(err)
          req.WriteResponse(writer, nil, err)
        } else {
          println("method: " + method)
          println(string(p))

          // TODO: call function and get reply
          var reply interface{}
          req.WriteResponse(writer, reply, err)
        }

      }
    }
  }
}

// Handles requests for application data.
type WebsocketHandler struct {
  Adapter *CommandAdapter
}

// Configure the service. Adds a handler for the websocket URL to
// the passed servemux.
func WebsocketService(mux *http.ServeMux, adapter *CommandAdapter) http.Handler {
  handler := WebsocketHandler{Adapter: adapter}
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

  ws := &WebsocketConnection{Conn: conn, Adapter: h.Adapter}
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
