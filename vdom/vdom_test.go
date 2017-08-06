package vdom

import (
  "os"
  "log"
  "testing"
  "io/ioutil"
  "golang.org/x/net/html"
)


func TestVdom(t *testing.T) {
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

  head := dom.Document.FirstChild.NextSibling.FirstChild
  body := head.NextSibling.NextSibling

  l1 := dom.CreateElement("link")
  err = dom.AppendChild(head, l1)
  if err != nil {
    t.Error(err)
  }

  foo := dom.CreateElement("foo")
  err = dom.AppendChild(l1, foo)
  if err != nil {
    t.Error(err)
  }

  l2 := dom.CreateElement("link")
  err = dom.AppendChild(head, l2)
  if err != nil {
    t.Error(err)
  }

  dom.RemoveChild(head, l1)

  meta := dom.CreateElement("meta")
  dom.InsertBefore(head, meta, l2)

  script := dom.CreateElement("script")
  err = dom.AppendChild(body, script)
  if err != nil {
    t.Error(err)
  }

  div := dom.CreateElement("div")
  err = dom.InsertBefore(body, div, script)
  if err != nil {
    t.Error(err)
  }

  err = html.Render(os.Stdout, dom.Document)
  if err != nil {
    t.Error(err)
  }
}

func TestDiff(t *testing.T) {
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

  head := dom.Document.FirstChild.NextSibling.FirstChild
  body := head.NextSibling.NextSibling

  div := dom.CreateElement("div")
  diff, err := dom.AppendDiff(body, div)
  if err != nil {
    t.Error(err)
  }

  log.Println(diff)

  /*
  d := Diff{
    Operation: APPEND_OP,
    Element: dom.GetId(head),
    Data: []byte(`<link />`)}
  log.Printf("%#v\n", d)
  */
}

