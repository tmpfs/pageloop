package model

import(
  "fmt"
  "strings"
  "net/url"
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

func (asset *AssetReference) AssertContainer() error {
  if asset.container == "" {
    return fmt.Errorf("Asset reference requires a container name")
  }
  return nil
}

func (asset *AssetReference) AssertApplication() error {
  if err := asset.AssertContainer(); err != nil {
    return err
  }
  if asset.application == "" {
    return fmt.Errorf("Asset reference requires an application name")
  }
  return nil
}

func (asset *AssetReference) AssertFile() error {
  if err := asset.AssertApplication(); err != nil {
    return err
  }
  if asset.url == "" {
    return fmt.Errorf("Asset reference requires a file url")
  }
  return nil
}

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
  ref = &AssetReference{
    container: parts[0],
    application: parts[1],
    url: u.Fragment,
    ref: uri}
  return
}
