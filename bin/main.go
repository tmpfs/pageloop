package main

import (
  "log"
  "github.com/tmpfs/pageloop"
)

func main() {
  var err error

  var apps []string
  apps = append(apps, "test/fixtures/mock-app")

  loop := &pageloop.PageLoop{}
	conf := pageloop.ServerConfig{AppPaths: apps, Addr: ":3577", Dev: true}
  _, err = loop.NewServer(conf)
  if err != nil {
    log.Fatal(err)
  }

	loop.Listen()
}
