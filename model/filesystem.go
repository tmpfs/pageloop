package model

import(
	"errors"
)

type ApplicationFileSystem interface {
	Open(url string) error
}

type UrlFileSystem struct {
	App *Application
}

func NewUrlFileSystem(app *Application) *UrlFileSystem {
	return &UrlFileSystem{App: app}
}

func (fs *UrlFileSystem) Open(url string) (*File, error) {
	file := fs.App.Urls[url]
	if file == nil {
		return nil, errors.New("File not found at url " + url)
	}
	return file, nil
}
