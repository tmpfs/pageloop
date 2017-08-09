// Data model for HTML pages written as template blocks.
package model

import(
  "os"
  //"io"
  //"fmt"
  //"log"
  //"bytes"
  "html/template"
  "golang.org/x/net/html"
  "github.com/tmpfs/pageloop/vdom"
)

const (
  HTML5 = "html"
)

type Application struct {
  Name string `json: "name"`
  Title string `json:"title"`
  Pages []Page `json:"pages"`
  Files []File `json:"files"`
  Base string `json:"base"`
  Urls map [string] File
}

type File struct {
  Path string `json:"path"`
  Directory bool `json:"directory"`
  Relative string `json:"relative"`
  Url string `json:"url"` 
  Index bool `json:"index"`
  info os.FileInfo
  data []byte
}

type Page struct {
  File
  DocType string `json:"doctype"`
  UserData map[string] interface{} `json:"data"`
  Blocks []Block  `json:"blocks"`
  Dom *vdom.Vdom
  file File
}

type Block struct {
  Title string  `json:"title"`
  Name string  `json:"name"`
  Content string  `json:"content"`
  Fragment bool
  Nodes []*html.Node
}

// Parse the file data into the virtual DOM of this page.
func (p *Page) Parse() (*vdom.Vdom, error) {
  var err error
  var dom = vdom.Vdom{}
  err = dom.Parse(p.file.data)
  if err != nil {
    return nil, err
  }

  p.Dom = &dom

  return p.Dom, nil
}

func (p *Page) Render() error {
  //fmt.Println("--- render function ---")
  //fmt.Println()

  data := p.file.data
  tpl := template.New(p.file.Relative)
  tpl, err := tpl.Parse(string(data))
  if err != nil {
    return err
  }
  //log.Println(app.Pages[0].UserData)
  tpl.Execute(os.Stdout, p.UserData)

  return nil
}

/*
func (p *Page) RemoveBlock(b Block) Page {
  for i, child := range p.Blocks {
    log.Printf("%#v\n", child)
    log.Printf("%#v\n", b)
    if child == b {
      before := p.Blocks[0:i]
      if p.Length() > i {
        after := p.Blocks[i+1:]
        p.Blocks = append(before, after...)
      }
      break
    }
  } 
  return *p
}
*/

func (p *Page) AddBlock(b Block) Page {
  p.Blocks = append(p.Blocks, b)
  return *p
}

func (p *Page) Length() int {
  return len(p.Blocks)
}
