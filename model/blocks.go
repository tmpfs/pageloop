// Data model for HTML pages written as template blocks.
package model

import(
  "os"
  //"io"
  //"fmt"
  //"log"
  "bytes"
  "golang.org/x/net/html"
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
  file File
}

type Block struct {
  Title string  `json:"title"`
  Name string  `json:"name"`
  Content string  `json:"content"`
  Fragment bool
  Nodes []*html.Node
}

// Parse blocks from the file data associated with a page.
func (p *Page) Parse() error {
  r := bytes.NewBuffer(p.file.data)

  _, err := html.Parse(r)
  if err != nil {
    return err
  }

  return err
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
