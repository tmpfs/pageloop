package pageloop

import (
  "net/http"
  . "github.com/tmpfs/pageloop/model"
)

type DirectoryList struct {
  Host *Host
}

func (d *DirectoryList) List(file *File, res http.ResponseWriter, req *http.Request) {
  if output, err := d.Host.DirectoryListing(file); err != nil {
    http.Error(res, err.Error(), http.StatusInternalServerError)
    return
  } else {
    // Send the directory listing
    send(res, req, file, output)
    return
  }
}
