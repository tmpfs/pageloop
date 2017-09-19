package handler

import(
  "log"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
)

// Mount an application such that it's published and source
// files are accessible over HTTP. This serves the published files
// as static files and serves two versions of the source file
// from in memory data. The src version is the file with any frontmatter
// stripped and the raw version includes frontmatter.
func MountApplication(mountpoints map[string]http.Handler, host *Host, app *Application) {
  listing := &DirList{Host: host}

	// Serve the static build files from the mountpoint path.
	url := app.PublishUrl()
	log.Printf("Serving app %s from %s", url, app.PublicDirectory())
  fileserver := http.FileServer(http.Dir(app.PublicDirectory()))
  mountpoints[url] = http.StripPrefix(url, PublicHandler{Listing: listing, App: app, FileServer: fileserver})
}
