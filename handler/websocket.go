package handler

import (
  // "fmt"
  "log"
	"net/http"
  "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
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
  Params []interface{} `json:"params"`
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
    // messageType, p, err := w.Conn.ReadMessage()
    req := &RpcRequest{}
    err := w.Conn.ReadJSON(req)
    if err != nil {
      log.Println(err)
      return
    }

    // fmt.Printf("%#v\n", request)

    method := req.Method

    // Could not find service method
    if method == "" {
      w.WriteError(req, CommandError(http.StatusBadRequest, "Service method name required"))
    } else {
      // Create an empty action to hold the arguments
      var act *Action = &Action{}
      if _, err := w.Adapter.FindService(method, act); err != nil {
        w.WriteError(req, err)
      } else {
        // Handle execution errors
        if result, err := w.Adapter.Execute(act); err != nil {
          w.WriteError(req, err)
        // Write out the command invocation result
        } else {
          res := &RpcResponse{Id: req.Id, Result: result.Data, Status: result.Status}
          w.WriteResponse(res)
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

  // fmt.Printf("%#v\n", conn)

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

  // Start reading from socket in goroutine
  go ws.ReadRequest()
}
