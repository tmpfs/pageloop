package pageloop

import (
  "os"
  "fmt"
  "log"
  "strings"
  "net/http"
  "path/filepath"
  . "github.com/tmpfs/pageloop/model"
)

var(
  // Maps application URLs to HTTP handlers.
  //
  // Because we want to mount and unmount applications and we cannot remove
  // a handler we have a single handler that defers to these handlers.
  mountpoints map[string] http.Handler

  // We need to know which requests go through the normal serve mux logic
  // so they do not collide with application requests.
  multiplex map[string] bool
)

// A mountpoint maps a path location indicating the source
// for an application and a URL that the application should
// be mounted at.
type Mountpoint struct {
  // Name of the parent container for the application.
  Container string `json:"container,omitempty" yaml:"container,omitempty"`
	// The URL location for the application mountpoint.
  Url string `json:"url" yaml:"url"`
	// The path to pass to the loader.
  Path string	`json:"path" yaml:"path"`
	// Description to pass to the application.
  Description string `json:"description" yaml:"description"`
  // Mark as a template
  Template bool `json:"template" yaml:"template"`
}

// Temporary map used when initializing loaded mountpoint definitions
// containing a container reference which was declared by string name
// in the mountpoint definition.
type MountpointMap struct {
  Container *Container
  Mountpoints []Mountpoint
}

type MountpointManager struct {
  Config *ServerConfig
}

func NewMountpointManager(c *ServerConfig) *MountpointManager {
	// Initialize mountpoint maps
	mountpoints = make(map[string] http.Handler)
	multiplex = make(map[string] bool)
  return &MountpointManager{Config: c}
}

// Delete a mountpoint for a userspace application and persist the list of mountpoints.
func (m *MountpointManager) DeleteApplicationMountpoint(url string) error {
  var conf *ServerConfig = m.Config.DeleteMountpoint(url)
  if err := m.Config.WriteFile(conf, ""); err != nil {
    return err
  }
  return nil
}

// Create and persist a mountpoint for a userspace application.
func (m *MountpointManager) CreateMountpoint(a *Application) (*Mountpoint, error) {
  var err error
  if a.Name == "" {
    return nil, fmt.Errorf("Cannot create a mountpoint without an application name")
  }

  if !ValidName(a.Name) {
    return nil, fmt.Errorf(
      "Application name %s is invalid, must match pattern %s", a.Name, NamePattern)
  }

  // Configure filesystem path for source files
  a.SetPath(filepath.Join(m.Config.SourceDirectory, a.Name))

  // Create source application directory
  if err := os.MkdirAll(a.SourceDirectory(), os.ModeDir | 0755); err != nil {
		return nil, err
	}

  var mt *Mountpoint = &Mountpoint{Path: a.Path, Url: a.Url, Description: a.Description}
  var conf *ServerConfig = m.Config.AddMountpoint(*mt)
  if err = m.Config.WriteFile(conf, ""); err != nil {
    return nil, err
  }
  return mt, nil
}

// Test if a mountpoint exists by URL.
func (m *MountpointManager) HasMountpoint(url string) bool {
  umu := strings.TrimSuffix(url, "/")
  if _, ok := multiplex[url]; ok {
    return true
  }
  if _, ok := multiplex[umu]; ok {
    return true
  }
  for _, m := range m.Config.Mountpoints {
    cmu := strings.TrimSuffix(m.Url, "/")
    if m.Url == url || cmu == umu {
      return true
    }
  }
  return false
}

// Mount all applications in a container.
func (m *MountpointManager) MountContainer(container *Container) {
	for _, a := range container.Apps {
		m.MountApplication(a)
	}
}

// Load a single mountpoint.
func (m *MountpointManager) LoadMountpoint(mountpoint Mountpoint, container *Container) (*Application, error) {
  var err error
  var apps []*Application
  var list []Mountpoint
  list = append(list, mountpoint)
  if apps, err = m.LoadMountpoints(list, container); err != nil {
    return nil, err
  }
  return apps[0], nil
}

// Iterates a list of mountpoints and creates an application for each mountpoint
// and adds it to the given container.
func (m *MountpointManager) LoadMountpoints(mountpoints []Mountpoint, container *Container) ([]*Application, error) {
  var err error
  var apps []*Application

  // iterate apps and configure paths
  for _, mt := range mountpoints {
		urlPath := mt.Url
		path := mt.Path
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

		app := NewApplication(urlPath, mt.Description)
    app.IsTemplate = mt.Template
		fs := NewUrlFileSystem(app)
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

// Mount an application such that it's published and source
// files are accessible over HTTP. This serves the published files
// as static files and serves two versions of the source file
// from in memory data. The src version is the file with any frontmatter
// stripped and the raw version includes frontmatter.
func (m *MountpointManager) MountApplication(app *Application) {
	// Serve the static build files from the mountpoint path.
	url := app.PublishUrl()
	log.Printf("Serving app %s from %s", url, app.PublicDirectory())
  fileserver := http.FileServer(http.Dir(app.PublicDirectory()))
  mountpoints[url] = http.StripPrefix(url, PublicHandler{App: app, FileServer: fileserver})

	// Serve the source files with frontmatter stripped.
	url = app.SourceUrl()
	log.Printf("Serving src %s from %s", url, app.SourceDirectory())
  mountpoints[url] = http.StripPrefix(url, SourceHandler{App: app})

	// Serve the raw source files.
	url = app.RawUrl()
	log.Printf("Serving raw %s from %s", url, app.SourceDirectory())
  mountpoints[url] = http.StripPrefix(url, SourceHandler{App: app, Raw: true})
}

// Unmount an application from the web server.
func (m *MountpointManager) UnmountApplication(app *Application) {
  delete(mountpoints, app.PublishUrl())
  delete(mountpoints, app.SourceUrl())
  delete(mountpoints, app.RawUrl())
}

