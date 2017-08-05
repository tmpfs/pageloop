package blocks

import(
  // . "fmt"
)

const (
  HTML5 = "<!doctype html>"
)

type Page struct {
  DocType string `json:"doctype"`
  UserData map[string] interface{} `json:"data"`
  Blocks []Block  `json:"blocks"`
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
