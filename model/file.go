package model

import (
	//"log"
  "os"
	"mime"
  "path/filepath"
  //"io/ioutil"
)

type File struct {
  Path string `json:"-"`
  Name string `json:"name"`
  Size int64 `json:"size,omitempty"`
  Url string `json:"url"`
  Directory bool `json:"dir,omitempty"`
  Relative string `json:"-"`
  Mime string `json:"mime"`

	// Owner application
	owner *Application

	// File stat information.
  info os.FileInfo

	// Raw source data.
	source []byte

	// Initially the source data but mutated later when
	// parsed from markdown or rendered from the vdom.
  data []byte
}

// Read only access to the source data outside this package.
func (f *File) Source() []byte {
	return f.source
}

// Read only access to the data outside this package.
func (f *File) Data() []byte {
	return f.data
}

// Read only access to the file info outside this package.
func (f *File) Info() os.FileInfo {
	return f.info
}

func getMimeType(path string) string {
	m := mime.TypeByExtension(filepath.Ext(path))
	if m == "" {
		m = "application/octet-stream"
	}
	return m
}

