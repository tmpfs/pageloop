package handler

import (
  // "fmt"
  "log"
	"net/http"
  "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
)

var(
  connections []*WebsocketConnection
  upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024}
)

type WebsocketConnection struct {
  Conn *websocket.Conn
}

func (w *WebsocketConnection) Read() {
  for {
    messageType, p, err := w.Conn.ReadMessage()
    if err != nil {
      log.Println(err)
      return
    }

    println(string(messageType))
    println(string(p))

    // TODO: handle write errors
    /*
    if err := conn.WriteMessage(messageType, p); err != nil {
      return
    }
    */
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

  ws := &WebsocketConnection{Conn: conn}
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
  go ws.Read()

  /*
  for {
    messageType, p, err := conn.ReadMessage()
    if err != nil {
      return
    }

    // println(string(p))

    // TODO: handle write errors
    if err := conn.WriteMessage(messageType, p); err != nil {
      return
    }
  }
  */
}
