package model

import(
  // "fmt"
  "strings"
  "net/url"
)

type Reference interface {
  ParseUrl(uri string) (ref Reference, err error)
}

// Represents a reference to an asset in the model hierarchy.
type AssetReference struct {
  // Name of a container
  Container string
  // Name of an application
  Application string
  // URL for a file
  Url string
  // Fully qualified URL to the asset, in the form: file://{container}/{application}#{url}
  Ref string
}

func (asset *AssetReference) ParseUrl(uri string) (ref Reference, err error) {
  var u *url.URL
  if u, err = url.Parse(uri); err != nil {
    return
  }
  parts := strings.Split(u.Path, "/")
  /*
  if len(parts) != 2 {
    err = fmt.Errorf("Invalid file reference %s", uri)
    return
  }
  */
  ref = &AssetReference{Container: u.Host, Application: parts[1], Url: u.Fragment, Ref: uri}
  return
}
