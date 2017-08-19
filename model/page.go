// Data model for HTML pages written as template blocks.
package model

import(
	//"fmt"
  "bytes"
	"strings"
	"regexp"
	//"net/url"
  "html/template"
  "encoding/json"
	"path/filepath"
  "gopkg.in/yaml.v2"
  "golang.org/x/net/html"
  "golang.org/x/net/html/atom"
  "github.com/tmpfs/pageloop/vdom"
	. "github.com/rhinoman/go-commonmark"
)

const(
	PageNone = iota
	PageHtml
	PageMarkdown

	Layout = "layout.html"
	Content = "content"

)

var(
	FragmentTop = regexp.MustCompile(`^<html><head></head><body>`)
	FragmentTail = regexp.MustCompile(`</body></html>$`)
)

type Page struct {
  Path string `json:"-"`
  Name string `json:"name"`
  Url string `json:"url"`
  Uri string `json:"uri"`
  Mime string `json:"mime"`
  Size int64 `json:"size,omitempty"`
  PageData map[string] interface{} `json:"data"`
	PageDataType int `json:"-"`
  Blocks []Block  `json:"blocks"`
  Dom *vdom.Vdom `json:"-"`

	Type int `json:"-"`

	// Owner application
	owner *Application

	// The underlying file data.
  file *File
}

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
  return p.file.data
}

// Parse the file data into the virtual DOM of this page.
func (p *Page) Parse(data []byte) (*vdom.Vdom, error) {
  var err error

	// See if we can find a layout
	layout := p.FindLayout()

	if p.Type == PageMarkdown {
		data = []byte(Md2Html(string(data), CMARK_OPT_DEFAULT))
	}

  var dom = vdom.Vdom{}
	dom.Document = &html.Node{Type: html.DocumentNode}

	if layout != nil {
		var nodes []*html.Node
		var context *html.Node = &html.Node{Type: html.ElementNode, Data: "body", DataAtom:atom.Body}
		if nodes, err = dom.ParseFragment(data, context); err != nil {
			return nil, err
		}
		for _, n := range nodes {
			dom.Document.AppendChild(n)
		}
	} else {
		// Parse as full document.
		err = dom.Parse(data)
	}
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
	var path string = p.Path
	var dir string = filepath.Dir(path)
	var appRoot string = strings.TrimSuffix(p.owner.Root.Path, "/")

	// Do not process layout files
	if p.Name == Layout {
		return nil
	}

	var search func(dir string) *Page
	search = func(dir string) *Page {
		var target string = filepath.Join(dir, Layout)
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

// The default function map for templates.
func (p *Page) DefaultFuncMap() template.FuncMap {
	funcs := make(template.FuncMap)

	// Get a URL relative to the application root mountpoint.
	funcs["root"] = func(path string) string {
		path = strings.TrimPrefix(path, SLASH)
		return p.owner.Url + path
	}

	// Get a URL relative to the current page being parsed.

	/*
	funcs["relative"] = func(path string) (string, error) {
		println("relative func called with path: " + path)
		var err error
		domain := "http://localhost"
		page := p.owner.Url + strings.TrimPrefix(p.Url, SLASH)
		base := domain + page
		println("relative page: " + page)
		println("relative base: " + base)
		var bu *url.URL
		var tu *url.URL
		var ru *url.URL
		if bu, err = url.Parse(base); err != nil {
			return "", err
		}
		if tu, err = url.Parse(page); err != nil {
			return "", err
		}
		ru = bu.ResolveReference(tu)
		println(ru.Path)
		return path, nil
	}
	*/

  // Render markdown inline in an HTML template
  funcs["markdown"] = func(md string) template.HTML {
    var data string = Md2Html(md, CMARK_OPT_DEFAULT)
    return template.HTML(data)
  }

	return funcs
}

// Extract template delimiters from a map.
func (p *Page) GetDelims() (string, string) {
	var template map[string] interface{}
	var delims map[string]interface{}
	var ok bool
	var left string
	var right string
	if template, ok = p.PageData["template"].(map [string] interface{}); ok {
		if delims, ok = template["delims"].(map [string] interface{}); ok {
			if left, ok = delims["left"].(string); ok {
				if right, ok = delims["right"].(string); ok {
					return left, right
				}
			}
		}
	}
	return "{{", "}}"
}

// Creates a template and parses template source data.
func (p *Page) ParseTemplate(path string, source []byte, funcs template.FuncMap, layout bool) (*template.Template, error) {
	tpl := template.New(path)
	if funcs != nil {
		tpl.Funcs(funcs)
	}

	// set delims if necessary
	if !layout && p.PageData != nil {
		left, right := p.GetDelims()
		tpl.Delims(left, right)
	}

	return tpl.Parse(string(source))
}

// Execute a template with the given data.
func (p *Page) ExecuteTemplate(tpl *template.Template, data map[string] interface{}) ([]byte, error) {
	var err error
	w := new(bytes.Buffer)
	if err = tpl.Execute(w, data); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// Render with a clean DOM and update the page's file reference
// with the new rendered data.
func (p *Page) Update() error {
	var err error
	var data []byte
	node := p.Dom.Clean(nil)
	if data, err = p.Render(p.Dom, node); err != nil {
		return err
	}
	p.file.data = data
	return nil
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

  if data, err = vdom.RenderRaw(node); err != nil {
    return nil, err
  }

	// Do not handle layout files
	if p.Name == Layout {
		return nil, nil
	}

	// Create and parse the primary template
	var primary func(layout bool) (*template.Template, error)
	primary = func(layout bool) (*template.Template, error) {
		var src []byte = data
		left, right := p.GetDelims()
		var define []byte = []byte(left + `define "content"` + right)
		var end []byte = []byte(left + `end` + right)
		if layout {
			src = append(define, data...)
			src = append(src, end...)
			//println(string(src))

			// FIXME: the output of calling golang.org/x/net/html#Render
			// FIXME: escapes the quotes in template nodes so we need
			// FIXME: to unescape them - there must be a better way ;)

			// FIXME: and if we actually use this we break markdown code blocks including markup :(
			//src = []byte(html.UnescapeString(string(src)))
		}
		return p.ParseTemplate(p.file.Path, src, p.DefaultFuncMap(), false)
	}

	// See if we can find a layout
	layout := p.FindLayout()
	if layout != nil {
		// Markdown documents when going via the vdom are rendered
		// with outer html, head and body elements. Markdown documents
		// that are part of a layout need these removed.
		if p.Type == PageMarkdown {
			data = FragmentTop.ReplaceAll(data, []byte{})
			data = FragmentTail.ReplaceAll(data, []byte{})
		}


		if tpl, err = primary(true); err != nil {
			return nil, err
		}
		// Configure the layout template
		file := layout.file
		var lyt *template.Template
		if lyt, err = p.ParseTemplate(file.Path, file.Source(false), p.DefaultFuncMap(), true); err != nil {
			return nil, err
		}

		// Add the content template
		for _, t := range tpl.Templates() {
			if t.Name() == Content {
				if _, err = lyt.AddParseTree(Content, t.Tree); err != nil {
					return nil, err
				}
				break
			}
		}

		// Execute the outer layout template
		if result, err = p.ExecuteTemplate(lyt, p.PageData); err != nil {
			return nil, err
		}
		data = result
	// Template without layout, execute directly
	} else {
		if tpl, err = primary(false); err != nil {
			return nil, err
		}
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
