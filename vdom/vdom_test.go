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

  dom, err := Parse(file)
  if err != nil {
    log.Fatal(err)
  }

  if dom == nil {
    t.Errorf("Expected vdom, got nil")
  }

  //log.Println(dom.Document.FirstChild)
  //log.Println(dom.Document.FirstChild.NextSibling.FirstChild)

  head := dom.Document.FirstChild.NextSibling.FirstChild

  err = dom.AppendChild(head, dom.CreateElement("link"))
  if err != nil {
    t.Error(err)
  }

  //log.Printf("%v", dom.Map)

  html.Render(os.Stdout, dom.Document)
}
