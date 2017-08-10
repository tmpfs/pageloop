package model

import (
  "fmt"
  "os"
  "path/filepath"
  "io/ioutil"
)

// Build directory.
var build string = "build"
var current string = "current"

// Abstract type to load files into an application.
type ApplicationPublisher interface {
  PublishApplication(app *Application) error
}

// Default implementation loads from the filesystem.
type FileSystemPublisher struct {}

// Loads the application assets from a filesystem directory path and 
// populate the given application with files and HTML pages.
func (p FileSystemPublisher) PublishApplication(app *Application) error {
  var err error
  var base string
  if base, err = os.Getwd(); err != nil {
    return err
  }
  dir := filepath.Join(base, build)
  dir = filepath.Join(dir, current)
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

  println("publishing app", app.Path)
  println(dir)

  fmt.Printf("files: %#v\n", app.Files)

  for _, f := range app.Files {
    if f.Path == app.Path {
      continue
    }
    fmt.Printf("%#v\n", f)
    var rel string
    println(f.Path)
    rel, err = filepath.Rel(app.Path, f.Path)
    if err != nil {
      return err
    }

    println(app.Path)
    println(f.Path)
    println("relative path")
    fmt.Println(rel)

    // Set output path and create parent directories
    out := filepath.Join(dir, rel)
    parent := filepath.Dir(out)
    println(parent)
    if f.info.Mode().IsDir() {
      if err = os.MkdirAll(parent, os.ModeDir | 0755); err != nil {
        return err
      }
    } 

    // Write out the file data
    mode := f.info.Mode()
    if err = ioutil.WriteFile(out, f.data, mode); err != nil {
      return err
    }
  }

  return nil
}
