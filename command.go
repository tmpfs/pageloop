package pageloop

import (
  "fmt"
  //"os/exec"
	//"regexp"
	//"strings"
  "net/http"
  //"mime"
  //"path/filepath"
	//"encoding/json"
  "github.com/tmpfs/pageloop/model"
)

var(
  adapter *CommandAdapter
)

type StatusError struct {
	Status int
	Message string
}

func (s StatusError) Error() string {
	return s.Message
}

func CommandError(status int, message string, a ...interface{}) *StatusError {
  if message == "" {
    message = http.StatusText(status)
  }
	return &StatusError{Status: status, Message: fmt.Sprintf(message, a...)}
}

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

// Create application.
func (b *CommandAdapter) CreateApplication(c *model.Container, a *model.Application) *StatusError {
  existing := c.GetByName(a.Name)
  if existing != nil {
    return CommandError(http.StatusPreconditionFailed, "Application %s already exists", a.Name)
  }

  // Get mountpoint URL.
  a.Url = a.MountpointUrl(c)

  // Mountpoint exists.
  exists := b.Root.HasMountpoint(a.Url)
  if exists {
    return CommandError(http.StatusPreconditionFailed, "Mountpoint URL %s already exists", a.Url)
  }

  // Create and save a mountpoint for the application.
  if mountpoint, err := b.Root.CreateMountpoint(a); err != nil {
    return CommandError(http.StatusInternalServerError, err.Error())
  } else {
    // Handle creating from a template.
    if a.Template != nil {
      // Find the template application.
      if source, err := b.Root.LookupTemplate(a.Template); err != nil {
        return CommandError(http.StatusBadRequest, err.Error())
      } else {
        // Copy template source files.
        if err := b.Root.CopyApplicationTemplate(a, source); err != nil {
          return CommandError(http.StatusInternalServerError, err.Error())
        }
      }
    }

    // Load and publish the app source files, note that we get a new application back
    // after loading the mountpoint.
    if app, err := b.Root.LoadMountpoint(*mountpoint, c); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    } else {
      // Mount the application
      b.Root.MountApplication(app)
    }
  }
  return nil
}

// List containers.
func (b *CommandAdapter) ListContainers() []*model.Container {
  return b.Root.Host.Containers
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
