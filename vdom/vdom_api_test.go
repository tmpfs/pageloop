package vdom

import (
  "log"
  "testing"
  "io/ioutil"
  //"golang.org/x/net/html"
)

// Tests the ability to use the the basic DOM API wrapper 
// functions, AppendChild, InsertBefore and RemoveChild.
func TestVdom(t *testing.T) {
  var err error
  var file []byte
  var dom *Vdom
  var expected string

  file, err = ioutil.ReadFile("../test/fixtures/vdom.html")
  if err != nil {
    log.Fatal(err)
  }

  dom = &Vdom{}
  err = dom.Parse(file)
  if err != nil {
    log.Fatal(err)
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

  data, err := dom.RenderToBytes(dom.Compact(dom.Document))
  //data, err := dom.RenderToBytes(dom.Document)
  if err != nil {
    t.Error(err)
  }

  expected = `<!DOCTYPE html><html data-id="0"><head data-id="0.0"><title data-id="0.0.0">Vdom Fixture</title><meta data-id="0.0.1"/><link data-id="0.0.2"/></head><body data-id="0.1"><div data-id="0.1.0"></div><script data-id="0.1.1"></script></body></html>`

  if string(data) != expected {
    t.Errorf("Unexpected markup, expected %s got %s", expected, string(data))
  }

  /*
  err = html.Render(os.Stdout, dom.Document)
  if err != nil {
    t.Error(err)
  }
  */
}
