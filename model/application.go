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

// Create a new file and publish it, the file cannot already exist on disc.
func (app *Application) Create(path string, content []byte) (*File, error) {
	// TODO: check path is not outside the application
  println("create file: " + path)
  println("create file app path: " + app.Path)

	var err error
	var fh *os.File
	// The file must not exist in order to create
	if fh, err = os.Open(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	if fh != nil {
		defer fh.Close()
	}
	var file *File = app.NewFile(path, nil, content)
  // TODO: remove file from list on error!
	if err := app.FileSystem.SaveFile(file); err != nil {
		return nil, err
	}

  // Must add before publish for all fields to be available
	app.Add(file)

	if err := app.FileSystem.PublishFile(app.Public, file, &DefaultPublishFilter{}); err != nil {
		return nil, err
	}

	return file, nil
}

// Update an existing file source and publish it, file must already exist on disc.
func (app *Application) Update(file *File, content []byte) error {
	var err error
	var fh *os.File
	// The file must exist in order to create
	if fh, err = os.Open(file.Path); err != nil {
		return err
	}
	defer fh.Close()
	file.source = content
	if err := app.FileSystem.SaveFile(file); err != nil {
		return err
	}
	if err := app.FileSystem.PublishFile(app.Public, file, &DefaultPublishFilter{}); err != nil {
		return err
	}
	return nil
}

// Delete a file.
//
// The file is removed from the URL map and the list of files
// for this application. If the file is also a page it is removed
// from the page list.
//
// Source and published versions are deleted from the filesystem.
func (app *Application) Del(file *File) error {
	// Remove from the URL map
	delete(app.Urls, file.Url)

	// Remove from the list of pages
	for i, p := range app.Pages {
		if p.file == file {
      before := app.Pages[0:i]
      after := app.Pages[i+1:]
      app.Pages = append(before, after...)
		}
	}

	// Remove from the list of files
	for i, f := range app.Files {
		if f == file {
      before := app.Files[0:i]
      after := app.Files[i+1:]
      app.Files = append(before, after...)
		}
	}

	/*
	if file.Directory {
		return app.FileSystem.RemoveAll(file)
	}
	*/

	return app.FileSystem.Remove(file)
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
			page := &Page{file: file, Path: file.Path, Type: pageType}
			file.page = page
			app.AddPage(page)
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

// Publish application files to the given directory.
func (app *Application) Publish(dir string) error {
  var err error

  // Render pages to the file data bytes.
	/*
	for _, page := range app.Pages {
		if err = page.Update(); err != nil {
			return err
		}
	}
	*/

  if err = app.FileSystem.Publish(dir, nil); err != nil {
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
func (app *Application) GetUrlFromPath(file *File, relative string) string {
	var url string = filepath.ToSlash(relative)
	if file.info.IsDir() && !strings.HasSuffix(url, SLASH) {
		url += "/"
	}
	if !strings.HasPrefix(url, SLASH) {
		url = "/"	+ url
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
  url := app.GetUrlFromPath(file, relative)
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
      err := yaml.Unmarshal(fm, &page.PageData)
      if err != nil {
        return nil, err
      }
      // strip frontmatter content from file data after parsing
      page.file.data = page.file.data[read:]
      page.file.source = page.file.source[read:]

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

