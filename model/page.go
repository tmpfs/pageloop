// Data model for HTML pages written as template blocks.
package model

import(
	//"fmt"
  "bytes"
	"strings"
  "html/template"
  "encoding/json"
	"path/filepath"
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

// Attempts to find a layout.html file for a page.
//
// Starts by looking in the directory containing the file 
// and walks the parent paths until it hits the application 
// root searching for a layout file.
func (p *Page) FindLayout() *Page {
	var layout string = "layout.html"
	var path string = p.Path
	var dir string = filepath.Dir(path)
	var appRoot string = strings.TrimSuffix(p.owner.Root.Path, "/")

	// Do not process layout files
	if p.Name == layout {
		return nil
	}

	var search func(dir string) *Page
	search = func(dir string) *Page {
		var target string = filepath.Join(dir, layout)
		for _, p := range p.owner.Pages {
			if p.file.Path == target {
				return p
			}
		}
		if dir == appRoot {
			return nil
		}
		return search(filepath.Dir(dir))
	}

	return search(dir)
}

func (p *Page) ParseTemplate(path string, source []byte, funcs template.FuncMap) (*template.Template, error) {
	tpl := template.New(path)
	if funcs != nil {
		tpl.Funcs(funcs)
	}
	return tpl.Parse(string(source))
}

// Parse a template file and execute it with the given data.
func (p *Page) ExecuteTemplate(tpl *template.Template, data map[string] interface{}) ([]byte, error) {
	var err error
	w := new(bytes.Buffer)
	if err = tpl.Execute(w, data); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// Render the current version of the virtual DOM to a byte 
// array. If page data is available parse the file as an 
// HTML template passing the page data to the template.
//
// Use this for the static rendered HTML for a page.
func (p *Page) Render(vdom *vdom.Vdom, node *html.Node) ([]byte, error) {
  var err error
  var data []byte
  var result []byte
	var tpl *template.Template
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
  //if p.PageDataType != DATA_NONE {

	if p.Name == "layout.html" {
		return nil, nil	
	}

	tpl, err = p.ParseTemplate(p.file.Path, p.file.Source(), nil)
	if err != nil {
		return nil, err
	}

	// see if we need to render as part of a layout
	layout := p.FindLayout()
	if layout != nil {
		file := layout.file
		var lyt *template.Template
		if lyt, err = p.ParseTemplate(file.Path, file.Source(), nil); err != nil {
			return nil, err
		}

		//println(tpl.DefinedTemplates())

		for _, t := range tpl.Templates() {
			if t.Name() == "content" {
				if _, err = lyt.AddParseTree("content", t.Tree); err != nil {
					return nil, err
				}
				break
			}
		}

		if result, err = p.ExecuteTemplate(lyt, p.PageData); err != nil {
			return nil, err
		}
		data = result
	// Template without layout
	} else {
		if result, err = p.ExecuteTemplate(tpl, p.PageData); err != nil {
			return nil, err
		}
		data = result
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
