// Collaborative realtime web based document manager.
package pageloop

import (
	"fmt"
  "log"
	"errors"
	"mime"
	"strings"
	"strconv"
  "net/http"
  "path/filepath"
	"regexp"
  "time"
  "github.com/tmpfs/pageloop/model"
)

const(
	HTML_MIME = "text/html; charset=utf-8"
)

var Name string = "pageloop"
var Version string = "1.0"

// Primary serve mux handler for built in endpoints.
var mux *http.ServeMux

var(
  adapter *CommandAdapter
  manager *MountpointManager
  listing *DirectoryList
)

type PageLoop struct {
  // Server configuration
  Config *ServerConfig

	// Underlying HTTP server.
	Server *http.Server `json:"-"`

	// Application host
	Host *model.Host
}

// Main HTTP server handler.
type ServerHandler struct {}

// The default server handler, defers to a multiplexer.
func (h ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var handler http.Handler
	var path string = req.URL.Path

  res.Header().Set("Access-Control-Allow-Origin", "*")

	// Look for serve mux mappings first
	for k, _ := range multiplex {
		if strings.HasPrefix(path, k) {
			handler, _ = mux.Handler(req)
			handler.ServeHTTP(res, req)
			return
		}
	}

	// Check for application mountpoints.
	//
	// Serve the highest score which is the longest
	// matching URL path.
	var score int
	for k, v := range mountpoints {
		if strings.HasPrefix(path, k) {
			if handler != nil && len(k) < score {
				continue
			}
			handler = v
			score = len(k)
		}
	}

	if handler == nil {
		handler = http.NotFoundHandler()
	}
	handler.ServeHTTP(res, req)
}

// Serves application source files from memory.
type ApplicationSourceHandler struct {
	App *model.Application
	Raw bool
}

// Serves application public files from disc.
type ApplicationPublicHandler struct {
	App *model.Application
  FileServer http.Handler
}

func (h ApplicationPublicHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  app := h.App
  path := "/" + req.URL.Path
  file := app.Urls[path]
	clean := strings.TrimSuffix(path, "/")
  // FIXME: this is rubbish
	indexPage := clean + "/index.html"
	indexMdPage := clean + "/index.md"
  if file != nil && file.Directory && app.Urls[indexPage] == nil && app.Urls[indexMdPage] == nil {
    listing.List(file, res, req)
    return
  }
  // Defer to file server for files
  h.FileServer.ServeHTTP(res, req)
}

func send (res http.ResponseWriter, req *http.Request, file *model.File, output []byte) {
	path := "/" + req.URL.Path
  base := filepath.Base(path)

  ext := filepath.Ext(file.Name)
  ct := mime.TypeByExtension(ext)
  if (ext == ".pdf") {
    res.Header().Set("Content-Disposition", "inline; filename=" + base)
  }
  res.Header().Set("Content-Type", ct)
  res.Header().Set("Content-Length", strconv.Itoa(len(output)))
  if (req.Method == http.MethodHead) {
    return
  }
  res.Write(output)
}

// Tests the request path and attempts to find a corresponding source file
// in the application files.
func (h ApplicationSourceHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	index := "index.html"
	urls := h.App.Urls
	path := "/" + req.URL.Path
	clean := strings.TrimSuffix(path, "/")
  trailing := clean + "/"
	indexPage := clean + "/" + index

	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var file *model.File

	// Exact match
	if urls[path] != nil {
		file = urls[path]
  } else if urls[trailing] != nil && !strings.HasSuffix(path, "/") {
    redirect := http.RedirectHandler(req.URL.Path + "/", http.StatusMovedPermanently)
    redirect.ServeHTTP(res, req)
    return
	// Normalized without a trailing slash
	} else if urls[clean] != nil {
		file = urls[clean]
	// Check for index page
	} else if urls[indexPage] != nil {
		file = urls[indexPage]
	}

	// TODO: write cache busting headers
	// TODO: handle directory requests (no data)
	if file != nil && !file.Info().IsDir() {
		output := file.Source(h.Raw)
    send(res, req, file, output)
		return
  // Handle directory listing
	} else if file != nil {
    listing.List(file, res, req)
    return
  }

	http.NotFound(res, req)
}

