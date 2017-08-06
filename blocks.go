/*
  Data model for HTML pages written as template blocks.
*/
package blocks

import(
  "os"
  //"io"
  "fmt"
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
  Name string  `json:"name"`
  Content string  `json:"content"`
  Fragment bool
  Nodes []*html.Node
}

type BlockNode interface {
  SetName(name string)
}

func (b *Block) SetName(name string) {
  b.Name = name
}

/*
type TextBlock struct {
  Block
}

type MarkdownBlock struct {
  Block
}

type HtmlBlock struct {
  Block
}
*/

/*
  Parse blocks from the file data associated with this page.
*/
func (p *Page) Parse() error {
  r := bytes.NewBuffer(p.file.data)

  doc, err := html.Parse(r)
  if err != nil {
    return err
  }

  var current Block = Block{Fragment: true}

  var addNode func (*html.Node)
  addNode = func(c *html.Node) {
    current.Nodes = append(current.Nodes, c)
  }

  var addBlock func(b Block)
  addBlock = func(b Block) {
    if (len(b.Nodes) > 0) {
      p.Blocks = append(p.Blocks, b)
    }
  }

  var f func(*html.Node)
  depth := 0
  f = func(n *html.Node) {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      if c.Type == html.DoctypeNode {
        p.DocType = c.Data
        addNode(c)
      } else if c.Type == html.ElementNode {
        // what is a block when parsing by default?
        // - header elements
        // - footer elements
        // - section elements
        // - article elements

        switch c.Data {
          case "header":
            fallthrough
          case "footer":
            fallthrough
          case "article":
            fallthrough
          case "section":
            // close current block
            addBlock(current)
            // add the node as the only child
            current = Block{}
            addNode(c)
            addBlock(current)
            // create a new fragment block
            current = Block{Fragment: true}
          default: 
            depth++
            f(c)
        }
      } else {
      }
    }
    depth--
  }

  f(doc)

  log.Println("block length", len(p.Blocks))

  for _, block := range p.Blocks {
    d := html.Node{Type: html.DocumentNode}
    fmt.Println("render block", block.Fragment)
    for _, c := range block.Nodes {

      c.Parent = nil
      c.NextSibling = nil
      c.PrevSibling = nil
      d.AppendChild(c)
    }
    //d.AppendChild()
    html.Render(os.Stdout, &d) 
    fmt.Println()
    fmt.Println()
  }

  /*
  z := html.NewTokenizer(r)
  depth := 0
  doc := html.Node{Type: html.DocumentNode}

  var err error

  var parent html.Node  = doc
  //var nextParent html.Node  = doc
  var parents []html.Node
  var node html.Node

  // this way we get the doctype if we call Parse()
  // the doctype is dropped during parsing and we only
  // have the node parse tree
  iterate:
    for {
      tt := z.Next()
      token := z.Token()
      //log.Printf("%#v", token)
      //log.Printf("depth %d", depth)
      //log.Printf("parent data %s\n", parent.Data)
      switch tt {
        case html.ErrorToken:
          err = z.Err()
          break iterate
        case html.DoctypeToken:
          p.DocType = token.Data
          log.Printf("doctype %s\n", p.DocType)
        case html.TextToken:
          if depth > 0 {
            node = html.Node{Type: html.TextNode, Data: token.Data, DataAtom: token.DataAtom}
            parent.AppendChild(&node)
            log.Printf("text node %s\n", token.Data)
          }
        case html.CommentToken:
            node = html.Node{Type: html.CommentNode, Data: token.Data, DataAtom: token.DataAtom}
            parent.AppendChild(&node)
        case html.SelfClosingTagToken:
            node = html.Node{Type: html.ElementNode, Data: token.Data, DataAtom: token.DataAtom, Attr: token.Attr}
            parent.AppendChild(&node)
        case html.StartTagToken, html.EndTagToken:
          if tt == html.StartTagToken {
            log.Printf("open element node %s\n", token.Data)
            node = html.Node{Type: html.ElementNode, Data: token.Data, DataAtom: token.DataAtom, Attr: token.Attr}
            parents = append(parents, node)
            parent.AppendChild(&node)
            //log.Printf("%#v", parent.FirstChild)
            parent = node
            depth++
          } else {
            log.Printf("close element node %s\n", token.Data)
            depth--
            parent = parents[len(parents) - 1]
            parents = parents[0:len(parents) - 1]
          }
      }

      //log.Printf("%#v", node)
      //log.Printf("%#v", node.Attr)
    }

  if err == io.EOF {
    log.Printf("RENDERING %#v", doc.FirstChild.FirstChild)
    log.Printf("RENDERING %#v", &doc)
    err = html.Render(os.Stdout, &doc)
    log.Printf("RENDER COMPLETE\n")
  }
  */


  //err = html.Render(os.Stdout, doc)

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
