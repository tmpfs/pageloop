package main

import (
  "os"
  "fmt"
  "log"
  "flag"
  "github.com/tmpfs/pageloop"
)

func printHelp () {
  fmt.Println("pageloop")
  flag.PrintDefaults()
  os.Exit(0)
}

func main() {
  var help *bool
  var configPath *string

  configPath = flag.String("config", "", "path to a yaml configuration file")
  help = flag.Bool("help", false, "print help")

  flag.Parse()

  println("config path: " + *configPath)

  if *help {
    printHelp()
  }

  var apps []pageloop.Mountpoint
  apps = append(apps, pageloop.Mountpoint{Path: "test/fixtures/mock-app", Description: "Mock application."})
  loop := &pageloop.PageLoop{}
  conf := pageloop.ServerConfig{Mountpoints: apps, Addr: ":3577", Dev: true}
  server, err := loop.NewServer(conf)
  if err != nil {
    log.Fatal(err)
  }
  loop.Listen(server)
}
