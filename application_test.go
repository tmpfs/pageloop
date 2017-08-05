package blocks

import (
  "os"
  "testing"
  "encoding/json"
)

func TestApplication(t *testing.T) {
  app := Application{Title: "Mock Application"}

  if app.Title == "" {
    t.Errorf("Unexpected empty application title")
  }

  app.Load("test/fixtures/mock-app")

  expected := 2
  if len(app.Files) != expected {
    t.Errorf("Unexpected number of files %d", len(app.Files))  
  }

  expected = 1
  if len(app.Pages) != expected {
    t.Errorf("Unexpected number of pages %d", len(app.Pages))  
  }
}

func TestApplicationJson(t *testing.T) {
  app := Application{Title: "Mock Application"}
  app.Load("test/fixtures/mock-app")
  b, err := json.Marshal(app)
  if err != nil {
    t.Errorf("%s", err)
  }

  os.Stdout.Write(append(b, 0x0A))

  //expected := `{"doctype":"\u003c!doctype html\u003e","data":null,"blocks":[{"title":"Mock Title","content":"Mock Content"}]}`

  //if string(b) != expected {
    //t.Errorf("unexpected JSON output")  
  //}
}
