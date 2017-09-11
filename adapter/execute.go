
package adapter

import (
  "fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

// Command functions.
//
// This is where the actual implementation for the various commands
// are using native types as arguments.
type CommandExecute struct {
  *CommandAdapter
}

// ROOT

// Meta information (/).
func (b *CommandExecute) ReadMeta() map[string]interface{} {
  status := make(map[string]interface{})
  status["name"] = b.Name
  status["version"] = b.Version
  return status
}

// CONTAINERS

// Read all containers in a host.
func (b *CommandExecute) ReadHost() []*Container {
  return b.Host.Containers
}

// Read a container.
func (b *CommandExecute) ReadContainer(container *Container) (*Container, *StatusError) {
  return container, nil
}

// APPLICATIONS

// Run an application build task.
func(b *CommandExecute) RunTask(a *Application, task string) (*Job, *StatusError) {
  var err error
  var job *Job
  // No build configuration of missing build task
  if !a.HasBuilder() || a.Builder.Tasks[task] == "" {
    return nil, CommandError(http.StatusNotFound, "")
  }

  // Run the task and get a job
  if job, err = a.Builder.Run(task, &TaskJobComplete{}); err != nil {
    // Send conflict if job already running, this is a bit flaky is Run()
    // starts returning errors for other reasons :(
    return nil, CommandError(http.StatusConflict, err.Error())
  }

  // Accepted for processing
  fmt.Printf("[job:%d] started %s\n", job.Number, job.Id)

  return job, nil
}

// Create application.
func (b *CommandExecute) CreateApp(container *Container, a *Application) (*Application, *StatusError) {
  // TODO: do not allow creating apps on non-user containers!
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

// Delete an application.
func (b *CommandExecute) DeleteApp(container *Container, app *Application) (*Application, *StatusError) {
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

// Read an application.
func (b *CommandExecute) ReadApplication(app *Application) *Application {
  return app
}

// Read the files for an application.
func (b *CommandExecute) ReadApplicationFiles(app *Application) []*File {
  return app.Files
}

// Read the pages for an application.
func (b *CommandExecute) ReadApplicationPages(app *Application) []*Page {
  return app.Pages
}

// FILES / PAGES

// Create a new file and publish it, the file cannot already exist on disc.
func (b *CommandExecute) CreateFile(a *Application, url string, content []byte) (*File, *StatusError) {
  var err error
	var file *File = a.Urls[url]

	if file != nil {
    return nil, CommandError(http.StatusConflict,"File already exists %s", url)
	}
  if a.ExistsConflict(url) {
    return nil, CommandError(http.StatusConflict,"File already exists, publish conflict on %s", url)
  }

  if file, err = a.Create(url, content); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }

  return file, nil
}

// Create a file from a template.
func (b *CommandExecute) CreateFileTemplate(a *Application, url string, template *ApplicationTemplate) (*File, *StatusError) {
  var err error
  var file *File
  var content []byte

  if file, err = b.Host.LookupTemplateFile(template); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }

  if file == nil {
    return nil, CommandError(http.StatusNotFound, "Template file %s does not exist", template.File)
  }

  content = file.Source(true)
  return b.CreateFile(a, url, content)
}

// Update file content.
func (b *CommandExecute) UpdateFile(a *Application, f *File, content []byte) (*File, *StatusError) {
  if err := a.Update(f, content); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }
  return f, nil
}

// Delete a file.
func (b *CommandExecute) DeleteFile(a *Application, url string) (*File, *StatusError) {
  var err error
  var file *File = a.Urls[url]
  if file == nil {
    return nil, CommandError(http.StatusNotFound, "")
  }
  if err = a.Del(file); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }
  return file, nil
}

// Batch delete files.
func (b *CommandExecute) DeleteFiles(app *Application, list UrlList) ([]*File, *StatusError) {
  var files []*File
  for _, url := range list {
    if file, err := b.DeleteFile(app, url); err != nil {
      return nil, err
    } else {
      files = append(files, file)
    }
  }
  return files, nil
}

// Move a file.
func (b *CommandExecute) MoveFile(app *Application, file *File, dest string) (*File, *StatusError) {
  if err := app.Move(file, dest); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }
  return file, nil
}

// Read a file.
func (b *CommandExecute) ReadFile(app *Application, url string) (*File, *StatusError) {
  var file *File
  file = app.Urls[url]
  // Cannot find the target file
  if file == nil {
    return nil, CommandError(http.StatusNotFound, "File %s not found", url)
  }
  return file, nil
}

// Read a page.
func (b *CommandExecute) ReadPage(app *Application, url string) (*Page, *StatusError) {
  file := app.Urls[url]
  // Cannot find the target file
  if file == nil {
    return nil, CommandError(http.StatusNotFound, "File %s not found", url)
  }
  // File is not a page type
  if file.Page() == nil {
    return nil, CommandError(http.StatusNotFound, "Page %s not found", url)
  }
  return file.Page(), nil
}

// JOBS

// List jobs.
func (b *CommandExecute) ReadActiveJobs() []*Job {
  return Jobs.Active
}

// Read a job.
func (b *CommandExecute) ReadJob(id string) (*Job, *StatusError) {
  var job *Job = Jobs.ActiveJob(id)
  if job == nil {
    return nil, CommandError(http.StatusNotFound, "")
  }
  return job, nil
}

// Abort an active job.
func(b *CommandExecute) DeleteJob(id string) (*Job, *StatusError) {
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
func (b *CommandExecute) ReadApplicationTemplates() []*Application {
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
