package handler

import (
  "fmt"
  "log"
	"net/http"
  "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
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

    // NOTE: client socket keepalive requests send the
    // NOTE: empty object so we ignore requests with no
    // NOTE: method to avoid sending error responses to
    // NOTE: client keepalive messages

    if method != "" {
      println(method)
      // Create an empty action to hold the arguments
      var act *Action = &Action{}
      if len(req.Params) >= 1 {
        req.Params[0].Assign(act)
      }

      // Make sure the service method exists
      if _, err := w.Adapter.FindService(method, act); err != nil {
        w.WriteError(req, err)
      } else {
        // Handle additional arguments
        if len(req.Arguments) > 0 {
          // TODO: use proper RPC arguments interface
          if method == "Container.CreateApp" {
            //println("test for create app arguments")
            //fmt.Printf("%#v\n", req.Arguments)

            if input, ok := req.Arguments[0].(map[string]interface{}); ok {
              if result, err := utils.ValidateInterface(SchemaAppNew, input); err != nil {
                w.WriteError(req, CommandError(http.StatusInternalServerError, err.Error()))
                continue
              } else {
                if !result.Valid() {
                  w.WriteError(req, CommandError(http.StatusBadRequest,result.Errors()[0].String()))
                  continue
                } else {
                  app := &Application{Name: input["name"].(string), Description: input["description"].(string)}
                  act.Push(app)
                }
              }
            }
          } else if method == "Application.DeleteFiles" {
            if input, ok := req.Arguments[0].([]interface{}); ok {
              list := UrlList{}
              for _, url := range input {
                list = append(list, url.(string))
              }
              act.Push(list)
            }
          } else if method == "File.Create" {
            // Create empty document
            var content []byte
            act.Push(content)
          } else if method == "File.CreateTemplate" {
            // fmt.Printf("%#v\n", req.Arguments)
            if input, ok := req.Arguments[0].(map[string]interface{}); ok {
              c := input["container"].(string)
              a := input["application"].(string)
              f := input["file"].(string)
              ref := &ApplicationTemplate{Container: c, Application: a, File: f}
              act.Push(ref)
            }
          }
        }

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
