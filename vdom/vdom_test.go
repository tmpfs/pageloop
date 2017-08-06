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

  l1 := dom.CreateElement("link")
  err = dom.AppendChild(head, l1)
  if err != nil {
    t.Error(err)
  }

  l2 := dom.CreateElement("link")
  err = dom.AppendChild(head, l2)
  if err != nil {
    t.Error(err)
  }

  dom.RemoveChild(head, l1)

  //log.Printf("%v", dom.Map)

  html.Render(os.Stdout, dom.Document)
}
