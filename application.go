package blocks

import (
  "os"
  "log"
  "bytes"
  "strings"
  "path/filepath"
  "io/ioutil"
  "encoding/json"
  "gopkg.in/yaml.v2"
)

const (
  JSON string = ".json"
  YAML string = ".yml"
)

var types = []string{YAML, JSON}

// Load an application using the given loader implementation, 
// if a nil loader is given the default file system loader is used.
func (app *Application) Load(path string, loader ApplicationLoader) Application {
  if loader == nil {
    loader = FileSystemLoader{}
  }
  loader.LoadApplication(path, app)
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
      app.GetUserData(&page)
      app.Pages[index] = page
      app.Pages[index].Parse()
    }
  }
  return *app
}

func (app *Application) GetUserData(page *Page) map[string] interface{} {
  page.UserData = make(map[string] interface{})

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
        log.Fatal(err)
      }
      // strip frontmatter content from file data after parsing
      page.file.data = page.file.data[read:]
    }
    return page.UserData
  }

  // external files
  dir, name := filepath.Split(page.file.Path)
  for _, dataType := range types {
    dataPath := TEMPLATE_FILE.ReplaceAllString(name, dataType)
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

      if dataType == JSON {
        err = json.Unmarshal(contents, &page.UserData)
        if err != nil {
          log.Fatal(err)
        }
      } else if dataType == YAML {
        err = yaml.Unmarshal(contents, &page.UserData)
        if err != nil {
          log.Fatal(err)
        }
      }
      break
    }
  }
  return page.UserData
}

/*
  Determine a URL from a relative path.
*/
func (app *Application) UrlFromPath(path string) string {
  var url string = strings.Join(strings.Split(path, string(os.PathSeparator)), "/")
  return url
}
