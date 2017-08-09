package model

import (
  //"os"
  "fmt"
  "log"
  "testing"
  //"encoding/json"
)

func TestApplication(t *testing.T) {
  app := Application{Title: "Mock Application"}

  if app.Title == "" {
    t.Errorf("Unexpected empty application title")
  }

  // app.Base = "test/fixtures/mock-app"

  app.Load("../test/fixtures/mock-app", nil)

  expected := 5
  if len(app.Files) != expected {
    t.Errorf("Unexpected number of files %d", len(app.Files))  
  }

  expected = 1
  if len(app.Pages) != expected {
    t.Errorf("Unexpected number of pages %d", len(app.Pages))  
  }

  if len(app.Pages[0].UserData) == 0 {
    t.Errorf("Unexpected empty user data for %s", app.Pages[0].file.Path)  
  }

  log.Println(app.Pages[0].file.Path)
  log.Println(app.Pages[0].UserData)

  fmt.Println("--- render result ---")
  fmt.Println()
  //app.Render(&app.Pages[0])
}
