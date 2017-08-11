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
	UrlPath string
	// The path to pass to the loader.
	Path string	
}

type PageLoop struct {
	Server *http.Server
	// All application mountpoints.
  Mountpoints []Mountpoint
}

type ServerConfig struct {
	Addr string

	// List of user application mountpoints.
  Mountpoints []Mountpoint

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

	// System applications to mount.
	l.Mountpoints = append(l.Mountpoints, Mountpoint{UrlPath: "/", Path: "data://app/home"})

	// Add user applications.
	l.Mountpoints = append(l.Mountpoints, config.Mountpoints...)

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
      return err
    }
		name := filepath.Base(path)

		if urlPath == "" {
			urlPath = "/app/" + name + "/"
		}

    app := model.Application{}

		var loader model.ApplicationLoader = model.FileSystemLoader{}

		// Load from bundled assets
		if dataScheme && !config.Dev {
			loader = model.AssetLoader{}
		}

    // Load the application files into memory
		println("Load apps")
		println(p)
		if err = app.Load(p, loader); err != nil {
			return err
		}

		/*
		} else {
			// Load from the bundled binary assets (data://)
			println("load bundled data://")
		}
		*/

    // Publish the application files to a build directory
		println("Publish apps")
    if err = app.Publish(nil); err != nil {
      return err
    }

    // Serve the static build files.
    url := urlPath
    log.Printf("Serving app %s from %s", url, app.Public)
    mux.Handle(url, http.StripPrefix(url, http.FileServer(http.Dir(app.Public))))

		// Serve the raw source files.
    url = urlPath + "source/"
    log.Printf("Serving app source %s from %s", url, p)
		sourceFileServer := http.StripPrefix(url, http.FileServer(http.Dir(p)))
		mux.HandleFunc(url, func(res http.ResponseWriter, req *http.Request) {
			//log.Println("got editor request")
			//log.Printf("%#v\n", req)
			//log.Printf("%#v\n", req.URL)
			// TODO: serve in-memory versions
			sourceFileServer.ServeHTTP(res, req)
		})

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

  return nil
}
