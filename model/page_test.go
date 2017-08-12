package model

import (
  //"os"
  "testing"
  "encoding/json"
)

func TestBlock(t *testing.T) {
  b := Block{Title: "mock block"}

  b.Content = "<html>"

  if b.Title == "" {
    t.Error("Unexpected empty title")
  }

  if b.Content == "" {
    t.Error("Unexpected empty content")
  }

  p := Page{}
  p.AddBlock(b)

  if p.Length() != 1 {
    t.Errorf("Unexpected block length %d expected %d after add", p.Length(), 1)
  }

  /*
  p.RemoveBlock(b)

  if p.Length() != 0 {
    t.Errorf("Unexpected block slice length after remove %d", p.Length())
  }
*/
  b1 := Block{Title: "Mock 1"}
  b2 := Block{Title: "Mock 2"}
  b3 := Block{Title: "Mock 3"}

  p.AddBlock(b1)
  p.AddBlock(b2)
  p.AddBlock(b3)

  //p.RemoveBlock(b2)

  //if p.Length() != 2 {
    //t.Errorf("Unexpected block slice length after remove %d", p.Length())
  //}

  //if p.Blocks[0] != b1 {
    //t.Errorf("Unexpected block %s at index 0", p.Blocks[0])
  //}

  //if p.Blocks[1] != b3 {
    //t.Errorf("Unexpected block %s at index 1", p.Blocks[1])
  //}
}

func TestJson(t *testing.T) {
  p := Page{}
  p.AddBlock(Block{Title: "Mock Title", Content: "Mock Content"})
  _, err := json.Marshal(p)
  if err != nil {
    t.Errorf("%s", err)
  }

  /*
  expected := `{"doctype":"\u003c!doctype html\u003e","data":null,"blocks":[{"title":"Mock Title","content":"Mock Content"}]}`

  if string(b) != expected {
    os.Stdout.Write(append(b, 0x0A))
    t.Errorf("unexpected JSON output")  
  }
  */
}
