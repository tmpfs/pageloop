package handler

import (
  "mime"
  "net/http"
  "path/filepath"
  "strconv"
  . "github.com/tmpfs/pageloop/model"
)

func send (res http.ResponseWriter, req *http.Request, file *File, output []byte) {
	path := "/" + req.URL.Path
  base := filepath.Base(path)

  ext := filepath.Ext(file.Name)
  ct := mime.TypeByExtension(ext)

  // TODO: remove this?
  if (ext == ".pdf") {
    res.Header().Set("Content-Disposition", "inline; filename=" + base)
  }

  res.Header().Set("Content-Type", ct)
  res.Header().Set("Content-Length", strconv.Itoa(len(output)))
  if (req.Method == http.MethodHead) {
    return
  }
  res.Write(output)
}

