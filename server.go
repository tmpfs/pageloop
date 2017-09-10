// Collaborative realtime web based document manager.
package pageloop

import (
	"fmt"
  "log"
	"mime"
	"strings"
	"strconv"
  "net/http"
  "path/filepath"
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

func send (res http.ResponseWriter, req *http.Request, file *model.File, output []byte) {
	path := "/" + req.URL.Path
  base := filepath.Base(path)

  ext := filepath.Ext(file.Name)
  ct := mime.TypeByExtension(ext)

  // TODO: remove this?
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
  adapter = &CommandAdapter{Host: l.Host, Mountpoints: manager}

	// Configure application containers.
	sys := model.NewContainer("system", "System applications.", true)
	tpl := model.NewContainer("template", "Application & document templates.", true)
	usr := model.NewContainer("user", "User applications.", false)

	l.Host.Add(sys)
	l.Host.Add(tpl)
	l.Host.Add(usr)

  // Initialize server mux
  mux = http.NewServeMux()

	// RPC global endpoint (/rpc/)
	rpc := NewRpcService(mux, l.Host)
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
    if _, err = manager.LoadMountpoints(c.Mountpoints, c.Container); err != nil {
      return nil, err
    }
  }

  // Mount containers and the applications within them
	manager.MountContainer(sys)
	manager.MountContainer(tpl)
	manager.MountContainer(usr)

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
		return fmt.Errorf("Cannot listen without a server, call NewServer().")
	}

	log.Printf("Listen %s", server.Addr)

  if err = server.ListenAndServe(); err != nil {
		return err
	}

	return nil
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
