package pageloop

import (
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

var defaultServerConfig *ServerConfig

// Represents a runtime configuration.
type ServerConfig struct {
  // Address for the web server to bind to.
	Addr string `json:"addr" yaml:"addr"`

	// List of application mountpoints.
  Mountpoints []Mountpoint `json:"mountpoints" yaml:"mountpoints"`

  // Directory for build publish preview
  PublishDirectory string `json:"publish" yaml:"publish"`

	// Load system assets from the file system, don't use
	// the embedded assets.
	Dev bool `json:"dev" yaml:"dev"`
}


// A mountpoint maps a path location indicating the source
// for an application and a URL that the application should
// be mounted at.
type Mountpoint struct {
  // Name of the parent container for the application.
  Container string `json:"container" yaml:"container"`
	// The URL location for the application mountpoint.
  Url string `json:"url" yaml:"url"`
	// The path to pass to the loader.
  Path string	`json:"path" yaml:"path"`
	// Description to pass to the application.
  Description string `json:"description" yaml:"description"`
}

// Public access to the default server config.
func DefaultServerConfig() *ServerConfig {
  return defaultServerConfig
}

// Load and merge a user supplied configuration file.
//
// Mountpoints are appended to the defaults and each mountpoint
// in the user configuration is added to the user container.
//
// User supplied configurations can currently only specify Addr
// and Mountpoints.
func (c *ServerConfig) Merge(path string) error {
  var err error
  var content []byte

  if content, err = ioutil.ReadFile(path); err != nil {
    return err
  }

  tempServerConfig := &ServerConfig{}
  if err := yaml.Unmarshal(content, tempServerConfig); err != nil {
    return err
  }

  if tempServerConfig.Addr != "" {
    c.Addr = tempServerConfig.Addr
  }

  if tempServerConfig.PublishDirectory != "" {
    c.PublishDirectory = tempServerConfig.PublishDirectory
  }

  for _, m := range tempServerConfig.Mountpoints {
    // Force user supplied applications into particular container
    m.Container = "user"
    c.Mountpoints = append(c.Mountpoints, m)
  }

  return nil
}

// Load system server configuration.
func init() {
  serverConfigFile := MustAsset("config.yml")
  defaultServerConfig = &ServerConfig{Addr: ":3577"}
  if err := yaml.Unmarshal(serverConfigFile, defaultServerConfig); err != nil {
    panic(err)
  }
}
