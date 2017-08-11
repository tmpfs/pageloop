package main

import (
  "log"
  "github.com/tmpfs/pageloop"
)

func main() {
  var err error

  var apps []pageloop.Mountpoint
	apps = append(apps, pageloop.Mountpoint{Path: "test/fixtures/mock-app"})

  loop := &pageloop.PageLoop{}
	conf := pageloop.ServerConfig{Mountpoints: apps, Addr: ":3577", Dev: true}
  _, err = loop.NewServer(conf)
  if err != nil {
    log.Fatal(err)
  }

	loop.Listen()
}
