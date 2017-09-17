package adapter

import (
  "strings"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

// List of URLs used for bulk file operations.
type UrlList []string

// Command functions

// ROOT

// Meta information (/).
func (b *CommandAdapter) ReadMeta() map[string]interface{} {
  return b.CommandExecute.ReadMeta()
}

// Stats information (/stats).
func (b *CommandAdapter) ReadStats() map[string]interface{} {
  return b.CommandExecute.ReadStats()
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
func(b *CommandAdapter) RunTask(c string, a string, task string) (*Job, *StatusError) {
  task = strings.TrimPrefix(task, SLASH)
  task = strings.TrimSuffix(task, SLASH)
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
func (b *CommandAdapter) ReadFile(c string, a string, f string) (*Application, *File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, nil, err
  } else {
    if file, err := b.CommandExecute.ReadFile(app, f); err != nil {
      return nil, nil, err
    } else {
      return app, file, nil
    }
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

// Move a file
// TODO: restore and test this
func (b *CommandAdapter) MoveFile(c, a, f, dest string) (*File, *StatusError) {
  if app, file, err :=  b.ReadFile(c, a, f); err != nil {
    return nil, err
  } else {
    if file, err := b.CommandExecute.MoveFile(app, file, dest); err != nil {
      return nil, err
    } else {
      return file, nil
    }
  }
}

// Delete a file.
func (b *CommandAdapter) DeleteFile(c, a, f string) (*File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.DeleteFile(app, f)
  }
}

// Delete a list of files.
func (b *CommandAdapter) DeleteFiles(c string, a string, l UrlList) ([]*File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.DeleteFiles(app, l)
  }
}

// Create file content.
func (b *CommandAdapter) CreateFile(c string, a string, f string, content []byte) (*File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.CreateFile(app, f, content)
  }
}

// Create file from a template.
func (b *CommandAdapter) CreateFileTemplate(c string, a string, f string, tpl *ApplicationTemplate) (*File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return b.CommandExecute.CreateFileTemplate(app, f, tpl)
  }
}

// Update file content.
//
// This can return a file or page so that when updating pages the page data
// can be updated on the client.
//
// TODO: create an interface for file/page types and return that rather than interface{}
func (b *CommandAdapter) UpdateFile(c string, a string, f string, content []byte) (interface{}, *StatusError) {
  if app, file, err :=  b.ReadFile(c, a, f); err != nil {
    return nil, err
  } else {
    if file, err := b.CommandExecute.UpdateFile(app, file, content); err != nil {
      return nil, err
    } else {
      if file.Page() != nil {
        return file.Page(), nil
      }
      return file, nil
    }
  }
}

// JOBS

// List jobs.
func (b *CommandAdapter) ReadActiveJobs() []*Job {
  return b.CommandExecute.ReadActiveJobs()
}

// Read a job.
func (b *CommandAdapter) ReadJob(id string) (*Job, *StatusError) {
  return b.CommandExecute.ReadJob(id)
}

// Abort an active job.
func(b *CommandAdapter) DeleteJob(id string) (*Job, *StatusError) {
  return b.CommandExecute.DeleteJob(id)
}

// MISC

// List all system templates and user applications
// that have been marked as a template.
func (b *CommandAdapter) ReadApplicationTemplates() []*Application {
  return b.CommandExecute.ReadApplicationTemplates()
}
