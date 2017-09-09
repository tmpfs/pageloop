package pageloop

import (
  //"fmt"
  //"os/exec"
	//"regexp"
	//"strings"
	//"net/http"
  //"mime"
  //"path/filepath"
	//"encoding/json"
  "github.com/tmpfs/pageloop/model"
)

var(
  adapter *CommandAdapter
)

// Abstraction that allows many different interfaces to
// the data model whether it is a string command interpreter,
// REST API endpoints, JSON RPC or any other bridge to the
// outside world.
//
// For simplicity with access over HTTP this implementation always
// returns errors with an associated HTTP status code.
type CommandAdapter struct {
  Root *PageLoop
}

type Command struct {
  Root *PageLoop
}

// List all system templates and user applications
// that have been marked as a template.
func (b *CommandAdapter) ListApplicationTemplates() []*model.Application {
  // Get built in and user templates
  c := b.Root.Host.GetByName("template")
  u := b.Root.Host.GetByName("user")
  list := append(c.Apps, u.Apps...)
  var apps []*model.Application
  for _, app := range list {
    if app.IsTemplate {
      apps = append(apps, app)
    }
  }
  return apps
}
