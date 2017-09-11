package adapter

import (
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

// Command functions

// ROOT

// Meta information.
func (b *CommandAdapter) Meta() map[string]interface{} {
  return b.CommandExecute.Meta()
}

// CONTAINERS

// Read all containers in a host.
func (b *CommandAdapter) ReadHost() []*Container {
  return b.CommandExecute.ReadHost()
}

// Read a container.
func (b *CommandAdapter) ReadContainer(c string) (*Container, *StatusError) {
  container := b.Host.GetByName(c)
  if container == nil {
    return nil, CommandError(http.StatusNotFound, "Container %s not found", c)
  }
  return b.CommandExecute.ReadContainer(container)
}

// APPLICATIONS

// Create application.
func (b *CommandAdapter) CreateApp(c string, a *Application) (*Application, *StatusError) {
  // TODO: do not allow creating apps on non-user containers!
  if container, err := b.ReadContainer(c); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.CreateApp(container, a)
  }
}

// Delete an application.
func (b *CommandAdapter) DeleteApp(c string, a string) (*Application, *StatusError) {
  if container, app, err := b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.DeleteApp(container, app)
  }
}

// Run an application build task.
func(b *CommandAdapter) RunAppTask(c string, a string, task string) (*Job, *StatusError) {
  if _, app, err := b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.RunTask(app, task)
  }
}

// Read an application.
func (b *CommandAdapter) ReadApplication(c string, a string) (*Container, *Application, *StatusError) {
  if container, err := b.ReadContainer(c); err != nil {
    return nil, nil, err
  } else {
    app :=  container.GetByName(a)
    if app == nil {
      return nil, nil, CommandError(http.StatusNotFound, "Application %s not found", a)
    }
    return container, b.CommandExecute.ReadApplication(app), nil
  }
}

// Read the files for an application.
func (b *CommandAdapter) ReadApplicationFiles(c string, a string) ([]*File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.ReadApplicationFiles(app), nil
  }
}

// Read the pages for an application.
func (b *CommandAdapter) ReadApplicationPages(c string, a string) ([]*Page, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.ReadApplicationPages(app), nil
  }
}

// FILES / PAGES

// Read a file.
func (b *CommandAdapter) ReadFile(c string, a string, f string) (*File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.ReadFile(app, f)
  }
}

// Read a page.
func (b *CommandAdapter) ReadPage(c string, a string, f string) (*Page, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.ReadPage(app, f)
  }
}

// JOBS

// List jobs.
func (b *CommandAdapter) ListJobs() []*Job {
  return b.CommandExecute.ListJobs()
}

// Read a job.
func (b *CommandAdapter) ReadJob(id string) (*Job, *StatusError) {
  return b.CommandExecute.ReadJob(id)
}

// Abort an active job.
func(b *CommandAdapter) AbortJob(id string) (*Job, *StatusError) {
  return b.CommandExecute.AbortJob(id)
}

// MISC

// List all system templates and user applications
// that have been marked as a template.
func (b *CommandAdapter) ListApplicationTemplates() []*Application {
  return b.CommandExecute.ListApplicationTemplates()
}

