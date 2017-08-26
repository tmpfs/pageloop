package main

import (
  "os"
  "fmt"
  "log"
  "flag"
  "github.com/tmpfs/pageloop"
)

var helpText []byte

func printHelp () {
  os.Stdout.Write(helpText)
  os.Exit(0)
}

func printVersion () {
  fmt.Printf("%s %s\n", pageloop.Name, pageloop.Version)
  os.Exit(0)
}

func main() {
  var err error
  var h *bool
  var help *bool
  var version *bool

  var config *string
  var addr *string

  addr = flag.String("addr", "", "")
  config = flag.String("config", "", "")

  h = flag.Bool("h", false, "")
  help = flag.Bool("help", false, "")
  version = flag.Bool("version", false, "")

  flag.Parse()

  if *h || *help {
    printHelp()
  }

  if *version {
    printVersion()
  }

  loop := &pageloop.PageLoop{}
  conf := pageloop.DefaultServerConfig()

  if *addr != "" {
    println("setting addr: " + *addr)
    conf.Addr = *addr
  }

  if *config != "" {
    // Merge user supplied config with the defaults
    if err = conf.Merge(*config); err != nil {
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

func init() {
  helpText = pageloop.MustAsset("pageloop.txt")
}
