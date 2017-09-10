package pageloop

import (
  "strings"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
)

// Serves application public files from disc.
type PublicHandler struct {
	App *Application
  FileServer http.Handler
}

func (h PublicHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  app := h.App
  path := "/" + req.URL.Path
  file := app.Urls[path]
	clean := strings.TrimSuffix(path, "/")
  // FIXME: this is rubbish
	indexPage := clean + "/index.html"
	indexMdPage := clean + "/index.md"
  if file != nil && file.Directory && app.Urls[indexPage] == nil && app.Urls[indexMdPage] == nil {
    listing.List(file, res, req)
    return
  }

  // Defer to file server for files
  h.FileServer.ServeHTTP(res, req)
}
