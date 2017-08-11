package model

import (
	"fmt"
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

  //fmt.Printf("files: %#v\n", app.Files)
	fmt.Printf("files: %#v\n", app.Public)

  for _, f := range app.Files {
    // Ignore the build directory
    if f.Path == app.Path {
      continue
    }

		println("Write path: " + f.Path)
		fmt.Println("Write path: ",f.info.Mode().IsDir())

    var rel string
    rel, err = filepath.Rel(app.Path, f.Path)
    if err != nil {
      return err
    }

    // Set output path and create parent directories
    out := filepath.Join(dir, rel)
		fmt.Println("Publish output", out)
    parent := filepath.Dir(out)
    if f.info.Mode().IsDir() {
      if err = os.MkdirAll(parent, os.ModeDir | 0755); err != nil {
        return err
      }
			continue
    } 

    // Write out the file data
    mode := f.info.Mode()
    if err = ioutil.WriteFile(out, f.data, mode); err != nil {
			println("Writing file...")
      return err
    }
  }

  return nil
}
