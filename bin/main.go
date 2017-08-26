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
  var err error
  var help *bool
  var configPath *string

  configPath = flag.String("config", "", "path to a yaml configuration file")
  help = flag.Bool("help", false, "print help")

  flag.Parse()

  if *help {
    printHelp()
  }

  loop := &pageloop.PageLoop{}
  conf := pageloop.DefaultServerConfig()

  if *configPath != "" {
    // Merge user supplied config with the defaults
    if err = conf.Merge(*configPath); err != nil {
      log.Fatal(err)
    }
  }

  server, err := loop.NewServer(conf)
  if err != nil {
    //log.Fatal(err)
    //fmt.Errorf(err)
    panic(err)
  }
  loop.Listen(server)
}
