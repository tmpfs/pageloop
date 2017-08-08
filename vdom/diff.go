package vdom
  
import(
  "encoding/json"
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

  // A key value pair when setting attributes.
  Attr html.Attribute

  // For the text operation an index into the element's child nodes to use 
  // to set the text.
  //Index int
}

// Encode the diff to JSON.
func (diff Diff) MarshalJSON() ([]byte, error) {
  var o map[string] interface{} = make(map[string] interface{})
  o["op"] = diff.Operation
  o["id"] = diff.Element
  o["type"] = diff.Type
  o["data"] = string(diff.Data)
  if diff.Attr.Key != "" {
    var a map[string] interface{} = make(map[string] interface{})
    a["key"] = diff.Attr.Key
    if diff.Attr.Namespace != "" {
      a["ns"] = diff.Attr.Namespace
    }
    if diff.Attr.Val != "" {
      a["val"] = diff.Attr.Val
    }
    o["attr"] = a
  }
  json, err := json.Marshal(&o)
  if err != nil {
    return nil, err
  }
  return json, nil
}

// Decode JSON into the diff.
func (diff *Diff) UnmarshalJSON(b []byte) error {
  //var temp *Diff
  var temp map[string] interface{} = make(map[string] interface{})
  if err := json.Unmarshal(b, &temp); err != nil {
    return err
  }

  diff.Operation = int(temp["op"].(float64))
  diff.Element = temp["id"].(string)
  diff.Type = html.NodeType(temp["type"].(float64))
  diff.Data = []byte(temp["data"].(string))

  if temp["attr"] != nil {
    var a map[string] interface{} = temp["attr"].(map[string] interface{})
    attr := html.Attribute{}
    attr.Key = a["key"].(string)
    if a["ns"] != nil {
      attr.Namespace = a["ns"].(string)
    }
    if a["val"] != nil {
      attr.Val = a["val"].(string)
    }

    diff.Attr = attr
  }

  return nil
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

func (p Patch) MarshalJSON() ([]byte, error) {
  return json.Marshal(&p.Diffs)
}

func (p *Patch) UnmarshalJSON(b []byte) error {
  return json.Unmarshal(b, &p.Diffs)
}
