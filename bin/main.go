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

  var a *string
  var addr *string
  var c *string
  var config *string
  var p *string
  var publish *string

  a = flag.String("a", "", "")
  addr = flag.String("addr", "", "")

  c = flag.String("c", "", "")
  config = flag.String("config", "", "")

  p = flag.String("p", "", "")
  publish = flag.String("publish", "", "")

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

  // Merge user supplied config with the defaults
  if *c != "" {
    if err = conf.Merge(*c); err != nil {
      log.Fatal(err)
    }
  } else if *config != "" {
    if err = conf.Merge(*config); err != nil {
      log.Fatal(err)
    }
  }

  // Must be after the merge, overrides all config files
  if *a != "" {
    conf.Addr = *a
  } else if *addr != "" {
    conf.Addr = *addr
  }

  if *p != "" {
    conf.PublishDirectory = *p
  } else if *publish != "" {
    conf.PublishDirectory = *publish
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
