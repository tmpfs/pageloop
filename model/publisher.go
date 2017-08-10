package model

import (
  "os"
  "path/filepath"
  //"io/ioutil"
)

// Build directory.
var build string = ".build"

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
  println("publishing app", app.Path)
  dir := filepath.Join(app.Path, build)
  fh, err := os.Open(dir)
  if err != nil {
    if !os.IsNotExist(err) {
      return err
    // Try to make the directory.
    } else {
      if err = os.Mkdir(dir, os.ModeDir); err != nil {
        return err
      }
    }
  }
  defer fh.Close()

  println(dir)

  return err
}
