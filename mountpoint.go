package pageloop

import (
  "os"
  "fmt"
  "path/filepath"
  . "github.com/tmpfs/pageloop/model"
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

