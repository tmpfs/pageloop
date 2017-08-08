package vdom

import(
  "errors"
  "golang.org/x/net/html"
)

// Apply a patch to the DOM.
//
// Returns a patch that reverses the operations. If an error occurs 
// you should try to apply the returned patch to rollback to the previous
// state. If no error is returned you can use the patch to undo the operation.
//
// For the memento undo pattern patches must be applied in the correct reverse 
// order.
func (vdom *Vdom) Apply(patch *Patch) (Patch, error) {
  var err error
  var out = Patch{}
  var tx *Diff

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
          tx, err = vdom.RemoveDiff(parent, n)
          if err != nil {
            return out, err
          }
          out.Add(tx)
          err = vdom.AppendChild(parent, n)
          tx.Element = vdom.GetId(n)
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
          tx, err = vdom.RemoveDiff(parent, n)
          if err != nil {
            return out, err
          }
          out.Add(tx)
          err = vdom.InsertBefore(parent, n, target)
          tx.Element = vdom.GetId(n)
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

        if target.NextSibling == nil {
          tx, err = vdom.AppendDiff(parent, target)
          if err != nil {
            return out, err
          }
        } else {
          tx, err = vdom.InsertDiff(parent, target, target.NextSibling)
          if err != nil {
            return out, err
          }
        }

        tx.Element = vdom.GetId(parent)
        out.Add(tx)

        err = vdom.RemoveChild(parent, target)
        if err != nil {
          return out, err
        }
      // TODO: transactions for attributes
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
