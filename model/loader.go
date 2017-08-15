package model

import (
  "os"
  "path/filepath"
  "io/ioutil"
)

// Abstract type to load files into an application.
type ApplicationLoader interface {
  LoadApplication(path string, app *Application) error
}

// Default implementation loads from the filesystem.
type FileSystemLoader struct {}

// Loads the application assets from a filesystem directory path and
// populate the given application with files and pages.
func (r FileSystemLoader) LoadApplication(dir string, app *Application) error {
  var err error
  err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
		var pageType int = PageNone

    if mode.IsDir() {
      file = File{Path: path, Directory: true, info: stat}
    } else if mode.IsRegular() {
      file = File{Path: path, info: stat}
      bytes, err := ioutil.ReadFile(path)
      if err != nil {
        return err
      }
      file.data = bytes
			file.source = bytes

			if TEMPLATE_FILE.MatchString(path) {
				pageType = PageHtml
			} else if MARKDOWN_FILE.MatchString(path) {
				pageType = PageMarkdown
			}

    }

		if pageType != PageNone {
			page := Page{file: &file, Path: path, Type: pageType}
			app.AddPage(&page)
		}

		//
		if dir == path {
			app.Root = &file
		} else {
			app.AddFile(&file)
		}

    return nil
  })
  return err
}
