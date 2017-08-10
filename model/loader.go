package model

import (
  "os"
  "path/filepath"
  "io/ioutil"
)

// Abstract type to load files into an application.
type ApplicationLoader interface {
  LoadApplication(path string, app *Application) error
}

// Default implementation loads from the filesystem.
type FileSystemLoader struct {}

// Loads the application assets from a filesystem directory path and 
// populate the given application with files and HTML pages.
func (r FileSystemLoader) LoadApplication(path string, app *Application) error {
  var err error
  err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    fh, err := os.Open(path)
    if err != nil {
      return err
    }
    defer fh.Close()

    stat, err := fh.Stat()
    if err != nil {
      return err
    }

    mode := stat.Mode()

    var file File
    if mode.IsDir() {
      file = File{Path: path, Directory: true, info: stat}
    } else if mode.IsRegular() {
      file = File{Path: path, info: stat}
      bytes, err := ioutil.ReadFile(path)
      if err != nil {
        return err
      }
      file.data = bytes
    }

    if TEMPLATE_FILE.MatchString(path) {
      page := Page{file: &file, Path: path}
      app.Pages = append(app.Pages, &page)
    }

    app.Files = append(app.Files, &file)

    return nil
  })
  return err
}
