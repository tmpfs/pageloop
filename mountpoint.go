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
  mountpoints[url] = http.StripPrefix(url, ApplicationPublicHandler{App: app, FileServer: fileserver})

	// Serve the source files with frontmatter stripped.
	url = app.SourceUrl()
	log.Printf("Serving src %s from %s", url, app.SourceDirectory())
  mountpoints[url] = http.StripPrefix(url, ApplicationSourceHandler{App: app})

	// Serve the raw source files.
	url = app.RawUrl()
	log.Printf("Serving raw %s from %s", url, app.SourceDirectory())
  mountpoints[url] = http.StripPrefix(url, ApplicationSourceHandler{App: app, Raw: true})
}

// Unmount an application from the web server.
func (m *MountpointManager) UnmountApplication(app *Application) {
  delete(mountpoints, app.PublishUrl())
  delete(mountpoints, app.SourceUrl())
  delete(mountpoints, app.RawUrl())
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
