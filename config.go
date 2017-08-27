package pageloop

import (
  "os"
  "path/filepath"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

var defaultServerConfig *ServerConfig

// Represents a runtime configuration.
type ServerConfig struct {

  // Address for the web server to bind to.
	Addr string `json:"addr,omitempty" yaml:"addr,omitempty"`

	// List of application mountpoints.
  Mountpoints []Mountpoint `json:"mountpoints" yaml:"mountpoints"`

  // Directory for generated source files
  SourceDirectory string `json:"source,omitempty" yaml:"source,omitempty"`

  // Directory for build publish preview
  PublishDirectory string `json:"publish,omitempty" yaml:"publish,omitempty"`

	// Load system assets from the file system, don't use
	// the embedded assets.
	Dev bool `json:"dev,omitempty" yaml:"dev,omitempty"`

  // User configuration merged with this config, only
  // available if merge has been called.
  userConfig *ServerConfig

  // Path used when calling merge to load a user configuration.
  userConfigPath string
}


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
}

// Public access to the default server config.
func DefaultServerConfig() *ServerConfig {
  return defaultServerConfig
}

// Gets a user config object if one was assigned during
// a merge operation otherwise creates an empty configuration
// and assigns it as the user configuration.
func (c *ServerConfig) UserConfig() *ServerConfig {
  if c.userConfig == nil {
    c.userConfig = &ServerConfig{}
  }
  return c.userConfig
}

// Add a mountpoint to the list of user configuration mountpoints
// and returns the user configuration.
func (c *ServerConfig) AddMountpoint(m Mountpoint) *ServerConfig {
  var conf *ServerConfig = c.UserConfig()
  // Append to the user configuration mountpoints
  conf.Mountpoints = append(conf.Mountpoints, m)
  // Append to the primary list for easier iteration
  c.Mountpoints = append(c.Mountpoints, m)
  return conf
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

  c.userConfig = tempServerConfig
  c.userConfigPath = path

  return nil
}

// Write a configuration to disc as YAML.
//
// When no path is given and merge has been called
// the file is written to the path that was used when
// loading a user configuration. Otherwise when there is no
// path and no loaded user configuration it writes
// to config.yml in the current working directory.
func (c *ServerConfig) WriteFile(conf *ServerConfig, path string) error {
  var err error
  var wd string
  if path == "" {
    if  c.userConfigPath == "" {
      if wd, err = os.Getwd(); err != nil {
        return err
      }
      path = filepath.Join(wd, "config.yml")
    } else {
      path = c.userConfigPath
    }
  }
  var content []byte
  if content, err = yaml.Marshal(conf); err != nil {
    return err
  }
  println("write config file: " + path)
  println("write config file: " + string(content))
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
