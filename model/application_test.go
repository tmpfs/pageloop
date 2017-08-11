package model

import (
  //"os"
  "fmt"
  "log"
  "testing"
  //"encoding/json"
)

func TestApplication(t *testing.T) {
  var err error
  app := Application{}

  err = app.Load("../test/fixtures/mock-app", nil)
  if err != nil {
    t.Error(err)
  }

  expected := 4
  if len(app.Files) != expected {
    t.Errorf("Unexpected number of files %d", len(app.Files))  
  }

  expected = 1
  if len(app.Pages) != expected {
    t.Errorf("Unexpected number of pages %d", len(app.Pages))  
  }

  if len(app.Pages[0].PageData) == 0 {
    t.Errorf("Unexpected empty user data for %s", app.Pages[0].file.Path)  
  }

  log.Println(app.Pages[0].file.Path)
  log.Println(app.Pages[0].PageData)

  fmt.Println("--- render result ---")
  fmt.Println()
  //app.Render(&app.Pages[0])
}
