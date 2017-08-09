package pageloop

import (
  "fmt"
  "net/http"
  "time"
)

type ServerHandler struct {
  
}

func (h ServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  mux := http.NewServeMux()
  //mux.Handle("/api/", apiHandler{})
  mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    println("running handler")
      // The "/" pattern matches everything, so we need to check
      // that we're at the root here.
      if req.URL.Path != "/" {
          http.NotFound(res, req)
          return
      }
      fmt.Fprintf(res, "Welcome to the home page!")
  })

  handler, _ := mux.Handler(req)
  handler.ServeHTTP(res, req)
}

func ServeHTTP() error {
  var err error
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
