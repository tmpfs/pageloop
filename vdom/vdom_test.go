package vdom

import (
  "os"
  "log"
  "testing"
  "io/ioutil"
  "golang.org/x/net/html"
)

func TestVdom(t *testing.T) {
  log.Println("testing")
  file, err := ioutil.ReadFile("../test/fixtures/vdom.html")
  if err != nil {
    log.Fatal(err)
  }
  log.Println(string(file))

  dom, err := Parse(file, GetSettings())
  if err != nil {
    log.Fatal(err)
  }

  if dom == nil {
    t.Errorf("Expected vdom, got nil")
  }

  dom.AppendChild(dom.Document.FirstChild.FirstChild, dom.CreateElement("link"))

  //log.Printf("%v", dom.Map)

  html.Render(os.Stdout, dom.Document)
}
