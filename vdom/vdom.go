// Virtual DOM implementation backed by the parse tree 
// returned from golang.org/x/net/html.
package vdom
  
import(
  //"log"
  "bytes"
  "strconv"
  "strings"
  "golang.org/x/net/html"
  "golang.org/x/net/html/atom"
)

var idAttribute string  = "data-id"

// The virtual DOM.
type Vdom struct {
  Document *html.Node
  Map map[string] *html.Node
}

// Parse an HTML document assigning virtual dom identifiers.
// Assigns each element an identifier attribute and adds entries 
// to the vdom Map for fast node lookup.
func (vdom *Vdom) Parse(b []byte) error {
  r := bytes.NewBuffer(b)
  doc, err := html.Parse(r)
  if err != nil {
    return err
  }

  //vdom := Vdom{Document: doc, Map: make(map[string] *html.Node)}
  vdom.Document = doc
  vdom.Map = make(map[string] *html.Node)
  var f func(n *html.Node, ids []int)
  f = func(n *html.Node, ids []int) {
    var id string
    var i int = 0
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      if c.Type == html.ElementNode {
        list := append(ids, i)
        id = intSliceToString(list)
        vdom.Map[id] = c
        vdom.SetAttr(c, html.Attribute{Key: idAttribute, Val: id})
        f(c, list)
        i++
      }
    }
  }

  var ids []int
  f(doc, ids)
  return nil
}

// Basic DOM API wrapper functions

// Append a child node.
func (vdom *Vdom) AppendChild(parent *html.Node, node *html.Node) error {
  var ids []int
  var err error
  var hasPreviousSibling bool
  // get identifier from last child element node
  if parent.LastChild != nil {
    el := findLastChildElement(parent)
    // might have children but no element nodes!
    if el != nil {
      ids, err = vdom.FindId(el)
      if err != nil {
        return err
      }
      // increment for the new id
      ids[len(ids) - 1]++
      hasPreviousSibling = true
    }
  }

  // get identifier from parent and start child index at zero
  if parent.LastChild == nil || !hasPreviousSibling {
    ids, err = vdom.FindId(parent)
    // now the new first child
    ids = append(ids, 0)

    if err != nil {
      return err
    }
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

// Extensions to the basic API

// Parse an HTML fragment and return the list of parsed nodes.
func (vdom *Vdom) ParseFragment(b []byte, context *html.Node) ([]*html.Node, error) {
  if context == nil {
    context = &html.Node{Type: html.ElementNode, Data: "body", DataAtom:atom.Body}
  }
  r := bytes.NewBuffer(b)
  nodes, err := html.ParseFragment(r, context)
  if err != nil {
    return nil, err
  }

  return nodes, err
}

// Create a new element.
func (vdom *Vdom) CreateElement(tagName string) *html.Node {
  node := html.Node{Type: html.ElementNode, Data: tagName}
  return &node
}

// Create a text node.
func (vdom *Vdom) CreateTextNode(text string) *html.Node {
  node := html.Node{Type: html.TextNode, Data: text}
  return &node
}

// Create a comment node.
func (vdom *Vdom) CreateCommentNode(comment string) *html.Node {
  node := html.Node{Type: html.CommentNode, Data: comment}
  return &node
}

// Create a doctype node.
func (vdom *Vdom) CreateDoctypeNode(doctype string) *html.Node {
  node := html.Node{Type: html.DoctypeNode, Data: doctype}
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
//
// If the attribute exists it is overwritten otherwise it is created.
func (vdom *Vdom) SetAttr(node *html.Node, attr html.Attribute) {
  index, existing := vdom.GetAttr(node, attr.Key)
  if existing != nil {
    existing.Val = attr.Val
    node.Attr[index] = *existing
  } else {
    node.Attr = append(node.Attr, attr)
  }
}

// Remove an attribute
func (vdom *Vdom) DelAttr(node *html.Node, attr html.Attribute) {
  for index, a := range node.Attr {
    var match bool = a.Key == attr.Key
    if a.Namespace != "" && attr.Namespace != "" {
      match = a.Namespace == attr.Namespace
    }

    if match {
      before := node.Attr[0:index]
      after := node.Attr[index + 1:len(node.Attr)]
      node.Attr = node.Attr[0:0]
      node.Attr = append(node.Attr, before...)
      node.Attr = append(node.Attr, after...)
      break
    }
  }
}

// Internal helper functions.

// Get the unique vdom identifier for an element extracted from the identifier attribute.
func (vdom *Vdom) GetId(node *html.Node) string {
  return vdom.GetAttrValue(node, idAttribute)
}

// Get the identifier for a node as a slice of integers.
func (vdom *Vdom) FindId(node *html.Node) ([]int, error) {
  id := vdom.GetId(node)
  return stringToIntSlice(id)
}

// Render a node to a byte slice, typically for debugging.
func (vdom *Vdom) RenderToBytes(node *html.Node) ([]byte, error) {
  w := new(bytes.Buffer)
  err := html.Render(w, node)
  if err != nil {
    return nil, err
  }
  return w.Bytes(), nil
}

// Diff / Patch functions

// Append a child node and return a diff that represents the operation.
func (vdom *Vdom) AppendDiff(parent *html.Node, node *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: APPEND, Element: vdom.GetId(parent), Type: node.Type}
  // convert to byte slice
  op.Data, err = vdom.RenderToBytes(node)
  return &op, err
}

// Insert a child node before another node and return a diff that represents the operation.
func (vdom *Vdom) InsertDiff(parent *html.Node, newChild *html.Node, oldChild *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: INSERT, Element: vdom.GetId(oldChild), Type: newChild.Type}

  // convert to byte slice
  op.Data, err = vdom.RenderToBytes(newChild)
  return &op, err
}

// Remove a node and return a diff that represents the operation.
func (vdom *Vdom) RemoveDiff(parent *html.Node, node *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: REMOVE, Element: vdom.GetId(node), Type: node.Type}

  // convert to byte slice
  op.Data, err = vdom.RenderToBytes(node)
  return &op, err
}

// Set an attribute and return a diff that represents the operation.
func (vdom *Vdom) SetAttrDiff(node *html.Node, attr html.Attribute) (*Diff, error) {
  var op Diff = Diff{Operation: ATTR_SET, Element: vdom.GetId(node), Attr: attr, Type: node.Type}
  return &op, nil
}

// Delete an attribute and return a diff that represents the operation.
func (vdom *Vdom) DelAttrDiff(node *html.Node, attr html.Attribute) (*Diff, error) {
  var op Diff = Diff{Operation: ATTR_DEL, Element: vdom.GetId(node), Attr: attr, Type: node.Type}
  return &op, nil
}

/*
// Set the text for a node.
func (vdom *Vdom) Text(parent *html.Node, node *html.Node, text []byte) error {
  node.Data = text
  return nil
}
*/

/*
func (vdom *Vdom) TextDiff(parent *html.Node) {

}
*/


// Private vdom methods

// Increments or decrements the identifiers for siblings that appear
// after the target node. Used when modifying the DOM to keep identifiers 
// sequential.
func (vdom *Vdom) adjustSiblings(node *html.Node, increment bool) error {
  for c := node.NextSibling; c != nil; c = c.NextSibling {
    //oldId := vdom.GetAttrValue(c, idAttribute)
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
    //delete(vdom.Map, oldId)
    vdom.Map[newId] = c
  }
  return nil
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

