// Virtual DOM implementation backed by the parse tree
// returned from golang.org/x/net/html.
//
// When a document is parsed each element node is given a
// unique identifier in the form 0.1.2 where each integer
// represents the child index. Typically the html element
// will be index 0.
//
// So that identifiers may be kept in sync you should call
// the wrapper API methods AppendChild, InsertBefore and
// RemoveChild. They each require an additional first
// argument which is the parent node to change.
package vdom

import(
  //"log"
  "bytes"
  "regexp"
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
        vdom.AttrSet(c, html.Attribute{Key: idAttribute, Val: id})
        f(c, list)
        i++
      }
    }
  }

  var ids []int
  f(doc, ids)
  return nil
}

// Clones a node and removes any attributes associated with
// the virtual DOM from elements in the tree.
//
// If the given node is nil then the document node is used.
func (vdom *Vdom) Clean(node *html.Node) *html.Node {
  if node == nil {
    node = vdom.Document
  }
  node = vdom.CloneNode(node, true)
  var f func(n *html.Node)
  f = func(n *html.Node) {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      if c.Type == html.ElementNode {
        vdom.AttrDel(c, html.Attribute{Key: idAttribute})
        f(c)
      }
    }
  }
  f(node)
  return node
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
  vdom.AttrSet(node, html.Attribute{Key: idAttribute, Val: id})
  parent.AppendChild(node)
  vdom.Map[id] = node
  return err
}

// Insert a child node before another node.
func (vdom *Vdom) InsertBefore(parent *html.Node, newChild *html.Node, oldChild * html.Node) error {
  var err error
  id := vdom.GetId(oldChild)
  vdom.AttrSet(newChild, html.Attribute{Key: idAttribute, Val: id})
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

// Clone a node, if the deep option is given all descendants are
// also cloned.
func (vdom *Vdom) CloneNode(node *html.Node, deep bool) *html.Node {
  if node == nil {
    node = vdom.Document
  }

  var clone func(n *html.Node) *html.Node
  clone = func(n *html.Node) *html.Node {
    var out = &html.Node{Type: n.Type, Data: n.Data[0:], DataAtom: n.DataAtom, Namespace: n.Namespace}
    for _, att := range n.Attr {
      out.Attr = append(out.Attr, html.Attribute{Key: att.Key, Val: att.Val, Namespace: att.Namespace})
    }
    return out
  }
  out := clone(node)
  var f func(n *html.Node, p *html.Node)
  f = func(n *html.Node, p *html.Node) {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      copy := clone(c)
      p.AppendChild(copy)
      if deep {
        f(c, copy)
      }
    }
  }
  f(node, out)
  return out
}

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
func (vdom *Vdom) AttrGet(node *html.Node, key string) (int, *html.Attribute) {
  for index, attr := range node.Attr {
    if attr.Key == key {
      return index, &attr
    }
  }
  return -1, nil
}

func (vdom *Vdom) AttrGetNs(node *html.Node, key string, ns string) (int, *html.Attribute) {
  for index, attr := range node.Attr {
    if attr.Key == key && attr.Namespace == ns {
      return index, &attr
    }
  }
  return -1, nil
}

// Get the value of an attribute.
func (vdom *Vdom) AttrGetValue(node *html.Node, key string) string {
  _, attr := vdom.AttrGet(node, key)
  if attr != nil {
    return attr.Val
  }
  return ""
}

// Set an Attribute.
//
// If the attribute exists it is overwritten otherwise it is created.
func (vdom *Vdom) AttrSet(node *html.Node, attr html.Attribute) {
  index, existing := vdom.AttrGet(node, attr.Key)
  if existing != nil {
    existing.Val = attr.Val
    node.Attr[index] = *existing
  } else {
    node.Attr = append(node.Attr, attr)
  }
}

