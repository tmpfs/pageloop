package vdom

import(
  //"log"
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
        var parent *html.Node = vdom.Map[diff.Id]
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
          err = vdom.AppendChild(parent, n)
          tx.Id = vdom.GetId(n)
          tx.Element = n
          out.Add(tx)
          if err != nil {
            return out, err
          }
        }
      case INSERT:
        var target *html.Node = vdom.Map[diff.Id]
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
          err = vdom.InsertBefore(parent, n, target)
          tx.Id = vdom.GetId(n)
          tx.Element = n
          out.Add(tx)
          if err != nil {
            return out, err
          }
        }
      case REMOVE:
        var target *html.Node = vdom.Map[diff.Id]
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

        err = vdom.RemoveChild(parent, target)
        tx.Id = vdom.GetId(parent)
        tx.Element = target
        out.Add(tx)
        if err != nil {
          return out, err
        }
      case ATTR_SET:
        var target *html.Node = vdom.Map[diff.Id]
        if target == nil {
          return out, errors.New("Missing target node for set attribute operation")
        }
        _, att := vdom.GetAttrNs(target, diff.Attr.Key, diff.Attr.Namespace)
        // revert to previous attribute value
        if att != nil {
          tx = &Diff{Operation: ATTR_SET, Id: vdom.GetId(target), Element: target, Attr: *att, Type: target.Type}
        // or delete if it didn't exist
        } else {
          tx = &Diff{Operation: ATTR_DEL, Id: vdom.GetId(target), Element: target, Attr: diff.Attr, Type: target.Type}
        }

        out.Add(tx)

        vdom.SetAttr(target, diff.Attr)
      case ATTR_DEL:
        var target *html.Node = vdom.Map[diff.Id]
        if target == nil {
          return out, errors.New("Missing target node for delete attribute operation")
        }
        _, att := vdom.GetAttrNs(target, diff.Attr.Key, diff.Attr.Namespace)
        // revert to previous attribute value
        if att != nil {
          tx = &Diff{Operation: ATTR_SET, Id: vdom.GetId(target), Element: target, Attr: *att, Type: target.Type}
        // or delete if it didn't exist
        } else {
          tx = &Diff{Operation: ATTR_DEL, Id: vdom.GetId(target), Element: target, Attr: diff.Attr, Type: target.Type}
        }

        out.Add(tx)

        vdom.DelAttr(target, diff.Attr)
    }
  }

  var reverse []Diff
  for i := len(out.Diffs) - 1; i >= 0; i-- {
    txn := out.Diffs[i] 
    reverse = append(reverse, txn)
  }
  out.Diffs = reverse

  return out, nil
}
