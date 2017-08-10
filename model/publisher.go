package model

import (
  //"os"
  //"path/filepath"
  //"io/ioutil"
)

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
  return err
}
