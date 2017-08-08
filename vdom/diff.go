package vdom
  
import(
  //"log"
  //"bytes"
  //"strconv"
  //"strings"
  "golang.org/x/net/html"
)

const (
  APPEND_OP = iota
  INSERT_OP
  REMOVE_OP
  ATTR_SET_OP
  ATTR_DEL_OP
  TEXT_OP
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
