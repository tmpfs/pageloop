package adapter

import (
  "fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

// Command functions

// ROOT

// Meta information.
func (b *CommandAdapter) Meta() map[string]interface{} {
  status := make(map[string]interface{})
  status["name"] = b.Name
  status["version"] = b.Version
  return status
}

// CONTAINERS

// Read all containers in a host.
func (b *CommandAdapter) ReadHost() []*Container {
  return b.Host.Containers
}

// Read a container.
func (b *CommandAdapter) ReadContainer(name string) (*Container, *StatusError) {
  c := b.Host.GetByName(name)
  if c == nil {
    return nil, CommandError(http.StatusNotFound, "")
  }
  return c, nil
}

// APPLICATIONS

func (b *CommandAdapter) ReadApplication(c string, name string) (*Application, *StatusError) {
  if container, err := b.ReadContainer(c); err != nil {
    return nil, err
  } else {
    app :=  container.GetByName(name)
    if app == nil {
      return nil, CommandError(http.StatusNotFound, "")
    }
    return app, nil
  }
}

func (b *CommandAdapter) ReadApplicationFiles(c string, name string) ([]*File, *StatusError) {
  if a, err :=  b.ReadApplication(c, name); err != nil {
    return nil, err
  } else {
    return a.Files, nil
  }
}

func (b *CommandAdapter) ReadApplicationPages(c string, name string) ([]*Page, *StatusError) {
  if a, err :=  b.ReadApplication(c, name); err != nil {
    return nil, err
  } else {
    return a.Pages, nil
  }
}

// JOBS

// List jobs.
func (b *CommandAdapter) ListJobs() []*Job {
  return Jobs.Active
}

// Read a job.
func (b *CommandAdapter) ReadJob(id string) (*Job, *StatusError) {
  var job *Job = Jobs.ActiveJob(id)
  if job == nil {
    return nil, CommandError(http.StatusNotFound, "")
  }
  return job, nil
}

// Abort an active job.
func(b *CommandAdapter) AbortJob(id string) (*Job, *StatusError) {
  var err error
  var job *Job = Jobs.ActiveJob(id)
  if job == nil {
    return nil, CommandError(http.StatusNotFound, "")
  }

  if err = Jobs.Abort(job); err != nil {
    return nil, CommandError(http.StatusConflict, "")
  }

  // Accepted for processing
  fmt.Printf("[job:%d] aborted %s\n", job.Number, job.Id)

  return job, nil
}

// MISC

// List all system templates and user applications
// that have been marked as a template.
func (b *CommandAdapter) ListApplicationTemplates() []*Application {
  // Get built in and user templates
  c := b.Host.GetByName("template")
  u := b.Host.GetByName("user")
  list := append(c.Apps, u.Apps...)
  var apps []*Application
  for _, app := range list {
    if app.IsTemplate {
      apps = append(apps, app)
    }
  }
  return apps
}

