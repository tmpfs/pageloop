package blocks

import (
  "os"
  "log"
  "strings"
  "path/filepath"
  "io/ioutil"
  "encoding/json"
)

/*
  Load an application using the given Reader implementation.

  If a nil reader is given the default file system reader is used.
*/
func (app *Application) Load(path string, reader Reader) Application {
  if reader == nil {
    reader = FileSystemReader{}
  }
  reader.Read(path, app)
  app.Urls = make(map[string] File)
  app.SetComputedFields(path)
  app.Merge()
  return *app
}


/*
  Set initial relative computed path and URL path.

  Also indicate whether a file is an index file and build the 
  map of URLs to files.
*/
func (app *Application) SetComputedFields(path string) Application {
  for _, file := range app.Files {
    // includes the leading slash
    file.Relative = strings.TrimPrefix(file.Path, path)
    if app.Name != "" {
      file.Relative = "/" + app.Name + file.Relative
    }
    file.Url = app.UrlFromPath(file.Relative)

    if INDEX_FILE.MatchString(file.Path) {
      file.Index = true
    }

    app.Urls[file.Url] = file
  }
  return *app
}

/*
  Merge user data with page structs loading user data from a JSON
  file with the same name of the HTML file that created the page.
*/
func (app *Application) Merge() Application {
  for index, page := range app.Pages {
    if TEMPLATE_FILE.MatchString(page.file.Path) {
      dir, name := filepath.Split(page.file.Path)
      dataPath := TEMPLATE_FILE.ReplaceAllString(name, ".json")
      dataPath = filepath.Join(dir, dataPath)
      fh, err := os.Open(dataPath)
      if err != nil {
        if !os.IsNotExist(err) {
          log.Fatal(err)
        }
      }
      if fh != nil {
        defer fh.Close()
        contents, err := ioutil.ReadFile(dataPath)
        if err != nil {
          log.Fatal(err)
        }
        page.UserData = make(map[string] interface{})
        err = json.Unmarshal(contents, &page.UserData)
        if err != nil {
          log.Fatal(err)
        }
        log.Printf("%+v\n", page.UserData)
        app.Pages[index] = page
      }
    }
  }

  return *app
}

/*
  Determine a URL from a relative path.
*/
func (app *Application) UrlFromPath(path string) string {
  var url string = strings.Join(strings.Split(path, string(os.PathSeparator)), "/")
  return url
}
