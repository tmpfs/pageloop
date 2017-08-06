package vdom
  
import(
  "log"
  "bytes"
  "strconv"
  "strings"
  "golang.org/x/net/html"
)

const (
  ADD_OP = iota
  REMOVE_OP
  ATTR_OP
)

type Vdom struct {
  Document *html.Node
  Map map[string] *html.Node
}

func (vdom *Vdom) AppendChild(parent *html.Node, node *html.Node) {
  parent.AppendChild(node)
}

func (vdom *Vdom) CreateElement(tagName string) *html.Node {
  node := html.Node{Type: html.ElementNode, Data: tagName}
  return &node
}


/*
func (vdom *Vdom) SetAttributes(t, attrs map[string] string) *html.Node {
  node := html.Node{Type: html.ElementNode, Data: tagName}
  return &node
}
*/

type Settings struct {
  IdAttribute string 
}

func GetSettings() Settings {
  return Settings{IdAttribute: "data-id"}
}

func GetIntSlice(id string) ([]int, error) {
  var out []int
  parts := strings.Split(id, ".")
  for _, num := range parts {
    s, err := strconv.Atoi(num)
    if err != nil {
      return nil, err
    }
    out = append(out, s)
  }
  return out, nil
}

func GetIdentifier(ids []int) string {
  id := ""
  for index, num := range ids {
    if index > 0 {
      id += "."
    }
    id += strconv.Itoa(num)
  }
  return id
}

func Parse(b []byte, settings Settings) (*Vdom, error) {
  r := bytes.NewBuffer(b)
  doc, err := html.Parse(r)
  if err != nil {
    return nil, err
  }

  dom := Vdom{Document: doc, Map: make(map[string] *html.Node)}
  var f func(n *html.Node, ids []int)
  f = func(n *html.Node, ids []int) {
    var id string
    var i int = 0
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      if c.Type == html.ElementNode {
        list := append(ids, i)
        id = GetIdentifier(list)
        //log.Printf("id: %s", id)
        mock, err := GetIntSlice(id)
        if err != nil {
          log.Fatal(err)
        }
        log.Println(mock)
        log.Println(len(mock))
        dom.Map[id] = c
        attr := html.Attribute{Key: settings.IdAttribute, Val: id}
        c.Attr = append(c.Attr, attr)
        f(c, list)
        i++
      }
    }
  }

  var ids []int
  f(doc, ids)
  return &dom, nil
}
