// Data model for HTML pages written as template blocks.
package model

import(
	//"fmt"
  "bytes"
  "html/template"
  "encoding/json"
  "gopkg.in/yaml.v2"
  "golang.org/x/net/html"
  "github.com/tmpfs/pageloop/vdom"
	. "github.com/rhinoman/go-commonmark"
)

const(
	PageNone = iota
	PageHtml
	PageMarkdown
)

type Block struct {
  Title string  `json:"title"`
  Name string  `json:"name"`
  Content string  `json:"content"`
  Fragment bool
  Nodes []*html.Node
}

// Get the HTML document data.
//
// For HTML documents the underlying data is used, for markdown 
// documents they are converted to HTML first and the underlying 
// markdown data is left untouched.
func (p *Page) Data() []byte {
	if p.Type == PageMarkdown {
		return []byte(Md2Html(string(p.file.data), CMARK_OPT_DEFAULT))
	}
  return p.file.data
}

// Parse the file data into the virtual DOM of this page.
func (p *Page) Parse(data []byte) (*vdom.Vdom, error) {
  var err error
  var dom = vdom.Vdom{}
  err = dom.Parse(data)
  if err != nil {
    return nil, err
  }

  p.Dom = &dom

  return p.Dom, nil
}

// Render the current version of the virtual DOM to a byte 
// array. If page data is available parse the file as an 
// HTML template passing the page data to the template.
//
// Use this for the static rendered HTML for a page.
func (p *Page) Render(vdom *vdom.Vdom, node *html.Node) ([]byte, error) {
  var err error
  var data []byte
  if vdom == nil {
    vdom = p.Dom
  }

  if node == nil {
    node = vdom.Document
  }
  
  if data, err = vdom.RenderToBytes(node); err != nil {
    return nil, err
  }

  // Parse the file as a go HTML template if we have some 
  // page data.
  if p.PageDataType != DATA_NONE {
    tpl := template.New(p.file.Relative)
    tpl, err = tpl.Parse(string(data))
    if err != nil {
      return nil, err
    }

    w := new(bytes.Buffer)
    tpl.Execute(w, p.PageData)

    // Overwrite data with the parsed template.
    data = w.Bytes()

    // Prepend frontmatter to the output.
    /*
    if p.PageDataType == DATA_YAML {
      var fm []byte
      if fm, err = p.MarshalPageData(); err != nil {
        return nil, err
      }

      // Add document --- dashes.
      var delimiter []byte = []byte("---\n")
      fm = append(delimiter, fm...)
      fm = append(fm, delimiter...)

      // Prepend to the DOM bytes.
      data = append(fm, data...)
    }
    */
  } 
  return data, nil
}

func (p *Page) MarshalPageData() ([]byte, error) {
  if p.PageDataType == DATA_YAML || p.PageDataType == DATA_YAML_FILE {
    return yaml.Marshal(&p.PageData)
  } else if(p.PageDataType == DATA_JSON_FILE) {
    return json.Marshal(&p.PageData)
  }
  return nil, nil
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
