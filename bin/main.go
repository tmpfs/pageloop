package main

import (
  "log"
  "github.com/tmpfs/pageloop"
)

func main() {
  err := pageloop.ServeHTTP()
  if err != nil {
    log.Fatal(err)
  }
}
