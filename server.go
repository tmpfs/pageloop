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

	sys := model.NewContainer("system", "System applications")
	usr := model.NewContainer("user", "User applications")
	tpl := model.NewContainer("template", "Application templates")
	snx := model.NewContainer("sandbox", "Playground")

	l.Host.Add(sys)
	l.Host.Add(usr)
	l.Host.Add(tpl)
	l.Host.Add(snx)

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
	system = append(system, Mountpoint{UrlPath: "/api/browser/", Path: "data://app/api/browser"})
	system = append(system, Mountpoint{UrlPath: "/api/docs/", Path: "data://app/api/docs"})
	system = append(system, Mountpoint{UrlPath: "/api/probe/", Path: "data://app/api/probe"})

  if err = l.LoadMountpoints(system, sys); err != nil {
    return nil, err
  }

	// Add user applications from configuration mountpoints.
  if err = l.LoadMountpoints(config.Mountpoints, usr); err != nil {
    return nil, err
  }

	l.MountContainer(sys)
	l.MountContainer(usr)

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

// Iterates a list of mountpoints and creates an application for each mountpoint
// and adds it to the given container.
func (l*PageLoop) LoadMountpoints(mountpoints []Mountpoint, container *model.Container) error {
  var err error
	// Application endpoints
	dataPattern := regexp.MustCompile(`^data://`)
  // iterate apps and configure paths
  for _, mt := range mountpoints {
		//var dataScheme bool
		urlPath := mt.UrlPath
		path := mt.Path
		if dataPattern.MatchString(path) {
			//dataScheme = true
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

		app := model.Application{Url: urlPath}

    // Load the application files into memory
		if err = app.Load(p, nil); err != nil {
			return err
		}

		// TODO: make publishing optional

    // Publish the application files to a build directory
    if err = app.Publish(nil); err != nil {
      return err
    }

		// Add to the container
		if err = container.Add(&app); err != nil {
			return err
		}
  }
	return nil
}

// Mount all applications in a container.
func (l *PageLoop) MountContainer(container *model.Container) {
	for _, a := range container.Apps {
		l.MountApplication(a)	
	}
}

// Mount an application from Public to Url.
func (l *PageLoop) MountApplication(app *model.Application) {

	// Serve the static build files from the mountpoint path.
	url := app.Url
	log.Printf("Serving app %s from %s", url, app.Public)
	mux.Handle(url, http.StripPrefix(url, http.FileServer(http.Dir(app.Public))))

	// Serve the raw source files.
	url = url + "-/source/"
	log.Printf("Serving src %s from %s", url, app.Path)
	sourceFileServer := http.StripPrefix(url, http.FileServer(http.Dir(app.Path)))
	mux.HandleFunc(url, func(res http.ResponseWriter, req *http.Request) {
		//log.Println("got source req")
		//log.Printf("%#v\n", req)
		//log.Printf("%#v\n", req.URL)
		// TODO: serve in-memory versions
		sourceFileServer.ServeHTTP(res, req)
	})
}

