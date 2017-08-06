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

func FindLastChildElement(parent *html.Node) *html.Node {
  for c := parent.LastChild; c != nil; c = c.PrevSibling {
    if c.Type == html.ElementNode {
      return c
    }
  }
  return nil
}

func (vdom *Vdom) AppendChild(parent *html.Node, node *html.Node) error {
  var ids []int
  var err error
  if parent.LastChild != nil {
    log.Println(parent.LastChild)
    ids, err = FindId(FindLastChildElement(parent))
    // increment for the new id
    ids[len(ids) - 1]++
  } else {
    // TODO: test adding to empty parent
    ids, err = FindId(parent)
    // now the new first child
    ids = append(ids, 0)
  }
  if err != nil {
    return err
  }
  id := GetIdentifier(ids)
  log.Println("ID: ", id)
  //node.Attr = append(node.Attr, html.Attribute{Key: idAttribute, Val: id})
  vdom.SetAttr(node, html.Attribute{Key: idAttribute, Val: id})
  parent.AppendChild(node)
  log.Println(parent.LastChild)
  return err
}

func (vdom *Vdom) CreateElement(tagName string) *html.Node {
  node := html.Node{Type: html.ElementNode, Data: tagName}
  return &node
}

func (vdom *Vdom) GetAttr(node *html.Node, key string) *html.Attribute {
  for _, attr := range node.Attr {
    if attr.Key == key {
      return &attr
    }
  }
  return nil
}

func (vdom *Vdom) SetAttr(node *html.Node, attr html.Attribute) {
  existing := vdom.GetAttr(node, attr.Key)
  if existing != nil {
    existing.Val = attr.Val 
  } else {
    node.Attr = append(node.Attr, attr)
  }
}

var idAttribute string  = "data-id"

func FindId(node *html.Node) ([]int, error) {
  var id string
  for _, attr := range node.Attr {
    log.Println(attr)
    if attr.Key == idAttribute {
      id = attr.Val
    }
  }
  log.Println("id: ", id)
  return GetIntSlice(id)
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

func Parse(b []byte) (*Vdom, error) {
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
        attr := html.Attribute{Key: idAttribute, Val: id}
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
