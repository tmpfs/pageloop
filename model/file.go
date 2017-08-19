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

  // This is a slash separated path to the source file
  // relative to the application base, it will start with
  // a leading slash. Directories will always have a trailing
  // slash.
  Url string `json:"url"`

  // This is a slash separated path to the public URI for
  // the published file, normally it will be the same as the
  // url however if a publish file filter mutates the path
  // then this will point to the mutated relative path.
  Uri string `json:"uri"`

  Directory bool `json:"dir,omitempty"`
  Relative string `json:"-"`
  Mime string `json:"mime"`

	// Owner application
	owner *Application

	// File stat information.
  info os.FileInfo

  // Frontmatter content
  frontmatter []byte

	// Raw source data (frontmatter is removed)
	source []byte

	// Initially the source data but mutated later when
	// parsed from markdown or rendered from the vdom.
  data []byte

	// A corresponding page if this file represents a page
	page *Page
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
func (f *File) Source(raw bool) []byte {
  if raw && f.frontmatter != nil {
	  return append(f.frontmatter, f.source...)
  }
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

