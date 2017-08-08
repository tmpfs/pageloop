// Virtual DOM implementation wrapping the parse tree 
// returned from x/net/html.
//
// Able to diff and patch using a collection of diff operations.
package vdom
  
import(
  //"os"
  //"io"
  //"fmt"
  //"log"
  "bytes"
  "strconv"
  "strings"
  "golang.org/x/net/html"
)

var idAttribute string  = "data-id"

// The virtual DOM.
type Vdom struct {
  Document *html.Node
  Map map[string] *html.Node
}

// Basic DOM API wrapper functions

// Append a child node.
func (vdom *Vdom) AppendChild(parent *html.Node, node *html.Node) error {
  var ids []int
  var err error
  if parent.LastChild != nil {
    ids, err = vdom.FindId(findLastChildElement(parent))
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
  id := intSliceToString(ids)
  vdom.SetAttr(node, html.Attribute{Key: idAttribute, Val: id})
  parent.AppendChild(node)
  vdom.Map[id] = node
  return err
}

// Insert a child node before another node.
func (vdom *Vdom) InsertBefore(parent *html.Node, newChild *html.Node, oldChild * html.Node) error {
  var err error
  //_, id := vdom.GetAttr(oldChild, idAttribute)
  id := vdom.GetId(oldChild)
  vdom.SetAttr(newChild, html.Attribute{Key: idAttribute, Val: id})
  parent.InsertBefore(newChild, oldChild)
  vdom.Map[id] = newChild
  err = vdom.adjustSiblings(newChild, true)
  return err
}

// Remove a node.
func (vdom *Vdom) RemoveChild(parent *html.Node, node *html.Node) error {
  var err error = vdom.adjustSiblings(node, false)
  if err != nil {
    return err
  }

  id := vdom.GetId(node)
  delete(vdom.Map, id)
  parent.RemoveChild(node)
  return err
}

// Diff / Patch functions

// Append a child node and return a diff that represents the operation.
func (vdom *Vdom) AppendDiff(parent *html.Node, node *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: APPEND_OP, Element: vdom.GetId(parent)}

  // convert to byte slice
  op.Data, err = renderToBytes(node)
  if err != nil {
    return nil, err
  }

  // append the child
  err = vdom.AppendChild(parent, node)
  if err != nil {
    return nil, err
  }

  return &op, err
}

// Insert a child node before another node and return a diff that represents the operation.
func (vdom *Vdom) InsertDiff(parent *html.Node, newChild *html.Node, oldChild *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: INSERT_OP, Element: vdom.GetId(oldChild)}

  // convert to byte slice
  op.Data, err = renderToBytes(newChild)
  if err != nil {
    return nil, err
  }

  err = vdom.InsertBefore(parent, newChild, oldChild)
  if err != nil {
    return nil, err
  }

  return &op, err
}

// Remove a node and return a diff that represents the operation.
func (vdom *Vdom) RemoveDiff(parent *html.Node, node *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: REMOVE_OP, Element: vdom.GetId(node)}

  // convert to byte slice
  op.Data, err = renderToBytes(node)

  err = vdom.RemoveChild(parent, node)
  if err != nil {
    return nil, err
  }

  return &op, err
}

func renderToBytes(node *html.Node) ([]byte, error) {
  w := new(bytes.Buffer)
  err := html.Render(w, node)
  if err != nil {
    return nil, err
  }
  return w.Bytes(), nil
}

// Extensions to the basic API

// Create a new element.
func (vdom *Vdom) CreateElement(tagName string) *html.Node {
  node := html.Node{Type: html.ElementNode, Data: tagName}
  return &node
}

// Get an Attribute.
//
// Returns the index of the attribute and a pointer to 
// the attribute.
func (vdom *Vdom) GetAttr(node *html.Node, key string) (int, *html.Attribute) {
  for index, attr := range node.Attr {
    if attr.Key == key {
      return index, &attr
    }
  }
  return -1, nil
}

// Get the value of an attribute, that is `html.Attribute.Val`.
func (vdom *Vdom) GetAttrValue(node *html.Node, key string) string {
  _, attr := vdom.GetAttr(node, key)
  if attr != nil {
    return attr.Val
  }
  return ""
}

// Set an Attribute.
func (vdom *Vdom) SetAttr(node *html.Node, attr html.Attribute) {
  index, existing := vdom.GetAttr(node, attr.Key)
  if existing != nil {
    existing.Val = attr.Val
    node.Attr[index] = *existing
  } else {
    node.Attr = append(node.Attr, attr)
  }
}

// Get the vdom identifier for an element extracted from the `data-id` attribute.
func (vdom *Vdom) GetId(node *html.Node) string {
  return vdom.GetAttrValue(node, idAttribute)
}

// Get the identifier for a node as a slice of integers.
func (vdom *Vdom) FindId(node *html.Node) ([]int, error) {
  id := vdom.GetId(node)
  return stringToIntSlice(id)
}

// Private vdom methods

// Increments or decrements the identifiers for siblings that appear
// after the target node. Used when modifying the DOM to keep identifiers 
// sequential.
func (vdom *Vdom) adjustSiblings(node *html.Node, increment bool) error {
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
    newId := intSliceToString(ids)
    vdom.SetAttr(c, html.Attribute{Key:idAttribute, Val: newId})
    delete(vdom.Map, oldId)
    vdom.Map[newId] = c
  }
  return nil
}


// Parse an HTML document assigning virtual dom identifiers.
// Assigns each element a `data-id` attribute and adds entries 
// to the vdom `Map` for fast node lookup.
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
        id = intSliceToString(list)
        dom.Map[id] = c
        dom.SetAttr(c, html.Attribute{Key: idAttribute, Val: id})
        f(c, list)
        i++
      }
    }
  }

  var ids []int
  f(doc, ids)
  return &dom, nil
}

// Helper functions

// Splits a string id to a slice of integers.
func stringToIntSlice(id string) ([]int, error) {
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

// Converts a slice of integers to a string.
func intSliceToString(ids []int) string {
  id := ""
  for index, num := range ids {
    if index > 0 {
      id += "."
    }
    id += strconv.Itoa(num)
  }
  return id
}

// Find the previous sibling that is of type element 
// starting with the last child of the parent.
func findLastChildElement(parent *html.Node) *html.Node {
  for c := parent.LastChild; c != nil; c = c.PrevSibling {
    if c.Type == html.ElementNode {
      return c
    }
  }
  return nil
}

