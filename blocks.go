/*
  Data model for HTML pages written as template blocks.
*/
package blocks

import(
  "os"
  // . "fmt"
)

const (
  HTML5 = "<!doctype html>"
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
  Content string `json:"content"`
}

type TextBlock struct {
  Block
}

type MarkdownBlock struct {
  Block
}

type HtmlBlock struct {
  Block
}

func (p *Page) RemoveBlock(b Block) Page {
  for i, child := range p.Blocks {
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

func (p *Page) AddBlock(b Block) Page {
  p.Blocks = append(p.Blocks, b)
  return *p
}

func (p *Page) Length() int {
  return len(p.Blocks)
}
