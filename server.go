// System for hyperfast HTML document editing.
//
// Stores HTML documents on the server as in-memory DOM 
// documents that can be modified on the client. The client 
// provides an editor view and a preview of the rendered 
// page loaded in an iframe.
package pageloop

import (
  "log"
	"errors"
  "net/http"
  "path/filepath"
	"regexp"
  "time"
  "github.com/tmpfs/pageloop/model"
  //"github.com/elazarl/go-bindata-assetfs"
)

var config ServerConfig
var mux *http.ServeMux

type Mountpoint struct {
	// The URL path component.
	UrlPath string `json:"url"`
	// The path to pass to the loader.
	Path string	`json:"path"`
}

type PageLoop struct {
	// Underlying HTTP server.
	Server *http.Server `json:"-"`

	// All application mountpoints.
  Mountpoints []Mountpoint `json:"-"`

	// All mounted applications.
	Apps []*model.Application `json:"apps"`
}

type ServerConfig struct {
	Addr string

	// List of user application mountpoints.
  Mountpoints []Mountpoint

	// Load system assets from the file system, don't use 
	// the embedded assets.
	Dev bool
}

type ServerHandler struct {}

// The default server handler, defers to a multiplexer.
func (h ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  handler, _ := mux.Handler(req)
  handler.ServeHTTP(res, req)
}

// Starts an HTTP server listening.
func (l *PageLoop) NewServer(config ServerConfig) (*http.Server, error) {
  var err error

  // Initialize server mux
  mux = http.NewServeMux()

	// System applications to mount.
	l.Mountpoints = append(l.Mountpoints, Mountpoint{UrlPath: "/", Path: "data://app/home"})

	// Add user applications.
	l.Mountpoints = append(l.Mountpoints, config.Mountpoints...)

	var apps []*model.Application
  if apps, err = l.LoadApps(config); err != nil {
    return nil, err
  }

  s := &http.Server{
    Addr:           config.Addr,
    Handler:        ServerHandler{},
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

	l.Server = s
	l.Apps = apps

  return s, nil
}

func (l *PageLoop) Listen() error {
	var err error
	s := l.Server
	if s == nil {
		return errors.New("Cannot listen without a server, call NewServer().")
	}
  if err = s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (l *PageLoop) LoadApps(config ServerConfig) ([]*model.Application, error) {
  var err error

	// Global endpoints
	rest := RestHandler{Loop: l}
	api := "/api/"
	mux.Handle(api, http.StripPrefix(api, rest))

	// Application endpoints
	var apps []*model.Application
	dataPattern := regexp.MustCompile(`^data://`)
  // iterate apps and configure paths
  for _, mt := range l.Mountpoints {
		var dataScheme bool
		urlPath := mt.UrlPath
		path := mt.Path
		if dataPattern.MatchString(path) {
			dataScheme = true
			path = dataPattern.ReplaceAllString(path, "data/")
		}
    var p string
    p, err = filepath.Abs(path)
    if err != nil {
      return nil, err
    }
		name := filepath.Base(path)

		if urlPath == "" {
			urlPath = "/app/" + name + "/"
		}

    app := model.Application{}
		apps = append(apps, &app)

		var loader model.ApplicationLoader = model.FileSystemLoader{}

		// Load from bundled assets
		if dataScheme && !config.Dev {
			// TODO: implement asset loader logic
			loader = model.AssetLoader{}
		}

    // Load the application files into memory
		if err = app.Load(p, loader); err != nil {
			return nil, err
		}

    // Publish the application files to a build directory
    if err = app.Publish(nil); err != nil {
      return nil, err
    }

    // Serve the static build files from the mountpoint path.
    url := urlPath
    log.Printf("Serving app %s from %s", url, app.Public)
    mux.Handle(url, http.StripPrefix(url, http.FileServer(http.Dir(app.Public))))

		// Serve the raw source files.
    url = urlPath + "-/source/"
    log.Printf("Serving source %s from %s", url, p)
		sourceFileServer := http.StripPrefix(url, http.FileServer(http.Dir(p)))
		mux.HandleFunc(url, func(res http.ResponseWriter, req *http.Request) {
			//log.Println("got source req")
			//log.Printf("%#v\n", req)
			//log.Printf("%#v\n", req.URL)
			// TODO: serve in-memory versions
			sourceFileServer.ServeHTTP(res, req)
		})

		// TODO: serve app editor application from /editor

		/*
		if config.Dev {
			//mux.Handle("/", http.FileServer(http.Dir("data")))
		} else {
			mux.Handle("/",
				http.FileServer(
					&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "data"}))
		}
		*/
  }

  return apps, nil
}
