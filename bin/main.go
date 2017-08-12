package main

import (
  "log"
  "github.com/tmpfs/pageloop"
)

func main() {
  var apps []pageloop.Mountpoint
	apps = append(apps, pageloop.Mountpoint{Path: "test/fixtures/mock-app"})
  loop := &pageloop.PageLoop{}
	conf := pageloop.ServerConfig{Mountpoints: apps, Addr: ":3577", Dev: true}
	server, err := loop.NewServer(conf)
  if err != nil {
    log.Fatal(err)
  }
	loop.Listen(server)
}
