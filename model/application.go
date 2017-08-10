package model

import (
  "os"
  "bytes"
  "strings"
  "path/filepath"
  "io/ioutil"
  "encoding/json"
  "gopkg.in/yaml.v2"
  "github.com/tmpfs/pageloop/vdom"
)

const (
  // Default doctype.
  HTML5 = "html"

  // Page data file extensions.
  JSON string = ".json"
  YAML string = ".yml"

  // Page data types.
  DATA_NONE = iota
  DATA_YAML
  DATA_YAML_FILE
  DATA_JSON_FILE
)

var types = []string{YAML, JSON}

type Application struct {
  Path string `json: "path"`

  // The public file system path for the HTTP server.
  Public string `json: "public"`

  Name string `json: "name"`
  Title string `json:"title"`
  Pages []*Page `json:"pages"`
  Files []*File `json:"files"`
  Base string `json:"base"`
  Urls map[string] *File
}

type File struct {
  Path string `json:"path"`
  Directory bool `json:"directory"`
  Relative string `json:"relative"`
  Url string `json:"url"` 
  Index bool `json:"index"`
  info os.FileInfo
  data []byte
}

type Page struct {
  Path string `json: "path"`
  DocType string `json:"doctype"`
  UserData map[string] interface{} `json:"data"`
  UserDataType int
  Blocks []Block  `json:"blocks"`
  Dom *vdom.Vdom
  file *File
}

// Load an application using the given loader implementation, 
// if a nil loader is given the default file system loader is used.
func (app *Application) Load(path string, loader ApplicationLoader) error {
  var err error
  if loader == nil {
    loader = FileSystemLoader{}
  }
  err = loader.LoadApplication(path, app)
  if err != nil {
    return err
  }
  app.Path = path
  app.Urls = make(map[string] *File)
  app.setComputedFields(path)
  if err = app.merge(); err != nil {
    return err
  }
  return nil
}

// Publish files using the given publisher implementation, if a nil 
// publisher is givem the default file system publisher is used.
func (app *Application) Publish(publisher ApplicationPublisher) error {
  var err error
  if publisher == nil {
    publisher = FileSystemPublisher{}
  }

  var data []byte

  // Render pages to the file data bytes.
  for _, page := range app.Pages {
    if data, err = page.Render(); err != nil {
      return err
    }

    page.file.data = data
  }

  // TODO: allow setting base path for publish
  if err = publisher.PublishApplication(app, ""); err != nil {
    return err
  }

  return nil
}

// Determine a URL from a relative path.
func (app *Application) UrlFromPath(path string) string {
  var url string = strings.Join(strings.Split(path, string(os.PathSeparator)), "/")
  return url
}

// Private methods

// Merge user data with page structs loading user data from a JSON 
// file with the same name of the HTML file that created the page.
func (app *Application) merge() error {
  var err error
  for index, page := range app.Pages {
    if TEMPLATE_FILE.MatchString(page.file.Path) {
      if _, err = app.getPageData(page); err != nil {
        return err
      }
      app.Pages[index] = page
      app.Pages[index].Parse()
    }
  }
  return nil
}

// Set initial relative computed path and URL path.
//
// Also indicate whether a file is an index file and build the 
// map of URLs to files.
func (app *Application) setComputedFields(path string) Application {
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

// Attempt to find user page data by first attempting to 
// parse embedded frontmatter YAML.
// 
// If there is no frontmatter data it attempts to 
// load data from a corresponding file with a .yml extension.
//
// Finally if a .json file exists it is parsed.
func (app *Application) getPageData(page *Page) (map[string] interface{}, error) {
  page.UserData = make(map[string] interface{})
  page.UserDataType = DATA_NONE

  // frontmatter
  if FRONTMATTER.Match(page.file.data) {
    var read int = 4
    var lines [][]byte = bytes.Split(page.file.data, []byte("\n"))
    var frontmatter [][]byte
    // strip leading ---
    lines = lines[1:]
    for _, line := range lines {
      read += len(line) + 1
      if FRONTMATTER_END.Match(line) {
        break 
      }
      frontmatter = append(frontmatter, line)
    }

    if len(frontmatter) > 0 {
      fm := bytes.Join(frontmatter, []byte("\n"))
      err := yaml.Unmarshal(fm, &page.UserData)
      if err != nil {
        return nil, err
      }
      // strip frontmatter content from file data after parsing
      page.file.data = page.file.data[read:]

      page.UserDataType = DATA_YAML
    }
    return page.UserData, nil
  }

  // external files
  dir, name := filepath.Split(page.file.Path)
  for _, dataType := range types {
    dataPath := TEMPLATE_FILE.ReplaceAllString(name, dataType)
    dataPath = filepath.Join(dir, dataPath)
    fh, err := os.Open(dataPath)
    if err != nil {
      if !os.IsNotExist(err) {
        return nil, err
      }
    }
    if fh != nil {
      defer fh.Close()
      contents, err := ioutil.ReadFile(dataPath)
      if err != nil {
        return nil, err
      }

      if dataType == JSON {
        err = json.Unmarshal(contents, &page.UserData)
        if err != nil {
          return nil, err
        }
        page.UserDataType = DATA_JSON_FILE
      } else if dataType == YAML {
        err = yaml.Unmarshal(contents, &page.UserData)
        if err != nil {
          return nil, err
        }
        page.UserDataType = DATA_YAML_FILE
      }
      break
    }
  }
  return page.UserData, nil
}

