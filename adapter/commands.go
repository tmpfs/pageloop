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
func (b *CommandAdapter) ReadContainer(c string) (*Container, *StatusError) {
  container := b.Host.GetByName(c)
  if container == nil {
    return nil, CommandError(http.StatusNotFound, "Container %s not found", c)
  }
  return container, nil
}

// APPLICATIONS

func (b *CommandAdapter) ReadApplication(c string, a string) (*Application, *StatusError) {
  if container, err := b.ReadContainer(c); err != nil {
    return nil, err
  } else {
    app :=  container.GetByName(a)
    if app == nil {
      return nil, CommandError(http.StatusNotFound, "Application %s not found", a)
    }
    return app, nil
  }
}

func (b *CommandAdapter) ReadApplicationFiles(c string, a string) ([]*File, *StatusError) {
  if app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return app.Files, nil
  }
}

func (b *CommandAdapter) ReadApplicationPages(c string, a string) ([]*Page, *StatusError) {
  if app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return app.Pages, nil
  }
}

// FILES / PAGES

func (b *CommandAdapter) ReadFile(c string, a string, f string) (*File, *StatusError) {
  if app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    file := app.Urls[f]
    if file == nil {
      return nil, CommandError(http.StatusNotFound, "File %s not found", f)
    }
    return file, nil
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

