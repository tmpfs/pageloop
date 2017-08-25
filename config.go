package pageloop

// Represents a runtime configuration.
type ServerConfig struct {

  // Address for the web server to bind to.
	Addr string `json:"addr" yaml:"addr"`

	// List of system application mountpoints.
  Mountpoints []Mountpoint `json:"mountpoints" yaml:"mountpoints"`

	// Load system assets from the file system, don't use
	// the embedded assets.
	Dev bool `json:"dev" yaml:"dev"`
}


// A mountpoint maps a path location indicating the source
// for an application and a URL that the application should
// be mounted at.
type Mountpoint struct {
	// The URL path component.
  UrlPath string `json:"url" yaml:"url"`
	// The path to pass to the loader.
  Path string	`json:"path" yaml:"path"`
	// Description to pass to the application.
  Description string `json:"description" yaml:"description"`
}

