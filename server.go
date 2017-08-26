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

var config ServerConfig
var mux *http.ServeMux

// Maps application URLs to HTTP handlers.
//
// Because we want to mount and unmount applications and we cannot remove
// a handler we have a single handler that defers to these handlers.
var mountpoints map[string] http.Handler

// We need to know which requests go through the normal serve mux logic
// so they do not collide with application requests.
var multiplex map[string] bool

type PageLoop struct {
	// Underlying HTTP server.
	Server *http.Server `json:"-"`

	// Application host
	Host *model.Host
}

// Temporary map used when initializing loaded mountpoint definitions
// containing a container reference which was declared by string name
// in the mountpoint definition.
type MountpointMap struct {
  Container *model.Container
  Mountpoints []Mountpoint
}

// Main HTTP server handler.
type ServerHandler struct {}

// The default server handler, defers to a multiplexer.
func (h ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var handler http.Handler
	var path string = req.URL.Path

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

// Tests the request path and attempts to find a corresponding source file
// in the application files.
func (h ApplicationSourceHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	index := "index.html"
	urls := h.App.Urls
	path := "/" + req.URL.Path
	clean := strings.TrimSuffix(path, "/")
	indexPage := clean + "/" + index

	// TODO: handle HEAD requests

	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var file *model.File

	// Exact match
	if urls[path] != nil {
		file = urls[path]
	// Normalized without a trailing slash
	} else if(urls[clean] != nil) {
		file = urls[clean]
	// Check for index page
	} else if(urls[indexPage] != nil) {
		file = urls[indexPage]
	}

	// TODO: write cache busting headers
	// TODO: handle directory requests (no data)
	if file != nil && !file.Info().IsDir() {
		ext := filepath.Ext(file.Name)
		ct := mime.TypeByExtension(ext)
		output := file.Source(h.Raw)
		res.Header().Set("Content-Type", ct)
		res.Header().Set("Content-Length", strconv.Itoa(len(output)))
		res.Write(output)
		return
	}

	http.NotFound(res, req)
}

// Creates an HTTP server.
func (l *PageLoop) NewServer(config *ServerConfig) (*http.Server, error) {
  var err error

  // Set up a host for our containers
	l.Host = model.NewHost()

	// Configure application containers.
	sys := model.NewContainer("system", "System applications.", true)
	tpl := model.NewContainer("template", "Application & document templates.", true)
	usr := model.NewContainer("user", "User applications.", false)
	snx := model.NewContainer("sandbox", "Playground.", false)

	l.Host.Add(sys)
	l.Host.Add(tpl)
	l.Host.Add(usr)
	l.Host.Add(snx)

	// Initialize mountpoint maps
	mountpoints = make(map[string] http.Handler)
	multiplex = make(map[string] bool)

  // Initialize server mux
  mux = http.NewServeMux()

	// RPC global endpoint (/rpc/)
	rpc := NewRpcService(l, mux)
	log.Printf("Serving rpc service from %s", RPC_URL)

	// REST API global endpoint (/api/)
	rest := NewRestService(l, mux)
	log.Printf("Serving rest service from %s", API_URL)

	multiplex[rpc.Url] = true
	multiplex[rest.Url] = true

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
    if err = l.LoadMountpoints(c.Mountpoints, c.Container); err != nil {
      return nil, err
    }
  }

  // Mount containers and the applications within them
	l.MountContainer(sys)
	l.MountContainer(tpl)
	l.MountContainer(usr)
	l.MountContainer(snx)

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
	// Bundled application endpoints
	dataPattern := regexp.MustCompile(`^data://`)

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
      return err
    }
		name := filepath.Base(path)

		// No mountpoint URL given so we assume an app
		// relative to the container
		if urlPath == "" {
			urlPath = fmt.Sprintf("/%s/%s/", container.Name, name)
		}

		app := model.NewApplication(urlPath, mt.Description)
		fs := model.NewUrlFileSystem(app)
		app.FileSystem = fs

    // Load the application files into memory
		if err = app.Load(p); err != nil {
			return err
		}

		// TODO: make publishing optional

    // Publish the application files to a build directory
    if err = app.Publish("public/" + container.Name); err != nil {
      return err
    }

		// Add to the container
		if err = container.Add(app); err != nil {
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

// Mount an application such that it's published and source
// files are accessible over HTTP. This serves the published files
// as static files and serves two versions of the source file
// from in memory data. The src version is the file with any frontmatter
// stripped and the raw version includes frontmatter.
func (l *PageLoop) MountApplication(app *model.Application) {
	// Serve the static build files from the mountpoint path.
	url := app.Url
	log.Printf("Serving app %s from %s", url, app.Public)
	mountpoints[url] = http.StripPrefix(url, http.FileServer(http.Dir(app.Public)))

	// Serve the source files with frontmatter stripped.
	url = "/apps/source/" + app.Container.Name + "/" + app.Name + "/"
	log.Printf("Serving src %s from %s", url, app.Path)
	mountpoints[url] = http.StripPrefix(url, ApplicationSourceHandler{App: app})

	// Serve the raw source files.
	url = "/apps/raw/" + app.Container.Name + "/" + app.Name + "/"
	log.Printf("Serving raw %s from %s", url, app.Path)
	mountpoints[url] = http.StripPrefix(url, ApplicationSourceHandler{App: app, Raw: true})
}

func init() {
  // Mime types set to those for code mirror modes
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".yml", "text/x-yaml")
	mime.AddExtensionType(".yaml", "text/x-yaml")
	mime.AddExtensionType(".md", "text/x-markdown")
	mime.AddExtensionType(".markdown", "text/x-markdown")
}
