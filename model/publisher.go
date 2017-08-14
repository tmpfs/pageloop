package model

import (
  "os"
	"strings"
  "path/filepath"
  "io/ioutil"
)

// Build directory.
var public string = "public"

// Abstract type to load files into an application.
type ApplicationPublisher interface {
  PublishApplication(app *Application, base string, filter PublishFilter) error
}

// Default implementation loads from the filesystem.
type FileSystemPublisher struct {}

type PublishFilter interface {
	Rename(path string) string
}

type DefaultFilter struct {}

func (f *DefaultFilter) Rename(path string) string {
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

// Publishes the application to a directory.
//
// Writes all application files using the current data bytes.
//
// Use base as the output directory, if base is the empty string a 
// public directory relative to the current working directory 
// is used.
func (p FileSystemPublisher) PublishApplication(app *Application, base string, filter PublishFilter) error {
  var err error
  var cwd string
	if filter == nil {
		filter = &DefaultFilter{}
	}
  if cwd, err = os.Getwd(); err != nil {
    return err
  }
  if base == "" {
    base = filepath.Join(cwd, public)
  }
  dir := filepath.Join(base, filepath.Base(app.Path))
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

    var rel string
    rel, err = filepath.Rel(app.Path, f.Path)
    if err != nil {
      return err
    }

    // Set output path and create parent directories
    out := filepath.Join(dir, rel)
		out = filter.Rename(out)
		// Remove this file from the output
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
  }

  return nil
}
