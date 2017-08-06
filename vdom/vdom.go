package vdom
  
import(
  //"log"
  "bytes"
  "strconv"
  "strings"
  "golang.org/x/net/html"
)

const (
  APPEND_OP = iota
  INSERT_OP
  REMOVE_OP
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

func FindPrevSiblingElement(parent *html.Node) *html.Node {
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
    ids, err = vdom.FindId(FindLastChildElement(parent))
    if err != nil {
      return err
    }
    // increment for the new id
    ids[len(ids) - 1]++
  } else {
    // TODO: test adding to empty parent
    ids, err = vdom.FindId(parent)
    // now the new first child
    ids = append(ids, 0)
  }
  if err != nil {
    return err
  }
  id := GetIdentifier(ids)
  vdom.SetAttr(node, html.Attribute{Key: idAttribute, Val: id})
  parent.AppendChild(node)
  vdom.Map[id] = node
  return err
}

func (vdom *Vdom) InsertBefore(parent *html.Node, newChild *html.Node, oldChild * html.Node) error {
  var err error
  _, id := vdom.GetAttr(oldChild, idAttribute)
  vdom.SetAttr(newChild, html.Attribute{Key: idAttribute, Val: id.Val})
  parent.InsertBefore(newChild, oldChild)
  //vdom.Map[id] = newChild
  err = vdom.AdjustSiblings(newChild, true)
  return err
}

func (vdom *Vdom) AdjustSiblings(node *html.Node, increment bool) error {
  for c := node.NextSibling; c != nil; c = c.NextSibling {
    oldId := vdom.GetAttrValue(c, idAttribute)
    ids, err := vdom.FindId(c)
    if err != nil {
      return err
    }
    if increment {
      ids[len(ids) - 1]++
    } else {
      ids[len(ids) - 1]--
    }
    newId := GetIdentifier(ids)
    vdom.SetAttr(c, html.Attribute{Key:idAttribute, Val: newId})
    delete(vdom.Map, oldId)
    vdom.Map[newId] = c
  }
  return nil
}

func (vdom *Vdom) RemoveChild(parent *html.Node, node *html.Node) error {
  var err error = vdom.AdjustSiblings(node, false)
  if err != nil {
    return err
  }

  _, id := vdom.GetAttr(node, idAttribute)
  delete(vdom.Map, id.Val)
  parent.RemoveChild(node)
  return err
}

func (vdom *Vdom) CreateElement(tagName string) *html.Node {
  node := html.Node{Type: html.ElementNode, Data: tagName}
  return &node
}

func (vdom *Vdom) GetAttrValue(node *html.Node, key string) string {
  _, attr := vdom.GetAttr(node, key)
  if attr != nil {
    return attr.Val
  }
  return ""
}

func (vdom *Vdom) GetAttr(node *html.Node, key string) (int, *html.Attribute) {
  for index, attr := range node.Attr {
    if attr.Key == key {
      return index, &attr
    }
  }
  return -1, nil
}

func (vdom *Vdom) SetAttr(node *html.Node, attr html.Attribute) {
  index, existing := vdom.GetAttr(node, attr.Key)
  if existing != nil {
    existing.Val = attr.Val
    node.Attr[index] = *existing
  } else {
    node.Attr = append(node.Attr, attr)
  }
}

var idAttribute string  = "data-id"

func (vdom *Vdom) FindId(node *html.Node) ([]int, error) {
  _, attr := vdom.GetAttr(node, idAttribute)
  id := attr.Val
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
