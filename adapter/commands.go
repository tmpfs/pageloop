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

// Create application.
func (b *CommandAdapter) CreateApp(c string, a *Application) (*Application, *StatusError) {
  // TODO: do not allow creating apps on non-user containers!
  if container, err := b.ReadContainer(c); err != nil {
    return nil, err
  } else {

    var app *Application

    existing := container.GetByName(a.Name)
    if existing != nil {
      return nil, CommandError(http.StatusPreconditionFailed, "Application %s already exists", a.Name)
    }

    // Get mountpoint URL.
    a.Url = a.MountpointUrl(container)

    // Mountpoint exists.
    exists := b.Mountpoints.HasMountpoint(a.Url)
    if exists {
      return nil, CommandError(http.StatusPreconditionFailed, "Mountpoint URL %s already exists", a.Url)
    }

    // Create and save a mountpoint for the application.
    if mountpoint, err := b.Mountpoints.CreateMountpoint(a); err != nil {
      return nil, CommandError(http.StatusInternalServerError, err.Error())
    } else {
      // Handle creating from a template.
      if a.Template != nil {
        // Find the template application.
        if source, err := b.Host.LookupTemplate(a.Template); err != nil {
          return nil, CommandError(http.StatusBadRequest, err.Error())
        } else {
          // Copy template source files.
          if err := a.CopyApplicationTemplate(source); err != nil {
            return nil, CommandError(http.StatusInternalServerError, err.Error())
          }
        }
      }

      // Load and publish the app source files, note that we get a new application back
      // after loading the mountpoint.
      if app, err = b.Mountpoints.LoadMountpoint(*mountpoint, container); err != nil {
        return nil, CommandError(http.StatusInternalServerError, err.Error())
      } else {
        // Return the new application reference
        return app, nil
      }
    }
    return app, nil
  }
}

// Delete an application.
func (b *CommandAdapter) DeleteApp(c string, a string) (*Application, *StatusError) {
  if container, app, err := b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    if app.Protected {
      return nil, CommandError(http.StatusForbidden, "Cannot delete protected application")
    }

    // Stop serving files for the application
    b.Mountpoints.UnmountApplication(app)

    // Delete the mountpoint
    if err := b.Mountpoints.DeleteApplicationMountpoint(app.Url); err != nil {
      return nil, CommandError(http.StatusInternalServerError, err.Error())
    }

    // Delete the files
    if err := app.DeleteApplicationFiles(); err != nil {
      return nil, CommandError(http.StatusInternalServerError, err.Error())
    }

    // Delete the in-memory application
    container.Del(app)

    return app, nil
  }
}

// APPLICATIONS

// Read an application.
func (b *CommandAdapter) ReadApplication(c string, a string) (*Container, *Application, *StatusError) {
  if container, err := b.ReadContainer(c); err != nil {
    return nil, nil, err
  } else {
    app :=  container.GetByName(a)
    if app == nil {
      return nil, nil, CommandError(http.StatusNotFound, "Application %s not found", a)
    }
    return container, app, nil
  }
}

// Read the files for an application.
func (b *CommandAdapter) ReadApplicationFiles(c string, a string) ([]*File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return app.Files, nil
  }
}

// Read the pages for an application.
func (b *CommandAdapter) ReadApplicationPages(c string, a string) ([]*Page, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    return app.Pages, nil
  }
}

// FILES / PAGES

// Read a file.
func (b *CommandAdapter) ReadFile(c string, a string, f string) (*File, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    file := app.Urls[f]
    if file == nil {
      return nil, CommandError(http.StatusNotFound, "File %s not found", f)
    }
    return file, nil
  }
}

// Read a page.
func (b *CommandAdapter) ReadPage(c string, a string, f string) (*Page, *StatusError) {
  if _, app, err :=  b.ReadApplication(c, a); err != nil {
    return nil, err
  } else {
    file := app.Urls[f]
    // Cannot find the target file
    if file == nil {
      return nil, CommandError(http.StatusNotFound, "File %s not found", f)
    }
    // File is not a page type
    if file.Page() == nil {
      return nil, CommandError(http.StatusNotFound, "Page %s not found", f)
    }
    return file.Page(), nil
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

