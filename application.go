package blocks

import (
  "os"
  "fmt"
  "path/filepath"
  "io/ioutil"
)

type File struct {
  path string `json:"path"`
  directory bool `json:"directory"`
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
      file = File{directory: true}
      fmt.Println(path)
      fmt.Println(stat)
    } else if mode.IsRegular() {
      file = File{path: path, info: stat}

      bytes, err := ioutil.ReadFile(path)
      if err != nil {
        return err
      }
      file.data = bytes
      fmt.Println(path)
    }
    return nil
  })

  return *app
}
