package handler

import (
  "fmt"
  "log"
	"net/http"
  "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
)

var(
  upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024}
)

// Handles requests for application data.
type WebsocketHandler struct {
  Adapter *CommandAdapter
}

// Configure the service. Adds a rest handler for the API URL to
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
  }

  fmt.Printf("%#v\n", conn)

  //for {
    //messageType, _, err := conn.ReadMessage()
    //if err != nil {
      //return
    //}

    //println(string(messageType))

    //if err := conn.WriteMessage(messageType, p); err != nil {
      //return err
    //}
  //}
}

