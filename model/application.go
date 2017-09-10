package model

import (
  "fmt"
	//"log"
  "os"
  "path"
  "strings"
  "io/ioutil"
  "path/filepath"
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

  SOURCE = "source"
  PUBLIC = "public"
)

var(
	types = []string{YAML, JSON}
)

type Application struct {
	// Mountpoint URL
	Url string `json:"url"`

  Path string `json:"-"`

  Name string `json:"name"`
  Description string `json:"description"`
  Pages []*Page `json:"-"`
  Files []*File `json:"-"`
	// The root file node, not included in the files slice.
	Root *File `json:"-"`
  Base string `json:"-"`
  Urls map[string] *File `json:"-"`

	FileSystem ApplicationFileSystem `json:"-"`

	// A protected application cannot be deleted.
	Protected bool `json:"protected,omitempty"`

  // Mark this application as a template
  IsTemplate bool `json:"is-template,omitempty"`

  ContainerName string `json:"container"`

	Container *Container `json:"-"`

	// A source template for this application
	Template *ApplicationTemplate `json:"template,omitempty"`

  // An application builder config loaded from build.yml.
  // For applications with no build file this is nil.
  Builder *BuildFile `json:"build,omitempty"`

  // Source file path
  sourcePath string

  // Public publish path
  publicPath string
}

func NewApplication(mountpoint, description string) *Application {
	return &Application{Url: mountpoint, Description: description}
}

func (app *Application) SourceDirectory() string {
  return app.sourcePath
}

