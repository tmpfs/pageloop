package model

import (
  "os"
  //"fmt"
	"mime"
  "strings"
  "path/filepath"
)

const OCTET_STREAM = "application/octet-stream"

type File struct {
  Path string `json:"-"`
  Name string `json:"name"`
  Size int64 `json:"size,omitempty"`
  PrettySize string `json:"filesize,omitempty"`

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
  Binary bool `json:"binary"`

  // Destination for file move operations
  Destination string `json:"destination,omitempty"`

	// Owner application
  Owner *Application `json:"-"`

	// A source template for this file
	Template *ApplicationTemplate `json:"template,omitempty"`

  // A reference to a file in the form: file://{container}/{application}#{url}
  Ref string `json:"ref,omitempty"`

  // An input value for the file content, passed in when creating or
  // updating files that are not binary
  Value string `json:"value,omitempty"`

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

type DirectoryListing struct {
  Directories int
  Files int
  Length int
  Parent *File
  Children []*File
}

func (f *File) DirectoryListing () *DirectoryListing {
  listing := &DirectoryListing{Parent: f}
  for _, child := range f.Owner.Urls {
    if child == f {
      continue
    }
    if (strings.HasPrefix(child.Url, f.Url)) {
      if child.Directory {
        listing.Directories++
      } else {
        listing.Files++
      }
      listing.Children = append(listing.Children, child)
    }
  }
  listing.Length = len(listing.Children)
  return listing
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

func (f *File) Page() *Page {
  return f.page
}

// Read only access to the source data outside this package.
func (f *File) Source(raw bool) []byte {
  if raw && f.frontmatter != nil {
	  return append(f.frontmatter, f.source...)
  }
	return f.source
}

// Set the file source bytes.
func (f *File) Bytes(src []byte) {
  f.source = src
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
		m = OCTET_STREAM
	}
	return m
}

func isBinaryMime(name, mimeType string) bool {
  // This treats .gitignore, .babelrc etc as text
  ext := filepath.Ext(name)
  if ext == name {
    return false
  }

  if mimeType == OCTET_STREAM {
    return true
  }

  // application/ special cases
  if strings.HasPrefix(mimeType, "application/json") ||
    strings.HasPrefix(mimeType, "application/javascript") ||
    strings.HasPrefix(mimeType, "application/xml") ||
    strings.HasSuffix(mimeType, "+xml") {
    return false
  }

  if strings.HasPrefix(mimeType, "text/") {
    return false
  }

  return true
}
