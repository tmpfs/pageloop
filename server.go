// Collaborative realtime web based document manager.
package pageloop

import (
	"fmt"
  "log"
	"mime"
  "net/http"
  "time"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/handler"
  . "github.com/tmpfs/pageloop/model"
)

const(
	HTML_MIME = "text/html; charset=utf-8"
)

var Name string = "pageloop"
var Version string = "1.0"

// Primary serve mux handler for built in endpoints.
var mux *http.ServeMux

var(
  manager *MountpointManager
)

type PageLoop struct {
  // Server configuration
  Config *ServerConfig

	// Underlying HTTP server.
	Server *http.Server `json:"-"`

  // Reference to the mountpoint manager
  MountpointManager *MountpointManager

	// Application host
	Host *Host
}

// Creates an HTTP server.
func (l *PageLoop) NewServer(config *ServerConfig) (*http.Server, error) {
  var err error

  l.Config = config

  // Set up a host for our containers
	l.Host = NewHost()

  l.MountpointManager = NewMountpointManager(l.Config, l.Host)

  // Initialize the command adapter
  adapter := &CommandAdapter{Host: l.Host, Mountpoints: l.MountpointManager}

	// Configure application containers.
	sys := NewContainer("system", "System applications.", true)
	tpl := NewContainer("template", "Application & document templates.", true)
	usr := NewContainer("user", "User applications.", false)

	l.Host.Add(sys)
	l.Host.Add(tpl)
	l.Host.Add(usr)

  // Initialize server mux
  mux = http.NewServeMux()
  var handler http.Handler

	// RPC global endpoint (/rpc/)
	handler = RpcService(mux, l.Host)
	l.MountpointManager.MountpointMap[RPC_URL] = handler
	log.Printf("Serving rpc service from %s", RPC_URL)

	// REST API global endpoint (/api/)
	handler = RestService(mux, adapter)
	l.MountpointManager.MountpointMap[API_URL] = handler
	log.Printf("Serving rest service from %s", API_URL)

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
	l.MountContainer(sys)
	l.MountContainer(tpl)
	l.MountContainer(usr)

  s := &http.Server{
    Addr:           config.Addr,
    Handler:        ServerHandler{MountpointManager: l.MountpointManager, Mux: mux},
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  return s, nil
}

// Mount all applications in a container.
func (l *PageLoop) MountContainer(container *Container) {
	for _, a := range container.Apps {
		MountApplication(l.MountpointManager.MountpointMap, l.Host, a)
	}
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
