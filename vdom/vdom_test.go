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

  var diff *Diff
  var err error
  var data []byte
  var expected string
  var dom *Vdom
  var p Patch = Patch{}


  file, err := ioutil.ReadFile("../test/fixtures/vdom.html")
  if err != nil {
    log.Fatal(err)
  }

  //log.Println(string(file))

  dom, err = Parse(file)
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
  diff, err = dom.AppendDiff(body, div)
  if err != nil {
    t.Error(err)
  }

  if diff.Operation != APPEND_OP {
    t.Errorf("Unexpected operation, expected %d got %d", APPEND_OP, diff.Operation)
  }

  expected = "<div></div>"
  if expected != string(diff.Data) {
    t.Errorf("Unexpected diff data, expected %s got %s", expected, string(diff.Data))
  }

  // div must be on the DOM for the insert before next
  // so we create it and omit it from the patch
  dom.AppendChild(body, div)

  log.Printf("%s\n", string(diff.Data))

  // insert paragraph before the div
  para := dom.CreateElement("p")
  diff, err = dom.InsertDiff(body, para, div)
  if err != nil {
    t.Error(err)
  }

  log.Printf("%s\n", string(diff.Data))

  if diff.Operation != INSERT_OP {
    t.Errorf("Unexpected operation, expected %d got %d", INSERT_OP, diff.Operation)
  }

  expected = "<p></p>"
  if expected != string(diff.Data) {
    t.Errorf("Unexpected diff data, expected %s got %s", expected, string(diff.Data))
  }

  p.Add(diff)

  // remove paragraph before the div
  diff, err = dom.RemoveDiff(body, para)
  if err != nil {
    t.Error(err)
  }

  log.Printf("%s\n", string(diff.Data))

  if diff.Operation != REMOVE_OP {
    t.Errorf("Unexpected operation, expected %d got %d", REMOVE_OP, diff.Operation)
  }

  expected = "<p></p>"
  if expected != string(diff.Data) {
    t.Errorf("Unexpected diff data, expected %s got %s", expected, string(diff.Data))
  }

  p.Add(diff)

  diff, err = dom.SetAttrDiff(para, html.Attribute{Key: "data-foo", Val: "bar"})
  if err != nil {
    t.Error(err)
  }

  if diff.Operation != ATTR_SET_OP {
    t.Errorf("Unexpected operation, expected %d got %d", ATTR_SET_OP, diff.Operation)
  }

  data, err = dom.RenderToBytes(para)
  if err != nil {
    t.Error(err)
  }

  expected = "<p data-foo=\"bar\"></p>"
  if expected != string(data) {
    t.Errorf("Unexpected diff data, expected %s got %s", expected, string(data))
  }

  p.Add(diff)

  log.Println(string(data))
  //log.Printf("%#v\n", diff)

  diff, err = dom.DelAttrDiff(para, html.Attribute{Key: "data-foo"})
  if err != nil {
    t.Error(err)
  }

  if diff.Operation != ATTR_DEL_OP {
    t.Errorf("Unexpected operation, expected %d got %d", ATTR_DEL_OP, diff.Operation)
  }

  data, err = dom.RenderToBytes(para)
  if err != nil {
    t.Error(err)
  }

  expected = "<p></p>"
  if expected != string(data) {
    t.Errorf("Unexpected diff data, expected %s got %s", expected, string(data))
  }

  p.Add(diff)

  log.Println(string(data))

  // apply the patch to perform the operations
  err = dom.Apply(&p)
  if err != nil {
    t.Error(err)
  }

  // debug
  err = html.Render(os.Stdout, dom.Document)
  if err != nil {
    t.Error(err)
  }
}

