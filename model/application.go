package model

import (
	//"log"
  "os"
  "bytes"
	"regexp"
  "strings"
  "path/filepath"
  "io/ioutil"
  "encoding/json"
  "gopkg.in/yaml.v2"
)

const (
	SLASH = "/"

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

type Application struct {
	// Mountpoint URL
	Url string `json:"url"`

  Path string `json:"-"`

  // The public file system path for the HTTP server.
  Public string `json:"-"`

  Name string `json:"name"`
  Description string `json:"description"`
  Pages []*Page `json:"-"`
  Files []*File `json:"-"`
	// The root file node, not included in the files slice.
	Root *File `json:"-"`
  Base string `json:"-"`
  Urls map[string] *File `json:"-"`

	FileSystem ApplicationFileSystem

	// A protected application cannot be deleted.
	Protected bool `json:"protected,omitempty"`

	Container *Container `json:"-"`
}

func NewApplication(mountpoint, description string) *Application {
	return &Application{Url: mountpoint, Description: description}
}

// Determine the page type for an input file path.
func (app *Application) GetPageType(path string) int {
	var pageType int = PageNone
	if !VENDOR.MatchString(path) {
		if TEMPLATE_FILE.MatchString(path) {
			pageType = PageHtml
		} else if MARKDOWN_FILE.MatchString(path) {
			pageType = PageMarkdown
		}
	}
	return pageType
}

// Delete a file
func (app *Application) Del(file *File) error {
	println("Deleting file: " + file.Url)
	return nil
}

// Add a file or page inspecting the file path to determine
// how to add the file.
//
// Note that pages also exist in the list of all files.
func (app *Application) Add(file *File) {
	var pageType int = app.GetPageType(file.Path)

	// TODO: merge and set computed properties!!!

	if file.Path == app.Path {
		app.Root = file
	} else {
		// Must add the file before page for computed proxied fields
		app.AddFile(file)

		// Add to the list of pages
		if pageType != PageNone {
			page := Page{file: file, Path: file.Path, Type: pageType}
			app.AddPage(&page)
		}
	}
}

// Create a new file.
func (app *Application) NewFile(path string, info os.FileInfo, data []byte) *File {
	return &File{Path: path, info: info, data: data, source: data}
}

// Add a file to this application.
func (app *Application) AddFile(file *File) int {
	app.setComputedFileFields(file)
	file.owner = app
	file.Mime = getMimeType(file.Path)
	app.Files = append(app.Files, file)
	return len(app.Files)
}

// Add a page to this application.
func (app *Application) AddPage(page *Page) int {
	app.setComputedPageFields(page)
	page.owner = app
	page.Mime = getMimeType(page.Path)
	app.Pages = append(app.Pages, page)
	return len(app.Pages)
}

// Load an application using the file system assigned to this application.
func (app *Application) Load(path string) error {
  var err error
	app.Name = filepath.Base(path)
  app.Path = path
  app.Urls = make(map[string] *File)

  err = app.FileSystem.Load(path)
  if err != nil {
    return err
  }

  if err = app.merge(); err != nil {
    return err
  }

  return nil
}

// Publish files using the given publisher implementation, if a nil
// publisher is given the default file system publisher is used.
func (app *Application) Publish(publisher ApplicationPublisher, base string) error {
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

  if err = publisher.PublishApplication(app, base, nil); err != nil {
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
		u = strings.TrimSuffix(u, SLASH)
		if u == url {
			return page
		}
	}
	return nil
}

// Private methods

// Determine a URL from a relative path.
func (app *Application) getUrlFromPath(file *File, relative string) string {
	var url string = filepath.ToSlash(relative)
	if file.info.IsDir() && !strings.HasSuffix(url, SLASH) {
		url += "/"
	}
  return url
}

// Get an absolute path from a relative URL reference.
func (app *Application) GetPathFromUrl(url string) string {
	base := app.Path
	parts := strings.Split(url, SLASH)
	//dest = filepath.Join(path.Split(dest))

	return base + SLASH + filepath.Join(parts...)
}

// Merge user data with page structs loading user data from a JSON
// file with the same name of the HTML file that created the page.
func (app *Application) merge() error {
  var err error
  for index, page := range app.Pages {
		if _, err = app.getPageData(page); err != nil {
			return err
		}
		if _, err = page.Parse(page.file.data); err != nil {
			return err
		}
		app.Pages[index] = page

		// Convert map[interface{}] interface{} to map[string] interface{}
		// recursively in parsed page data.
		//
		// This is necessary as the YAML unmarshaller converts nested maps
		// using interface{} keys and interface{} keys are not recognised when
		// marshalling to JSON. We don't want to define a type for page data as
		// it is intended to be arbitrary.
		var coerce func(m map[string] interface{})
		coerce = func(m map[string] interface{}) {
			for key, value := range m {
				if val, ok := value.(map[interface{}] interface{}); ok {
					var r map[string]interface{}
					r = make(map[string]interface{})
					for k, v := range val {
						r[k.(string)] = v
					}
					coerce(r)
					m[key] = r
				}
			}
		}
		coerce(page.PageData)
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
	if strings.HasSuffix(url, SLASH) {
		url = strings.TrimSuffix(url, SLASH)
	}
	file.Url = url
	if !file.info.IsDir() {
		file.Size = file.info.Size()
	}
}

// Set computed fields for files.
func (app *Application) setComputedFileFields(file *File) {
	path := app.Path
	app.setFileFields(file, path)
	app.Urls[file.Url] = file
}

// Set computed fields for pages, the underlying file must
// have had it's computed fields set.
func (app *Application) setComputedPageFields(page *Page) {
	file := page.file
	page.Name = file.Name
	page.Url = file.Url
	page.Size = file.info.Size()
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
      println("parsing frontmatter")
      println(page.Path)
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
    dataPath := name
    if TEMPLATE_FILE.MatchString(name) {
      dataPath = TEMPLATE_FILE.ReplaceAllString(name, dataType)
    } else if MARKDOWN_FILE.MatchString(name) {
      dataPath = MARKDOWN_FILE.ReplaceAllString(name, dataType)
    }
    dataPath = filepath.Join(dir, dataPath)
    // Failed to change file extension
    if dataPath == page.file.Path {
      return nil, nil
    }

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