// Creates an HTTP server.
func (l *PageLoop) NewServer(config *ServerConfig) (*http.Server, error) {
  var err error

  // Initialize the command adapter
  l.Config = config

  // Set up a host for our containers
	l.Host = model.NewHost()

  manager = NewMountpointManager(config)

  listing = &DirectoryList{Host: l.Host}

  // TODO: remove Root reference
  adapter = &CommandAdapter{Root: l, Host: l.Host, Mountpoints: manager}

	// Configure application containers.
	sys := model.NewContainer("system", "System applications.", true)
	tpl := model.NewContainer("template", "Application & document templates.", true)
	usr := model.NewContainer("user", "User applications.", false)

	l.Host.Add(sys)
	l.Host.Add(tpl)
	l.Host.Add(usr)

	// Initialize mountpoint maps
	mountpoints = make(map[string] http.Handler)
	multiplex = make(map[string] bool)

  // Initialize server mux
  mux = http.NewServeMux()

	// RPC global endpoint (/rpc/)
	rpc := NewRpcService(l, mux)
	log.Printf("Serving rpc service from %s", RPC_URL)

	// REST API global endpoint (/api/)
	rest := NewRestService(mux)
	log.Printf("Serving rest service from %s", API_URL)

	multiplex[strings.TrimSuffix(rpc.Url, "/")] = true
	multiplex[strings.TrimSuffix(rest.Url, "/")] = true

  // Collect mountpoints by container name
  var collection map[string] *MountpointMap = make(map[string] *MountpointMap)
  for _, m := range config.Mountpoints {
    c := l.Host.GetByName(m.Container)
    if c == nil {
      return nil, fmt.Errorf("Unknown container %s", m.Container)
    }
    if collection[m.Container] == nil {
      collection[m.Container] = &MountpointMap{Container: c}
    }
    collection[m.Container].Mountpoints = append(collection[m.Container].Mountpoints, m)
  }

  // Load mountpoints
  for _, c := range collection {
    if _, err = l.LoadMountpoints(c.Mountpoints, c.Container); err != nil {
      return nil, err
    }
  }

  // Mount containers and the applications within them
	l.MountContainer(sys)
	l.MountContainer(tpl)
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

// Load a single mountpoint.
func (l *PageLoop) LoadMountpoint(mountpoint Mountpoint, container *model.Container) (*model.Application, error) {
  var err error
  var apps []*model.Application
  var list []Mountpoint
  list = append(list, mountpoint)
  if apps, err = l.LoadMountpoints(list, container); err != nil {
    return nil, err
  }
  return apps[0], nil
}

// Iterates a list of mountpoints and creates an application for each mountpoint
// and adds it to the given container.
func (l *PageLoop) LoadMountpoints(mountpoints []Mountpoint, container *model.Container) ([]*model.Application, error) {
  var err error
	// Bundled application endpoints
	dataPattern := regexp.MustCompile(`^data://`)

  var apps []*model.Application

  // iterate apps and configure paths
  for _, mt := range mountpoints {
		urlPath := mt.Url
		path := mt.Path
		if dataPattern.MatchString(path) {
			path = dataPattern.ReplaceAllString(path, "data/")
		}
    var p string
    p, err = filepath.Abs(path)
    if err != nil {
      return nil, err
    }
		name := filepath.Base(path)

		// No mountpoint URL given so we assume an app
		// relative to the container
		if urlPath == "" {
			urlPath = fmt.Sprintf("/%s/%s/", container.Name, name)
		}

		app := model.NewApplication(urlPath, mt.Description)
    app.IsTemplate = mt.Template
		fs := model.NewUrlFileSystem(app)
		app.FileSystem = fs

    // Load the application files into memory
		if err = app.Load(p); err != nil {
			return nil, err
		}

    // Only publish if the build file has not explicitly
    // enabled build at boot time
    var shouldPublish = true
    if app.HasBuilder() && !app.Builder.Boot {
      shouldPublish = false
    }

    if shouldPublish {
      // Publish the application files to a build directory
      if err = app.Publish(app.PublicDirectory()); err != nil {
        return nil, err
      }
    }

		// Add to the container
		if err = container.Add(app); err != nil {
			return nil, err
		}

    apps = append(apps, app)
  }
	return apps, nil
}

// Mount all applications in a container.
func (l *PageLoop) MountContainer(container *model.Container) {
	for _, a := range container.Apps {
		manager.MountApplication(a)
	}
}

func init() {
  // Mime types set to those for code mirror modes
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".babelrc", "application/json")
	mime.AddExtensionType(".yml", "text/x-yaml")
	mime.AddExtensionType(".yaml", "text/x-yaml")
	mime.AddExtensionType(".md", "text/x-markdown")
	mime.AddExtensionType(".markdown", "text/x-markdown")
}
