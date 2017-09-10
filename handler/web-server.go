package handler

import (
  "strings"
  "net/http"
  . "github.com/tmpfs/pageloop/core"
)

// Main HTTP server handler.
type ServerHandler struct {
  Mux *http.ServeMux
  // Reference to the mountpoint manager
  MountpointManager *MountpointManager
}

// The default server handler, defers to a multiplexer.
func (h ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var handler http.Handler
	var path string = req.URL.Path

  res.Header().Set("Access-Control-Allow-Origin", "*")

  var system []string
  system = append(system, API_URL, RPC_URL)
	// Look for system services first
	for _, u := range system {
		if strings.HasPrefix(path, u) {
			handler, _ = h.Mux.Handler(req)
			handler.ServeHTTP(res, req)
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
	handler.ServeHTTP(res, req)
}
