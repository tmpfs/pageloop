package model

import(
  "os"
	"errors"
	"strings"
	"net/http"
  "path/filepath"
  "io/ioutil"
)

// Build directory.
var public string = "public"

// Represents types that reference an application.
type ApplicationReference interface {
	App() *Application
}

// Type for a file system that references files by URL
// relative to an application.
type ApplicationFileSystem interface {
	ApplicationReference
	http.FileSystem

	LoadFile(path string) (*File, error)
	Load(dir string) error

	PublishFile(dir string, f *File, filter FileFilter) error
	Publish(dir string, filter FileFilter) error

	// Remove the source and published files for a file reference
	Remove(f *File) error
	RemoveAll(f *File) error
}

// Default file system that uses the underlying host file system.
type UrlFileSystem struct {
	app *Application
}

// Represents types that can change the name or path of files.
//
// Implementations may return the empty string to indicate the 
// file should be ignored.
type FileFilter interface {
	Rename(path string) string
}

type DefaultPublishFilter struct {}

// Default file filter used during publishing.
func (f *DefaultPublishFilter) Rename(path string) string {
	name := filepath.Base(path)
	if name == Layout {
		return ""
	}
	ext := filepath.Ext(path)
	if ext == ".md" || ext == ".markdown" {
		name = strings.TrimSuffix(name, ext)
		return filepath.Join(filepath.Dir(path), name + ".html")
	}
	return path
}

// Create a new URL file system.
func NewUrlFileSystem(app *Application) *UrlFileSystem {
	return &UrlFileSystem{app: app}
}

// Get the application reference.
func (fs *UrlFileSystem) App() *Application {
	return fs.app
}

// Get a pointer to a file and error if the file does not exist.
func (fs *UrlFileSystem) Open(url string) (http.File, error) {
	file := fs.App().Urls[url]
	if file == nil {
		return nil, errors.New("File not found at url " + url)
	}
	return file, nil
}

// Loads a file from disc and returns a new file reference 
// to the underlying file.
//
// The file reference has it's data and source set to the 
// loaded file contents.
func (fs *UrlFileSystem) LoadFile(path string) (*File, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	stat, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	mode := stat.Mode()

	var file File

	if mode.IsDir() {
		file = File{Path: path, Directory: true, info: stat}
	} else if mode.IsRegular() {
		file = File{Path: path, info: stat}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		file.data = bytes
		file.source = bytes
	}

	return &file, nil
}

// Recursively loads all files from the given directory and 
// adds them to the application.
func (fs *UrlFileSystem) Load(dir string) error {
	var app *Application = fs.App()
  var err error
  err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    var file *File
		if file, err = fs.LoadFile(path); err != nil {
			return err
		}
		app.Add(file)
    return nil
  })
  return err
}

// Publish a single file relative to the given directory.
func (fs *UrlFileSystem) PublishFile(dir string, f *File, filter FileFilter) error {
	app := fs.App()
	var err error
	rel, err := filepath.Rel(app.Path, f.Path)
	if err != nil {
		return err
	}

	// Set output path and create parent directories
	out := filepath.Join(dir, rel)

	println("publishing to: " + out)

	out = filter.Rename(out)
	// Ignore publishing this file
	if out == "" {
		return nil
	}

	parent := out
	isDir := f.info.Mode().IsDir()
	if !isDir {
		parent = filepath.Dir(out)
	}
	if err = os.MkdirAll(parent, os.ModeDir | 0755); err != nil {
		return err
	}

	// Write out the file data
	if !isDir {
		mode := f.info.Mode()
		if err = ioutil.WriteFile(out, f.data, mode); err != nil {
			return err
		}
	}
	return nil
}

// Publishes the application to a directory.
//
// Writes all application files using the current data bytes.
//
// Use dir as the output directory, if dir is the empty string a
// public directory relative to the current working directory
// is used.
//
// If a nil file filter is given the default publish filter is used.
func (fs *UrlFileSystem) Publish(dir string, filter FileFilter) error {
	var app *Application = fs.App()
  var err error
  var cwd string
	if filter == nil {
		filter = &DefaultPublishFilter{}
	}
  if cwd, err = os.Getwd(); err != nil {
    return err
  }
  if dir == "" {
    dir = filepath.Join(cwd, public)
  }
  dir = filepath.Join(dir, filepath.Base(app.Path))
  fh, err := os.Open(dir)
  if err != nil {
    if !os.IsNotExist(err) {
      return err
    // Try to make the directory.
    } else {
      if err = os.MkdirAll(dir, os.ModeDir | 0755); err != nil {
        return err
      }
    }
  }
  defer fh.Close()

	// TODO: remove this and assign outside the publisher
  app.Public = dir

  for _, f := range app.Files {
    // Ignore the build directory
    if f.Path == app.Path {
      continue
    }
		if err = fs.PublishFile(dir, f, filter); err != nil {
			return err
		}
  }
  return nil
}

// Remove the source and published files from the file system.
func (fs *UrlFileSystem) Remove(f *File) error {
	app := fs.App()
	src := f.Path
	pub := filepath.Join(app.Public, f.Relative)

	println("fs delete: " + src)
	println("fs delete: " + pub)

	if err := os.Remove(pub); err != nil {
		return err
	}

	return os.Remove(src)
}

// Recursively deletes a directory including source and 
// published versions of the files.
func (fs *UrlFileSystem) RemoveAll(f *File) error {
	app := fs.App()
	src := f.Path
	pub := filepath.Join(app.Public, f.Relative)

	println("fs delete: " + src)
	println("fs delete: " + pub)

	return nil
}
