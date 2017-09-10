package pageloop

import (
  //"os"
  //"fmt"
  //"log"
  "strings"
  "net/http"
  //"path/filepath"
  . "github.com/tmpfs/pageloop/model"
)

// Serves application source files from memory.
type ApplicationSourceHandler struct {
	App *Application
	Raw bool
}

// Tests the request path and attempts to find a corresponding source file
// in the application files.
func (h ApplicationSourceHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	index := "index.html"
	urls := h.App.Urls
	path := "/" + req.URL.Path
	clean := strings.TrimSuffix(path, "/")
  trailing := clean + "/"
	indexPage := clean + "/" + index

	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var file *File

	// Exact match
	if urls[path] != nil {
		file = urls[path]
  } else if urls[trailing] != nil && !strings.HasSuffix(path, "/") {
    redirect := http.RedirectHandler(req.URL.Path + "/", http.StatusMovedPermanently)
    redirect.ServeHTTP(res, req)
    return
	// Normalized without a trailing slash
	} else if urls[clean] != nil {
		file = urls[clean]
	// Check for index page
	} else if urls[indexPage] != nil {
		file = urls[indexPage]
	}

	// TODO: write cache busting headers
	// TODO: handle directory requests (no data)
	if file != nil && !file.Info().IsDir() {
		output := file.Source(h.Raw)
    send(res, req, file, output)
		return
  // Handle directory listing
	} else if file != nil {
    listing.List(file, res, req)
    return
  }

	http.NotFound(res, req)
}

