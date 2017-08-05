package blocks

import (
  "os"
  "path/filepath"
  "io/ioutil"
)

type Reader interface {
  Read(path string, app *Application) Application
}

type FileSystemReader struct {}

/*
  Reads an application's assets from a filesystem directory path.
*/
func (r FileSystemReader) Read(path string, app *Application) Application {
  filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
      page := Page{file: file}
      app.Pages = append(app.Pages, page)
    }

    app.Files = append(app.Files, file)

    return nil
  })
  return *app
}

