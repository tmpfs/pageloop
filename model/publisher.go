package model

import (
  "os"
  "path/filepath"
  "io/ioutil"
)

// Build directory.
var public string = "public"

// Abstract type to load files into an application.
type ApplicationPublisher interface {
  PublishApplication(app *Application, base string) error
}

// Default implementation loads from the filesystem.
type FileSystemPublisher struct {}

// Loads the application assets from a filesystem directory path and 
// populate the given application with files and HTML pages.
//
// Use base as the output directory, if base is the empty string a 
// public directory relative to the current working directory 
// is used.
func (p FileSystemPublisher) PublishApplication(app *Application, base string) error {
  var err error
  var cwd string
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
