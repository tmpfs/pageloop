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
	//"github.com/gorilla/rpc"
	//"github.com/gorilla/rpc/json"
  //"github.com/elazarl/go-bindata-assetfs"
)

const(
	HTML_MIME = "text/html; charset=utf-8"
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

	// Application host
	Host *model.Host
	//Container *model.Container
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

// Creates an HTTP server.
func (l *PageLoop) NewServer(config ServerConfig) (*http.Server, error) {
  var err error

	l.Host = model.NewHost()

	// Configure application container.
	//l.Container = model.NewContainer()

	sys := model.NewContainer("System applications", "")
	usr := model.NewContainer("User applications", "")

	l.Host.Add("system", sys)
	l.Host.Add("user", usr)

  // Initialize server mux
  mux = http.NewServeMux()

	// RPC global endpoint (/rpc/)
	NewRpcService(l, mux)
	log.Printf("Serving rpc service from %s", RPC_URL)

	// REST API global endpoint (/api/)
	NewRestService(l, mux)
	log.Printf("Serving rest service from %s", API_URL)

	// System applications to mount.
	var system []Mountpoint
	system = append(system, Mountpoint{UrlPath: "/", Path: "data://app/home"})
  if err = l.loadApps(system, sys); err != nil {
    return nil, err
  }

	// Add user applications from configuration mountpoints.
  if err = l.loadApps(config.Mountpoints, usr); err != nil {
    return nil, err
  }

  s := &http.Server{
    Addr:           config.Addr,
    Handler:        ServerHandler{},
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  return s, nil
}

// Start the HTTP server listening.
func (l *PageLoop) Listen(server *http.Server) error {
	var err error
	if server == nil {
		return errors.New("Cannot listen without a server, call NewServer().")
	}

	log.Printf("Listen %s", server.Addr)

  if err = server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Load application mountpoints.
func (l *PageLoop) loadApps(mountpoints []Mountpoint, container *model.Container) error {
  var err error

	// Application endpoints
	dataPattern := regexp.MustCompile(`^data://`)
  // iterate apps and configure paths
  for _, mt := range mountpoints {
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
			// TODO: implement asset loader logic
			//loader = model.AssetLoader{}
		}

    // Load the application files into memory
		if err = app.Load(p, loader); err != nil {
			return err
		}

    // Publish the application files to a build directory
    if err = app.Publish(nil); err != nil {
      return err
    }

		// Add to the container
		if err = container.Add(&app); err != nil {
			return err
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

  return nil
}
