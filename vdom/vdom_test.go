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

  /*
  err = html.Render(os.Stdout, dom.Document)
  if err != nil {
    t.Error(err)
  }
  */
}

func TestDiff(t *testing.T) {
  file, err := ioutil.ReadFile("../test/fixtures/vdom.html")
  if err != nil {
    log.Fatal(err)
  }

  //log.Println(string(file))

  dom, err := Parse(file)
  if err != nil {
    log.Fatal(err)
  }

  if dom == nil {
    t.Errorf("Expected vdom, got nil")
  }

  head := dom.Document.FirstChild.NextSibling.FirstChild
  body := head.NextSibling.NextSibling

  // append a div
  div := dom.CreateElement("div")
  diffa, err := dom.AppendDiff(body, div)
  if err != nil {
    t.Error(err)
  }

  if diffa.Operation != APPEND_OP {
    t.Errorf("Unexpected operation, expected %d got %d", APPEND_OP, diffa.Operation)
  }

  log.Printf("%s\n", string(diffa.Data))
  log.Printf("%#v\n", diffa)

  // insert paragraph before the div
  para := dom.CreateElement("p")
  diffi, err := dom.InsertDiff(body, para, div)
  if err != nil {
    t.Error(err)
  }

  log.Printf("%s\n", string(diffi.Data))
  log.Printf("%#v\n", diffi)

  if diffi.Operation != INSERT_OP {
    t.Errorf("Unexpected operation, expected %d got %d", INSERT_OP, diffi.Operation)
  }

  // remove paragraph before the div
  diffr, err := dom.RemoveDiff(body, para)
  if err != nil {
    t.Error(err)
  }

  log.Printf("%s\n", string(diffr.Data))
  log.Printf("%#v\n", diffr)

  if diffr.Operation != REMOVE_OP {
    t.Errorf("Unexpected operation, expected %d got %d", REMOVE_OP, diffr.Operation)
  }

  // debug
  err = html.Render(os.Stdout, dom.Document)
  if err != nil {
    t.Error(err)
  }
}

