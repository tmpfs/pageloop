package core

import (
  "os"
  "path/filepath"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

const(
	API_URL = "/api/"
	RPC_URL = "/rpc/"
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

  // User configuration merged with this config, only
  // available if merge has been called.
  userConfig *ServerConfig

  // Path used when calling merge to load a user configuration.
  userConfigPath string
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
  return conf
}

// Attempt to delete a user mountpoint for the given URL.
func (c *ServerConfig) DeleteMountpoint(url string) *ServerConfig {
  var conf *ServerConfig = c.UserConfig()
  for i, m := range conf.Mountpoints {
    if url == m.Url {
      before := conf.Mountpoints[0:i]
      after := conf.Mountpoints[i+1:]
      conf.Mountpoints = append(before, after...)
    }
  }
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

  for _, m := range tempServerConfig.Mountpoints {
    // Force user supplied applications into particular container
    m.Container = "user"
    //c.Mountpoints = append(c.Mountpoints, m)
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
  // TODO: preserve permissions
  if err = ioutil.WriteFile(path, content, 0644); err != nil {
    return err
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