// Remove an attribute
func (vdom *Vdom) AttrDel(node *html.Node, attr html.Attribute) {
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
  return vdom.AttrGetValue(node, idAttribute)
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

// Render a node unescaped to a byte slice, typically for debugging.
func (vdom *Vdom) RenderRaw(node *html.Node) ([]byte, error) {
  w := new(bytes.Buffer)
  err := RenderUnsafe(w, node)
  if err != nil {
    return nil, err
  }
  return w.Bytes(), nil
}


// Render a node to a byte slice but do not perform HTML escaping.
/*
func (vdom *Vdom) RenderToBytesUnsafe(node *html.Node) ([]byte, error) {
  w := new(bytes.Buffer)

	var render func(node *html.Node)
	render = func(node *html.Node) {
		println("render: " + node.Data)
		switch(node.Type) {
			case html.DoctypeNode:
				w.Write([]byte(`<!doctype ` + node.Data + `>`))
			case html.TextNode:
				w.Write([]byte(node.Data))
			case html.CommentNode:
				w.Write([]byte(`<!--` + node.Data + `-->`))
			case html.DocumentNode:
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					render(c)
				}
			case html.ElementNode:
				w.Write([]byte(`<` + node.Data))
				for _, att := range node.Attr {
					w.Write([]byte(` ` + att.Key + `='` + att.Val + `'`))
				}
				w.Write([]byte(`>`))
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					render(c)
				}
				w.Write([]byte(`</` + node.Data + `>`))
		}
	}

	render(node)

  return w.Bytes(), nil
}
*/

// Compacts a DOM tree by removing text nodes between
// element nodes that consist solely of whitespace.
//
// Removes whitespace from descendants of the target node
// and returns the target node with in place modifications
// if you need to keep the original you should copy it first.
//
// The text nodes are preserved in the tree; their Data is
// set to the empty string.
func (vdom *Vdom) Compact(node *html.Node) *html.Node {
  var whitespace = regexp.MustCompile(`^\s+$`)
  for c := node.FirstChild; c != nil; c = c.NextSibling {
    var next *html.Node = c.NextSibling
    var prev *html.Node = c.PrevSibling
    if (c.Type == html.TextNode && whitespace.MatchString(c.Data)) {
      var removes bool = (next == nil || next.Type == html.ElementNode) && (prev == nil || prev.Type == html.ElementNode)
      if removes {
        c.Data = ""
        //node.RemoveChild(c)
      }
    }
    if (c.Type == html.ElementNode) {
      vdom.Compact(c)
    }
  }
  return node
}

// Diff / Patch functions

// Get a diff that represents the append child operation.
func (vdom *Vdom) DiffAppend(parent *html.Node, node *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: APPEND, Id: vdom.GetId(parent), Element: parent, Type: node.Type}
  // convert to byte slice
  op.Data, err = vdom.RenderToBytes(node)
  return &op, err
}

// Get a diff that represents the insert before operation.
func (vdom *Vdom) DiffInsert(parent *html.Node, newChild *html.Node, oldChild *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: INSERT, Id: vdom.GetId(oldChild), Element: oldChild, Type: newChild.Type}

  // convert to byte slice
  op.Data, err = vdom.RenderToBytes(newChild)
  return &op, err
}

// Get a diff that represents the remove child operation.
func (vdom *Vdom) DiffRemove(parent *html.Node, node *html.Node) (*Diff, error) {
  var err error
  var op Diff = Diff{Operation: REMOVE, Id: vdom.GetId(node), Element: node, Type: node.Type}

  // convert to byte slice
  op.Data, err = vdom.RenderToBytes(node)
  return &op, err
}

// Get a diff that represents the set attribute operation.
func (vdom *Vdom) DiffAttrSet(node *html.Node, attr html.Attribute) (*Diff, error) {
  var op Diff = Diff{Operation: ATTR_SET, Id: vdom.GetId(node), Element: node, Attr: attr, Type: node.Type}
  return &op, nil
}

// Get a diff that represents the delete attribute operation.
func (vdom *Vdom) DiffAttrDel(node *html.Node, attr html.Attribute) (*Diff, error) {
  var op Diff = Diff{Operation: ATTR_DEL, Id: vdom.GetId(node), Element: node, Attr: attr, Type: node.Type}
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
    vdom.AttrSet(c, html.Attribute{Key:idAttribute, Val: newId})
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

