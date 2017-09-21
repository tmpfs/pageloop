package handler

import (
  //"fmt"
  "log"
	"net/http"
  "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
  //. "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

var(
  connections []*WebsocketConnection
  upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024}
)

type WebsocketConnection struct {
  Adapter *CommandAdapter
  Conn *websocket.Conn
}

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

func (w *WebsocketConnection) ReadRequest() {
  for {
    req := &RpcRequest{}
    err := w.Conn.ReadJSON(req)
    if err != nil {
      log.Println(err)
      return
    }

    method := req.Method
    if method != "" {
      println("got websocket request")
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
