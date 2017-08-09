package model

import (
  "testing"
)

func TestLoadApplication(t *testing.T) {
  l := FileSystemLoader{}
  var err error
  err = l.LoadApplication("../test/fixtures/mock-app", &Application{})
  if err != nil {
    t.Error(err)
  }
}
