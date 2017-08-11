package pageloop

import (
  "log"
	"errors"
  "net/http"
  "path/filepath"
  "time"
  "github.com/tmpfs/pageloop/model"
  "github.com/elazarl/go-bindata-assetfs"
)

var config ServerConfig
var mux *http.ServeMux

type PageLoop struct {
	Server *http.Server
}

type ServerConfig struct {
	Addr string
  AppPaths []string
	// Load system assets from the file system, don't use 
	// the embedded assets
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

  if err = l.LoadApps(config); err != nil {
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

func (l *PageLoop) LoadApps(config ServerConfig) error {
  var err error
  // iterate apps and configure paths
  for _, path := range config.AppPaths {
    var p string
    p, err = filepath.Abs(path)
    if err != nil {
      return err
    }
    name := filepath.Base(path)

    app := model.Application{}

    // Load the application files into memory
    if err = app.Load(p, nil); err != nil {
      return err
    }

    // Publish the application files to a build directory
    if err = app.Publish(nil); err != nil {
      return err
    }

    // Serve the static build files.
    url := "/apps/" + name + "/"
    log.Printf("Serving %s from %s", url, app.Public)
    mux.Handle(url, http.StripPrefix(url, http.FileServer(http.Dir(app.Public))))

		//mux.Handle("/api/", apiHandler{})
    url = "/editor/apps/" + name + "/"
    log.Printf("Serving %s from %s", url, p)
		sourceFileServer := http.FileServer(http.Dir(p))
		mux.HandleFunc(url, func(res http.ResponseWriter, req *http.Request) {
			log.Println("got editor request")
			log.Printf("%#v\n", req)
			log.Printf("%#v\n", req.URL)
			sourceFileServer.ServeHTTP(res, req)
		})

		if config.Dev {
			mux.Handle("/", http.FileServer(http.Dir("data")))
		} else {
			mux.Handle("/",
				http.FileServer(
					&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "data"}))
		}
  }

  return nil
}
