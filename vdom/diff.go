package vdom
  
import(
  "golang.org/x/net/html"
)

const (
  APPEND = iota
  INSERT
  REMOVE
  ATTR_SET
  ATTR_DEL
  TEXT
)

// Diff is an operation on the virtual DOM that results 
// in a change to the DOM.
//
// Whilst it represents a difference in the DOM this is 
// essentially the memento pattern, storing changes so 
// that they may be applied and reverted.
type Diff struct {
  // Operation constant
  Operation int

  // Id of the primary target element.
  // 
  // For the append operation it is the parent to append to.
  // For the insert operation it is the old child node (parent is inferred).
  // For the remove operation it is the node to remove.
  // For the attr operation it is the target node.
  // For the text operation it is the parent element.
  Element string

  // A node type associated with the data.
  Type html.NodeType

  // HTML fragment data (append and insert only) or data for the text operation.
  //
  // The remove operation may propagate this with the node being removed so 
  // that the operation can be reversed.
  Data []byte

  // For the text operation an index into the element's child nodes to use 
  // to set the text.
  Index int

  // A key value pair when setting attributes.
  Attr html.Attribute
}

// Patch is a slice of diff operations.
type Patch struct {
  Diffs []Diff
}

// Add a diff to the patch, returns the length of the diff slice.
func (p *Patch) Add(diff *Diff) int {
  p.Diffs = append(p.Diffs, *diff)
  return len(p.Diffs)
}
