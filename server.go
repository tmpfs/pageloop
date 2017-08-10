package pageloop

import (
  "log"
  "net/http"
  "path/filepath"
  "time"
  "github.com/tmpfs/pageloop/model"
)

var config ServerConfig
var mux *http.ServeMux

type PageLoop struct {}

type ServerConfig struct {
  AppPaths []string
}

type ServerHandler struct {}

func (h ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  handler, _ := mux.Handler(req)
  handler.ServeHTTP(res, req)
}

func (l *PageLoop) ServeHTTP(config ServerConfig) error {
  var err error

  // Initialize server mux
  mux = http.NewServeMux()

  //mux.Handle("/api/", apiHandler{})
  mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    // The "/" pattern matches everything, so we need to check
    // that we're at the root here.
    if req.URL.Path != "/" {
      http.NotFound(res, req)
      return
    }
    //fmt.Fprintf(res, "Welcome to the home page!")
  })

  if err = l.LoadApps(config); err != nil {
    return err
  }

  s := &http.Server{
    Addr:           ":3577",
    Handler:        ServerHandler{},
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  err = s.ListenAndServe()
  return err
}

func (l *PageLoop) LoadApps(config ServerConfig) error {
  var err error
  // iterate apps and configure paths
  for _, path := range config.AppPaths {

    var p string
    p, err = filepath.Abs(path)
    if err != nil {
      return err
    }
    name := filepath.Base(path)
    url := "/apps/" + name + "/"

    app := model.Application{}

    // Load the application files into memory
    if err = app.Load(p, nil); err != nil {
      return err
    }

    // Publish the application files to a build directory
    if err = app.Publish(nil); err != nil {
      return err
    }

    // Serve the static build files.
    log.Printf("Serving '%s' from %s", url, app.Public)
    mux.Handle(url, http.StripPrefix(url, http.FileServer(http.Dir(app.Public))))
  }

  return nil
}
