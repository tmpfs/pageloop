package model

import(
  "os"
	"errors"
	"net/http"
  "path/filepath"
  "io/ioutil"
)

// Represents types that reference an application.
type ApplicationReference interface {
	App() *Application
}

// Type for a file system that references files by URL
// relative to an application.
type ApplicationFileSystem interface {
	ApplicationReference
	http.FileSystem
}

// Default file system that uses the underlying host file system.
type UrlFileSystem struct {
	app *Application
}

// Create a new URL file system.
func NewUrlFileSystem(app *Application) *UrlFileSystem {
	return &UrlFileSystem{app: app}
}

// Get the application reference.
func (fs *UrlFileSystem) App() *Application {
	return fs.app
}

// Get a pointer to a file and error if the file does not exist.
func (fs *UrlFileSystem) Open(url string) (http.File, error) {
	file := fs.App().Urls[url]
	if file == nil {
		return nil, errors.New("File not found at url " + url)
	}
	return file, nil
}

// Loads a file from disc and returns a new file reference 
// to the underlying file.
//
// The file reference has it's data and source set to the 
// loaded file contents.
func (fs *UrlFileSystem) LoadFilePath(path string) (*File, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	stat, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	mode := stat.Mode()

	var file File

	if mode.IsDir() {
		file = File{Path: path, Directory: true, info: stat}
	} else if mode.IsRegular() {
		file = File{Path: path, info: stat}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		file.data = bytes
		file.source = bytes
	}

	return &file, nil
}

// Recursively loads all files from the given directory and 
// adds them to the application.
func (fs *UrlFileSystem) Load(dir string) error {
	var app *Application = fs.App()
  var err error
  err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    var file *File
		if file, err = fs.LoadFilePath(path); err != nil {
			return err
		}
		app.Add(file)
    return nil
  })
  return err
}
