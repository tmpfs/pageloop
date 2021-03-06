// Collaborative realtime web based document manager.
package pageloop

import (
	"fmt"
  "log"
	"mime"
  "net/http"
  "time"
  //. "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/handler"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/service"
  . "github.com/tmpfs/pageloop/rpc"
)

type PageLoop struct {
  // Server configuration
  Config *ServerConfig

	// Underlying HTTP server
	Server *http.Server `json:"-"`

  // Server multiplexer
  Mux *http.ServeMux `json:"-"`

  // Reference to the mountpoint manager
  MountpointManager *MountpointManager `json:"-"`

	// Application host
	Host *Host

  // Map of services
  Services *ServiceMap
}

// Creates an HTTP server.
func (l *PageLoop) NewServer(config *ServerConfig) (*http.Server, error) {
  var handler http.Handler

  // Configuration for the server
  l.Config = config

  // Initialize server multiplexer
  l.Mux = http.NewServeMux()

  // Set up a host for our containers
	l.Host = NewHost()

  // Manager for application mountpoints.
  //
  // Application mountpoints are dynamic (they can be added and removed at runtime)
  // so they need special care.
  l.MountpointManager = NewMountpointManager(l.Config, l.Host)

  l.initServices()

	// Configure application containers.
	sys := NewContainer("system", "System applications.", true)
	tpl := NewContainer("template", "Application & document templates.", true)
	usr := NewContainer("user", "User applications.", false)
	l.Host.Add(usr)
	l.Host.Add(sys)
	l.Host.Add(tpl)

	// REST API global endpoint (/api/)
	handler = RestService(l.Mux, l.Services, l.Host, l.MountpointManager)
	l.MountpointManager.MountpointMap[API_URL] = handler
	log.Printf("Serving rest service from %s", API_URL)

	// Websocket global endpoint (/ws/)
	handler = WebsocketService(l.Mux, l.Services, l.Host, l.MountpointManager)
	l.MountpointManager.MountpointMap[WEBSOCKET_URL] = handler
	log.Printf("Serving websocket service from %s", WEBSOCKET_URL)

	// RPC global endpoint (/rpc/)
  /*
	handler = RpcService(l.Mux, l.Host)
	l.MountpointManager.MountpointMap[RPC_URL] = handler
	log.Printf("Serving rpc service from %s", RPC_URL)
  */

  // Collect mountpoints by container name
  if collection, err := l.MountpointManager.Collect(config.Mountpoints, config.UserConfig().Mountpoints); err != nil {
    return nil, err
  // Load the mountpoints using the container map
  } else {
    // Discarding the returned list of applications
    if _, err := l.MountpointManager.LoadCollection(collection); err != nil {
      return nil, err
    }
  }

  // Mount containers and the applications within them
	l.MountContainer(sys)
	l.MountContainer(tpl)
	l.MountContainer(usr)

  s := &http.Server{
    Addr:           config.Addr,
    Handler:        ServerHandler{MountpointManager: l.MountpointManager, Mux: l.Mux},
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

// Initialize services
func (l *PageLoop) initServices() {
  l.Services = &ServiceMap{}

  srv := new(RpcServices)
  core := new(CoreService)
  host := new(HostService)
  ctx := new(ContainerService)
  app := new(AppService)
  zip := new(ArchiveService)
  file := new(FileService)
  job := new(JobService)
  tpl := new(TemplateService)

  srv.Services = l.Services
  srv.Router = DefaultRouter

  host.Host = l.Host
  ctx.Host = l.Host
  app.Host = l.Host
  zip.Host = l.Host
  file.Host = l.Host
  tpl.Host = l.Host

  ctx.Mountpoints = l.MountpointManager
  app.Mountpoints = l.MountpointManager

  l.Services.MustRegister(core, "Core")
  l.Services.MustRegister(host, "Host")
  l.Services.MustRegister(ctx, "Container")
  l.Services.MustRegister(app, "Application")
  l.Services.MustRegister(zip, "Archive")
  l.Services.MustRegister(file, "File")
  l.Services.MustRegister(job, "Job")
  l.Services.MustRegister(tpl, "Template")
  l.Services.MustRegister(srv, "Service")
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
