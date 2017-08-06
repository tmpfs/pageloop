/*
  Data model for HTML pages written as template blocks.
*/
package blocks

import(
  "os"
  "log"
  "bytes"
  "golang.org/x/net/html"
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

/*
  Parse blocks from the file data associated with this page.
*/
func (p *Page) Parse() error {
  r := bytes.NewBuffer(p.file.data)

  /*
  doc, err := html.Parse(r)
  if err != nil {
    log.Fatal(err)
  }
  */

  /*
  var f func(*html.Node)
  f = func(n *html.Node) {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      log.Printf("type: %d\n", c.Type)
      if c.Type == html.ElementNode {
        log.Println(c.Data)
        log.Println(c.Attr)
        for _, attrs := range c.Attr {
          if attrs.Key == "id" {
            log.Println("got id attribute")
          }
        }
      }

      f(c)
    }
  }

  */

  z := html.NewTokenizer(r)
  depth := 0
  doc := html.Node{Type: html.DocumentNode}

  var parent html.Node  = doc
  var current []html.Node
  var node html.Node

  // this way we get the doctype if we call Parse()
  // the doctype is dropped during parsing and we only
  // have the node parse tree
  for {
    tt := z.Next()
    token := z.Token()
    log.Printf("%#v", token)
    log.Printf("depth %d", depth)
    switch tt {
      case html.ErrorToken:
        return z.Err()
      case html.DoctypeToken:
        p.DocType = token.Data
        log.Printf("doctype %s\n", p.DocType)
      case html.TextToken:
        if depth > 0 {
          node = html.Node{Type: html.TextNode, Parent: &parent}
        }
      case html.StartTagToken, html.EndTagToken:
        tn, _ := z.TagName()
        if tt == html.StartTagToken {
          node = html.Node{Type: html.ElementNode, Parent: &parent, Data: name}
          current = append(current, node)
          nextParent = node
          depth++
        } else {
          depth--
        }
    }

    parent = nextParent

    log.Printf("%#v", node)
  }

  //f(doc)

  return nil
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