func (app *Application) PublicDirectory() string {
  return app.publicPath
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

// Tests for files that would conflict when published, for example
// document.md and document.html would be published to the same location
// so we have to test for this when creating new files.
func (app *Application) ExistsConflict(url string) bool {
  var ext string = path.Ext(url)
  if TEMPLATE_FILE.MatchString(url) {
    url = strings.TrimSuffix(url, ext) + ".md"
    if app.Urls[url] != nil {
      return true
    }
  } else if MARKDOWN_FILE.MatchString(url) {
    url = strings.TrimSuffix(url, ext)
    if app.Urls[url + ".htm"] != nil || app.Urls[url + ".html"] != nil {
      return true
    }
  }
  return false
}

// Create a new file and publish it, the file cannot already exist on disc.
func (app *Application) Create(url string, content []byte) (*File, error) {
	path := app.GetPathFromUrl(url)

  /*
  println("create file: " + url)
  println("create file: " + path)
  println("create file app path: " + app.Path)
  */

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

  isDir := strings.HasSuffix(url, SLASH)

  file := app.NewFile(path, nil, content)
  if isDir {
    file.Directory = true
  }
	if err := app.FileSystem.SaveFile(file); err != nil {
		return nil, err
	}

  // Must add before publish for all fields to be available
	app.Add(file)

	if err := app.FileSystem.PublishFile(app.PublicDirectory(), file, &DefaultPublishFilter{}); err != nil {
		return nil, err
	}

	return file, nil
}

// Move a file to a new URL
func (app *Application) Move(file *File, newUrl string) error {
  println("move new url: " + newUrl)
  u := path.Clean(newUrl)
  if !strings.HasPrefix(u, "/") {
    u = "/" + u
  }

  if app.Urls[u] != nil {
    return fmt.Errorf("Cannot move file, destination %s exists", newUrl)
  }

  pth := app.GetPathFromUrl(u)

  // Move the source and published files
	if err := app.FileSystem.MoveFile(file, u, pth, nil); err != nil {
		return err
	}

  file.Path = pth
  delete(app.Urls, file.Url)
  app.setComputedFileFields(file)
  return nil
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
  if file.page != nil {
    if err := file.page.ParsePageData(); err != nil {
      return err
    }
  } else {
    // Update in-memory file data
    file.data = content
  }
	if err := app.FileSystem.SaveFile(file); err != nil {
		return err
	}
	if err := app.FileSystem.PublishFile(app.PublicDirectory(), file, &DefaultPublishFilter{}); err != nil {
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
func (app *Application) Add(file *File) error {
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
      if err := page.ParsePageData(); err != nil {
        return err
      }
		}
	}
  return nil
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
  file.Binary = isBinaryMime(file.Name, file.Mime)
	app.Files = append(app.Files, file)
	return len(app.Files)
}

// Add a page to this application.
func (app *Application) AddPage(page *Page) int {
	app.setComputedPageFields(page)
	page.owner = app
	page.Mime = getMimeType(page.Path)
  page.Binary = isBinaryMime(page.Name, page.Mime)
	app.Pages = append(app.Pages, page)
	return len(app.Pages)
}

func (app *Application) SetPath(path string) {
  app.Path = path
  app.sourcePath = filepath.Join(path, SOURCE)
  app.publicPath = filepath.Join(path, PUBLIC)
}

// Attempt to delete all application source and publish files.
func (app *Application) DeleteApplicationFiles() error {
  // Delete published files
  if err := os.RemoveAll(app.PublicDirectory()); err != nil {
    return err
  }
  // Delete source files, note this is the top-level directory not the
  // source sub-directory.
  if err := os.RemoveAll(app.Path); err != nil {
    return err
  }
  return nil
}

// Copy source files from another source application into this application.
func (app *Application) CopyApplicationTemplate(source *Application) error {
  var err error
  var files []*File = source.Files
  var prefix string

  for _, f := range files {
    if !f.Directory {
      rel := strings.TrimPrefix(f.Relative, prefix)
      out := filepath.Join(app.SourceDirectory(), rel)
      parent := filepath.Dir(rel)
      if parent != "/" {
        parent = filepath.Join(app.Path, parent)
        if err = os.MkdirAll(parent, os.ModeDir | 0755); err != nil {
          return err
        }
      }
      content := f.Source(true)
      ioutil.WriteFile(out, content, f.Info().Mode())
    }
  }
  return nil
}

// Load an application using the file system assigned to this application.
func (app *Application) Load(path string) error {
  var err error
  var builder *BuildFile
	app.Name = filepath.Base(path)
  app.SetPath(path)
  app.Urls = make(map[string] *File)

  if builder, err = ReadBuildFile(app); err != nil {
    return err
  }

  if builder != nil {
    app.Builder = builder
  }

  if err = app.FileSystem.Load(app.sourcePath); err != nil {
    return err
  }

  return nil
}

func (app *Application) HasBuilder() bool {
  return app.Builder != nil
}

func (app *Application) Build(done TaskComplete) {
  app.Builder.Build(done)
}

// Publish application files to the given directory.
func (app *Application) Publish(dir string) error {
  // Defer to a build task if defined
  if app.HasBuilder() {
    app.Build(&DefaultTaskComplete{})
    return nil
  }
  return app.FileSystem.Publish(dir, nil)
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
  url = path.Clean(url)
	base := app.SourceDirectory()
	parts := strings.Split(url, SLASH)
	//dest = filepath.Join(path.Split(dest))

	return base + SLASH + filepath.Join(parts...)
}

// Build a mountpoint URL by convention based on a container
// name.
func (app *Application) MountpointUrl(c *Container) string {
  if c == nil {
    c = app.Container
  }
  return "/apps/www/" + c.Name + "/" + app.Name + "/"
}

func (app *Application) PublishUrl() string {
  return app.Url
}

func (app *Application) SourceUrl() string {
	return "/apps/src/" + app.Container.Name + "/" + app.Name + "/"
}

func (app *Application) RawUrl() string {
  return "/apps/raw/" + app.Container.Name + "/" + app.Name + "/"
}

// Private methods

// Extract a name, relative path and URL for a file.
func (app *Application) getFileFields(file *File, base string) (string, string, string) {
	name := filepath.Base(file.Path)
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
	if file.info.IsDir() && !strings.HasSuffix(url, SLASH) {
		// url = strings.TrimSuffix(url, SLASH)
    url = url + SLASH
	}
	file.Url = url
	if !file.info.IsDir() {
		file.Size = file.info.Size()
	}
}

// Set computed fields for files.
func (app *Application) setComputedFileFields(file *File) {
	app.setFileFields(file, app.SourceDirectory())
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

