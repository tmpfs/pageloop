package model

import (
  "os"
	"fmt"
  "bytes"
	"errors"
	"regexp"
  "strings"
  "path/filepath"
  "io/ioutil"
  "encoding/json"
  "gopkg.in/yaml.v2"
  "github.com/tmpfs/pageloop/vdom"
)

const (
  // Page data file extensions.
  JSON string = ".json"
  YAML string = ".yml"

  // Page data types.
  DATA_NONE = iota
  DATA_YAML
  DATA_YAML_FILE
  DATA_JSON_FILE
)

var(
	types = []string{YAML, JSON}
	ptn string = `^[a-zA-Z0-9]+[a-zA-Z0-9-]*`
	re = regexp.MustCompile(ptn)
)

// Represents a collection of applications.
type Container struct {
	Apps []*Application `json:"apps"`
}

// Add an application to the container, the application must 
// have the Name field set and it must not already exist in 
// the container list.
//
// Application names may only contain lowercase, uppercase, hyphens 
// and digits. They may not begin with a hyphen.
func (c *Container) Add(app *Application) error {
	if app.Name == "" {
		return errors.New("Application name is required to add to container")
	}

	if !re.MatchString(app.Name) {
		return errors.New(fmt.Sprintf("Application name must match pattern %s", ptn))
	}

	var exists *Application = c.GetByName(app.Name)
	if exists != nil {
		return errors.New(fmt.Sprintf("Application exists with name %s", app.Name))
	}

	c.Apps = append(c.Apps, app)

	return nil
}

// Get an application by name.
func (c *Container) GetByName(name string) *Application {
	for _, app := range c.Apps {
		if app.Name == name {
			return app
		}
	}
	return nil
}

type Application struct {
  Path string `json:"-"`

  // The public file system path for the HTTP server.
  Public string `json:"-"`

  Name string `json:"name"`
  Pages []*Page `json:"-"`
  Files []*File `json:"-"`
	// The root file node, not included in the files slice.
	Root *File `json:"-"`
  Base string `json:"-"`
  Urls map[string] *File `json:"-"`
}

type File struct {
  Path string `json:"-"`
  Name string `json:"name"` 
  Size int64 `json:"size,omitempty"` 
  Url string `json:"url"` 
  Directory bool `json:"dir,omitempty"`
  Relative string `json:"-"`
	//Mime string `json:"mime,omitempty"`
  //Index bool `json:"index"`
  info os.FileInfo
  data []byte
}

type Page struct {
  Path string `json:"-"`
  Name string `json:"name"` 
  Url string `json:"url"` 
  Size int64 `json:"size,omitempty"` 
  PageData map[string] interface{} `json:"data"`
	PageDataType int `json:"-"`
  Blocks []Block  `json:"blocks"`
  Dom *vdom.Vdom `json:"-"`
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
	app.Name = filepath.Base(path)
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
    node := page.Dom.Clean(nil)
    if data, err = page.Render(page.Dom, node); err != nil {
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

// Get a file pointer by URL.
func (app *Application) GetFileByUrl(url string) *File {
	return app.Urls[url]
}

// Get a page pointer by URL.
func (app *Application) GetPageByUrl(url string) *Page {
	for _, page := range app.Pages {
		u := page.Url
		u = strings.TrimSuffix(u, "/")
		if u == url {
			return page
		}
	}
	return nil
}

// Private methods

// Determine a URL from a relative path.
func (app *Application) getUrlFromPath(file *File, relative string) string {
  var url string = strings.Join(strings.Split(relative, string(os.PathSeparator)), "/")
	if file.info.IsDir() && !strings.HasSuffix(url, "/") {
		url += "/"
	}
  return url
}

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

// Extract a name, relative path and URL for a file.
func (app *Application) getFileFields(file *File, base string) (string, string, string) {
	name := file.info.Name()
	relative := strings.TrimPrefix(file.Path, base)
  url := app.getUrlFromPath(file, relative)
	return name, relative, url
}

// Set computed fields on a file.
func (app *Application) setFileFields(file *File, base string) {
	name, relative, url := app.getFileFields(file, base)
	file.Name = name
	file.Relative = relative
	// Normalize directories for URLs without the /
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	file.Url = url
	if !file.info.IsDir() {
		file.Size = file.info.Size()
	}
}

// Set computed fields.
//
func (app *Application) setComputedFields(path string) {
  for _, file := range app.Files {
		app.setFileFields(file, path)
    app.Urls[file.Url] = file
	}

  for _, page := range app.Pages {
		file := page.file
		page.Name = file.Name
		page.Url = file.Url
		page.Size = file.info.Size()
	}
}

// Attempt to find user page data by first attempting to 
// parse embedded frontmatter YAML.
// 
// If there is no frontmatter data it attempts to 
// load data from a corresponding file with a .yml extension.
//
// Finally if a .json file exists it is parsed.
func (app *Application) getPageData(page *Page) (map[string] interface{}, error) {
  page.PageData = make(map[string] interface{})
  page.PageDataType = DATA_NONE

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
      err := yaml.Unmarshal(fm, &page.PageData)
      if err != nil {
        return nil, err
      }
      // strip frontmatter content from file data after parsing
      page.file.data = page.file.data[read:]

      page.PageDataType = DATA_YAML
    }
    return page.PageData, nil
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
        err = json.Unmarshal(contents, &page.PageData)
        if err != nil {
          return nil, err
        }
        page.PageDataType = DATA_JSON_FILE
      } else if dataType == YAML {
        err = yaml.Unmarshal(contents, &page.PageData)
        if err != nil {
          return nil, err
        }
        page.PageDataType = DATA_YAML_FILE
      }
      break
    }
  }
  return page.PageData, nil
}

