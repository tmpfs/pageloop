package model

import(
  "fmt"
  "strings"
  "net/url"
  "net/http"
  . "github.com/tmpfs/pageloop/util"
)

type Reference interface {
  ParseUrl(uri string) (ref Reference, err error)
}

// Represents a reference to an asset in the model hierarchy.
type AssetReference struct {
  // Fully qualified URL to the asset, in the form: file://domain.com/{container}/{application}#{url}
  ref string
  // Name of a container
  container string
  // Name of an application
  application string
  // URL for a file
  url string
}

// Returns an error if there is no container id.
func (asset *AssetReference) assertContainer() error {
  if asset.container == "" {
    return fmt.Errorf("Asset reference requires a container name")
  }
  return nil
}

// Returns an error if there is no container or application id.
func (asset *AssetReference) assertApplication() error {
  if asset.application == "" {
    return fmt.Errorf("Asset reference requires an application name")
  }
  return nil
}

// Returns an error if there is no container, application or file id.
func (asset *AssetReference) assertFile() error {
  if asset.url == "" {
    return fmt.Errorf("Asset reference requires a file url")
  }
  return nil
}

// Attempt to find a container matching this reference.
func (asset *AssetReference) FindContainer(host *Host) (*Container, *StatusError) {
  if err := asset.assertContainer(); err != nil {
    return nil, CommandError(http.StatusBadRequest, err.Error())
  }
  c := host.GetByName(asset.container)
  if c == nil {
    return nil, CommandError(http.StatusNotFound, "Container %s not found", asset.container)
  }
  return c, nil
}

// Attempt to find an application matching this reference.
func (asset *AssetReference) FindApplication(host *Host) (*Container, *Application, *StatusError) {
  if container, err := asset.FindContainer(host); err != nil {
    return nil, nil, err
  } else {
    if err := asset.assertApplication(); err != nil {
      return nil, nil, CommandError(http.StatusBadRequest, err.Error())
    }
    application := container.GetByName(asset.application)
    if application == nil {
      return nil, nil, CommandError(http.StatusNotFound, "Application %s not found", asset.application)
    }
    return container, application, nil
  }
}

// Attempt to find a file matching this reference.
func (asset *AssetReference) FindFile(host *Host) (*Container, *Application, *File, *StatusError) {
  if container, application, err := asset.FindApplication(host); err != nil {
    return nil, nil, nil, err
  } else {
    if err := asset.assertFile(); err != nil {
      return nil, nil, nil, CommandError(http.StatusBadRequest, err.Error())
    }
    file := application.Urls[asset.url]
    if file == nil {
      return nil, nil, nil, CommandError(http.StatusNotFound, "File %s not found", asset.url)
    }
    return container, application, file, nil
  }
}

// Parse a URL into this asset reference.
func (asset *AssetReference) ParseUrl(uri string) (ref Reference, err error) {
  var u *url.URL
  if u, err = url.Parse(uri); err != nil {
    return
  }
  path := strings.TrimPrefix(u.Path, "/")
  parts := strings.Split(path, "/")
  if len(parts) != 2 {
    err = fmt.Errorf("Invalid reference %s", uri)
    return
  }
  asset.container = parts[0]
  asset.application = parts[1]
  asset.url = u.Fragment
  asset.ref = uri
  ref = asset
  return
}
