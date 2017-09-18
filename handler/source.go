package handler

import (
  "strings"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
)

// Serves application source files from memory.
type SourceHandler struct {
  Listing *DirList
	App *Application
	Raw bool
}

const(
	IndexName = "index.html"
)

// Tests the request path and attempts to find a corresponding source file
// in the application files.
func (h SourceHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	urls := h.App.Urls
	path := "/" + req.URL.Path
	trailing := strings.TrimSuffix(path, "/") + "/"
	index := trailing + IndexName

	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var file *File

	// Exact match
	if urls[path] != nil {
		file = urls[path]
  // Handler trailing slash redirect
  } else if urls[trailing] != nil && !strings.HasSuffix(path, "/") {
    redirect := http.RedirectHandler(req.URL.Path + "/", http.StatusMovedPermanently)
    redirect.ServeHTTP(res, req)
    return
	} else if urls[index] != nil {
		file = urls[index]
	}

  res.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

	if file != nil && !file.Info().IsDir() {
		output := file.Source(h.Raw)
    send(res, req, file, output)
		return
  // Handle directory listing
	} else if file != nil {
    h.Listing.List(file, res, req)
    return
  }

	http.NotFound(res, req)
}

