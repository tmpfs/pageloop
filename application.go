package blocks

import (
  "os"
  //"fmt"
  "regexp"
  "path/filepath"
  "io/ioutil"
)

var templateFile = regexp.MustCompile(`\.html?$`)

type File struct {
  Path string `json:"path"`
  Directory bool `json:"directory"`
  info os.FileInfo
  data []byte
}

type Application struct {
  Title string `json:"title"`
  Pages []Page `json:"pages"`
  Files []File `json:"files"`
}

func (app *Application) Load(path string) Application {
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

    if templateFile.MatchString(path) {
      page := Page{file: file}
      app.Pages = append(app.Pages, page)
    } else {
      app.Files = append(app.Files, file)
    }

    return nil
  })

  return *app
}
