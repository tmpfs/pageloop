package model

import (
  "testing"
)

func TestLoadApplication(t *testing.T) {
  l := FileSystemLoader{}
  var app *Application
  var err error
  app, err = l.LoadApplication("../test/fixtures/mock-app", &Application{})
  if err != nil {
    t.Error(err)
  }
  if app == nil {
    t.Error("Unexpected nil application from LoadApplication")
  }
}
