package main

import (
  "log"
	"net/http"
  "github.com/tmpfs/pageloop"
)

func main() {
  var err error
	var server *http.Server

  var apps []string
  apps = append(apps, "test/fixtures/mock-app")

  loop := &pageloop.PageLoop{}
	conf := pageloop.ServerConfig{AppPaths: apps, Addr: ":3577", Dev: true}
  server, err = loop.ServeHTTP(conf)
  if err != nil {
    log.Fatal(err)
  }
}
