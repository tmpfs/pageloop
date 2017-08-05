package blocks

import (
  //"os"
  "testing"
  //"encoding/json"
)

func TestApplication(t *testing.T) {
  app := Application{Title: "Mock Application"}

  if app.Title == "" {
    t.Errorf("Unexpected empty application title")
  }
}

func TestApplicationJson(t *testing.T) {
  /*
  p := Page{DocType: HTML5}
  p.AddBlock(Block{Title: "Mock Title", Content: "Mock Content"})
  b, err := json.Marshal(p)
  if err != nil {
    t.Errorf("%s", err)
  }

  os.Stdout.Write(append(b, 0x0A))

  expected := `{"doctype":"\u003c!doctype html\u003e","data":null,"blocks":[{"title":"Mock Title","content":"Mock Content"}]}`

  if string(b) != expected {
    t.Errorf("unexpected JSON output")  
  }
  */
}
