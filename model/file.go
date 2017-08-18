package model

import (
  "os"
	"mime"
  "path/filepath"
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

// TODO: http.File implementation
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.info, nil
}

func (f *File) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (f *File) Close() error {
	return nil
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

