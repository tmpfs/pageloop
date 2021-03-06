// Data model for HTML pages written as template blocks.
package model

import(
  //"fmt"
  "os"
  "bytes"
	"strings"
	"regexp"
	//"net/url"
  "io/ioutil"
  "html/template"
  "encoding/json"
	"path/filepath"
  "gopkg.in/yaml.v2"
  "golang.org/x/net/html"
  "golang.org/x/net/html/atom"
  "github.com/tmpfs/pageloop/vdom"
	. "github.com/rhinoman/go-commonmark"
  . "github.com/tmpfs/pageloop/util"
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
  Binary bool `json:"binary"`
  Size int64 `json:"size,omitempty"`
  PrettySize string `json:"filesize,omitempty"`
  PageData map[string] interface{} `json:"data"`
	PageDataType int `json:"-"`
  Blocks []Block  `json:"blocks,omitempty"`
  Dom *vdom.Vdom `json:"-"`

	Type int `json:"-"`

	// Owner application
  Owner *Application `json:"-"`

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
  var name string = Layout
  if p.PageData != nil {
    // Layout disabled in page data
    if layout, ok := p.PageData["layout"].(bool); ok {
      if !layout {
        return nil
      }
    }

    // String value specifies name of the layout file
    if layoutName, ok := p.PageData["layout"].(string); ok {
      if layoutName != "" {
        name = layoutName
      }
    }
  }
	var path string = p.Path
	var dir string = filepath.Dir(path)
	var appRoot string = strings.TrimSuffix(p.Owner.SourceDirectory(), "/")

	// Do not process layout files
	//if p.Name == name {
		//return nil
	//}

	var search func(dir string) *Page
	search = func(dir string) *Page {
		var target string = filepath.Join(dir, name)
		for _, p := range p.Owner.Pages {
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
		return p.Owner.Url + path
	}

  // Render markdown inline in an HTML template
  funcs["markdown"] = func(md string) template.HTML {
    var data string = Md2Html(md, CMARK_OPT_DEFAULT)
    return template.HTML(data)
  }

  // Pretty print byte sizes
  funcs["prettybytes"] = func(size int64) string {
    return PrettyBytes(size)
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
func (p *Page) ExecuteTemplate(tpl *template.Template, data interface{}) ([]byte, error) {
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

  // User data has template flag set, can be used
  // to disable template parsing
  if userflag, ok := p.PageData["template"].(bool); ok {
    // Template parsing disabled!
    if !userflag {
      return data, nil
    }
  }

	// Do not handle layout files
  /*
	if p.Name == Layout {
		return nil, nil
	}
  */

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

			// FIXME: the output of calling golang.org/x/net/html#Render
			// FIXME: escapes the quotes in template nodes so we need
			// FIXME: to unescape them - there must be a better way ;)

			// FIXME: and if we actually use this we break markdown code blocks including markup :(
			//src = []byte(html.UnescapeString(string(src)))
		}
		return p.ParseTemplate(p.file.Path, src, p.DefaultFuncMap(), false)
	}

  doInclude := func(tpl *template.Template) error {
    var err error
    var content []byte
    var partial *template.Template
    if includes, ok := p.PageData["includes"].([]interface{}); ok {
      for _, inc := range includes {
        if includePath, ok := inc.(string); ok {
          includePath = filepath.Clean(includePath)
          fullPath := filepath.Join(p.Owner.SourceDirectory(), strings.TrimPrefix(includePath, "/"))
          if content, err = ioutil.ReadFile(fullPath); err != nil {
            return err
          }
          if partial, err = p.ParseTemplate(fullPath, content, p.DefaultFuncMap(), false); err != nil {
            return err
          }
          for _, t := range partial.Templates() {
            tpl.AddParseTree(t.Name(), t.Tree)
          }
        }
      }
    }
    return nil
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

		// Add the content template to the parse tree if we are not the layout itself
    if layout != p {
      for _, t := range tpl.Templates() {
        if t.Name() == Content {
          if _, err = lyt.AddParseTree(Content, t.Tree); err != nil {
            return nil, err
          }
          break
        }
      }
    }

    doInclude(lyt)

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
    doInclude(tpl)
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

func (page *Page) ParsePageData () error {
  var err error
  if _, err = page.parsePageData(); err != nil {
    return err
  }

  if page.file.data != nil {
    if _, err = page.Parse(page.file.data); err != nil {
      return err
    }
  }

  //app.Pages[index] = page

  // Convert map[interface{}] interface{} to map[string] interface{}
  // recursively in parsed page data.
  //
  // This is necessary as the YAML unmarshaller converts nested maps
  // using interface{} keys and interface{} keys are not recognised when
  // marshalling to JSON. We don't want to define a type for page data as
  // it is intended to be arbitrary.
  var coerce func(m map[string] interface{})
  coerce = func(m map[string] interface{}) {
    for key, value := range m {
      if val, ok := value.(map[interface{}] interface{}); ok {
        var r map[string]interface{}
        r = make(map[string]interface{})
        for k, v := range val {
          r[k.(string)] = v
        }
        coerce(r)
        m[key] = r
      }
    }
  }
  coerce(page.PageData)
  return nil
}

// Attempt to find user page data by first attempting to
// parse embedded frontmatter YAML.
//
// If there is no frontmatter data it attempts to
// load data from a corresponding file with a .yml extension.
//
// Finally if a .json file exists it is parsed.
func (page *Page) parsePageData() (map[string] interface{}, error) {
  page.PageData = make(map[string] interface{})
  page.PageDataType = DATA_NONE

  // frontmatter
  if FRONTMATTER.Match(page.file.source) {
    var read int = 4
    var lines [][]byte = bytes.Split(page.file.source, []byte("\n"))
    var frontmatter [][]byte
    // strip leading ---
    lines = lines[1:]
    for _, line := range lines {
      read += len(line) + 1
      if FRONTMATTER_END.Match(line) {
        break
      }
      frontmatter = append(frontmatter, line)
    }

    if len(frontmatter) > 0 {
      fm := bytes.Join(frontmatter, []byte("\n"))
      err := yaml.Unmarshal(fm, &page.PageData)
      if err != nil {
        return nil, err
      }

      //println(string(page.file.data))

      // strip frontmatter content from file data after parsing
      //page.file.data = page.file.data[read:]
      page.file.source = page.file.source[read:]
      fm = append([]byte("---\n"), fm...)
      fm = append(fm, []byte("\n---\n")...)
      page.file.frontmatter = fm

      //println(string(fm))

      page.PageDataType = DATA_YAML
    }
    return page.PageData, nil
  }

  // external files
  dir, name := filepath.Split(page.file.Path)
  for _, dataType := range types {
    dataPath := name
    if TEMPLATE_FILE.MatchString(name) {
      dataPath = TEMPLATE_FILE.ReplaceAllString(name, dataType)
    } else if MARKDOWN_FILE.MatchString(name) {
      dataPath = MARKDOWN_FILE.ReplaceAllString(name, dataType)
    }
    dataPath = filepath.Join(dir, dataPath)
    // Failed to change file extension
    if dataPath == page.file.Path {
      return nil, nil
    }

    fh, err := os.Open(dataPath)
    if err != nil {
      if !os.IsNotExist(err) {
        return nil, err
      }
    }
    if fh != nil {
      defer fh.Close()
      contents, err := ioutil.ReadFile(dataPath)
      if err != nil {
        return nil, err
      }

      if dataType == JSON {
        err = json.Unmarshal(contents, &page.PageData)
        if err != nil {
          return nil, err
        }
        page.PageDataType = DATA_JSON_FILE
      } else if dataType == YAML {
        err = yaml.Unmarshal(contents, &page.PageData)
        if err != nil {
          return nil, err
        }
        page.PageDataType = DATA_YAML_FILE
      }
      break
    }
  }

  return page.PageData, nil
}
