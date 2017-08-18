package model

type ApplicationFileSystem interface {
	Open(url string) error
}

type UrlFileSystem struct {
	App *Application
}

//func (*fs ApplicationFileSystem)
