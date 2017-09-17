package handler

import (
  "strings"
  "strconv"
  "net/http"
  // "github.com/gorilla/websocket"
  . "github.com/tmpfs/pageloop/core"
)

// Main HTTP server handler.
type ServerHandler struct {
  Mux *http.ServeMux
  // Reference to the mountpoint manager
  MountpointManager *MountpointManager
}

type ResponseWriterProxy struct {
  Response http.ResponseWriter
}

func (w *ResponseWriterProxy) Header() http.Header {
  return w.Response.Header()
}

func (w *ResponseWriterProxy) WriteHeader(status int) {
  w.Response.WriteHeader(status)
  Stats.Http.Add(strconv.Itoa(status), 1)
  Stats.Http.Add("responses", 1)
}

func (w *ResponseWriterProxy) Write(data []byte) (int, error) {
  if written, err := w.Response.Write(data); err != nil {
    return 0, err
  } else {
    Stats.Http.Add("body-out", int64(written))
    return written, nil
  }
}

// The default server handler, defers to a multiplexer.
func (h ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var handler http.Handler
	var path string = req.URL.Path

  // println("got request: " + path)

  /*
  if websocket.IsWebSocketUpgrade(req) {
    println("do upgrade")
    return
  }
  */

  res.Header().Set("Access-Control-Allow-Origin", "*")

  Stats.Http.Add("requests", 1)

  if req.ContentLength > -1 {
    Stats.Http.Add("body-in", req.ContentLength)
  }

  proxy := &ResponseWriterProxy{Response: res}

  var system []string
  system = append(system, API_URL, RPC_URL, WEBSOCKET_URL)
	// Look for system services first
	for _, u := range system {
		if strings.HasPrefix(path, u) {
			handler, _ = h.Mux.Handler(req)
			handler.ServeHTTP(proxy, req)
			return
		}
	}

	// Check for application mountpoints.
	//
	// Serve the highest score which is the longest
	// matching URL path.
	var score int
	for k, v := range h.MountpointManager.MountpointMap {
		if strings.HasPrefix(path, k) {
			if handler != nil && len(k) < score {
				continue
			}
			handler = v
			score = len(k)
		}
	}

	if handler == nil {
		handler = http.NotFoundHandler()
	}
	handler.ServeHTTP(proxy, req)
}
