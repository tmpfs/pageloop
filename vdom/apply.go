package vdom

import(
  "errors"
  "golang.org/x/net/html"
)

// Apply a patch to the DOM.
func (vdom *Vdom) Apply(patch *Patch) (Patch, error) {
  var err error
  var out = Patch{}

  // parse an HTML fragment from the diff data
  var getNodes func (diff *Diff) ([]*html.Node, error)
  getNodes = func(diff *Diff) ([]*html.Node, error) {
    var nodes []*html.Node
    var err error
    switch diff.Type {
      case html.ElementNode:
        // parse the node data
        nodes, err = vdom.ParseFragment(diff.Data, nil)
        if err != nil {
          return nodes, err
        }
      case html.DoctypeNode:
        fallthrough
      case html.TextNode:
        fallthrough
      case html.CommentNode:
        fallthrough
      default:
        var node *html.Node = &html.Node{Type: diff.Type, Data: string(diff.Data)}
        nodes = append(nodes, node)
    }

    return nodes, err
  }

  // iterate and attempt to apply operations
  for _, diff := range patch.Diffs {
    switch diff.Operation {
      case APPEND:
        var parent *html.Node = vdom.Map[diff.Element]
        if parent == nil {
          return out, errors.New("Missing parent node for append operation")
        }

        var nodes []*html.Node
        var err error
        
        nodes, err = getNodes(&diff)
        if err != nil {
          return out, err
        }
        for _, n := range nodes {
          err = vdom.AppendChild(parent, n)
          if err != nil {
            return out, err
          }
        }
      case INSERT:
        var target *html.Node = vdom.Map[diff.Element]
        if target == nil {
          return out, errors.New("Missing target node for insert before operation")
        }

        // infer the parent node
        var parent *html.Node = target.Parent
        if parent == nil {
          return out, errors.New("Missing parent node for insert before operation (node may be detached)")
        }

        var nodes []*html.Node
        var err error
        
        nodes, err = getNodes(&diff)
        if err != nil {
          return out, err
        }
        for _, n := range nodes {
          err = vdom.InsertBefore(parent, n, target)
          if err != nil {
            return out, err
          }
        }
      case REMOVE:
        var target *html.Node = vdom.Map[diff.Element]
        if target == nil {
          return out, errors.New("Missing target node for remove operation")
        }

        // infer the parent node
        var parent *html.Node = target.Parent
        if parent == nil {
          return out, errors.New("Missing parent node for remove operation (node may be detached)")
        }

        err = vdom.RemoveChild(parent, target)
        if err != nil {
          return out, err
        }
      case ATTR_SET:
        var target *html.Node = vdom.Map[diff.Element]
        if target == nil {
          return out, errors.New("Missing target node for set attribute operation")
        }
        vdom.SetAttr(target, diff.Attr)
      case ATTR_DEL:
        var target *html.Node = vdom.Map[diff.Element]
        if target == nil {
          return out, errors.New("Missing target node for delete attribute operation")
        }
        vdom.DelAttr(target, diff.Attr)
    }
  }

  return out, nil
}
