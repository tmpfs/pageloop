package adapter

import (
  "fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/util"
)

// Handler for asynchronous background tasks.
type TaskJobComplete struct {}

func (tj *TaskJobComplete) Done(err error, job *Job) {
  // TODO: send reply to the client over websocket
  fmt.Printf("[job:%d] completed %s\n", job.Number, job.Id)
  Jobs.Stop(job)
}

// TODO: implement action generation and execution
const(
  // Basic CRUD operations
  OperationCreate = iota
  OperationRead
  OperationUpdate
  OperationDelete
)

// A command action is a simple representation of a command invocation
// it can be used to execute a command without any object references.
type CommandAction struct {
  Operation int
  ContainerName string
  ApplicationName string
  FileUrl string
}

// Abstraction that allows many different interfaces to
// the data model whether it is a string command interpreter,
// REST API endpoints, JSON RPC or any other bridge to the
// outside world.
//
// For simplicity with access over HTTP this implementation always
// returns errors with an associated HTTP status code.
type CommandAdapter struct {
  Host *Host
  Mountpoints *MountpointManager
}

// List jobs.
func (b *CommandAdapter) ListJobs() []*Job {
  return Jobs.Active
}

// List containers.
func (b *CommandAdapter) ListContainers() []*Container {
  return b.Host.Containers
}

// List applications in a container.
func (b *CommandAdapter) ListApplications(c *Container) []*Application {
  return c.Apps
}

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

// Create application.
func (b *CommandAdapter) CreateApplication(c *Container, a *Application) (*Application, *StatusError) {

  var app *Application

  existing := c.GetByName(a.Name)
  if existing != nil {
    return nil, CommandError(http.StatusPreconditionFailed, "Application %s already exists", a.Name)
  }

  // Get mountpoint URL.
  a.Url = a.MountpointUrl(c)

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
    if app, err = b.Mountpoints.LoadMountpoint(*mountpoint, c); err != nil {
      return nil, CommandError(http.StatusInternalServerError, err.Error())
    } else {
      // Return the new application reference
      return app, nil
    }
  }
  return app, nil
}

// Delete an application.
func (b *CommandAdapter) DeleteApplication(c *Container, a *Application) (*Application, *StatusError) {
  if a.Protected {
    return nil, CommandError(http.StatusForbidden, "Cannot delete protected application")
  }

  // Stop serving files for the application
  b.Mountpoints.UnmountApplication(a)

  // Delete the mountpoint
  if err := b.Mountpoints.DeleteApplicationMountpoint(a.Url); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }

  // Delete the files
  if err := a.DeleteApplicationFiles(); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }

  // Delete the in-memory application
  c.Del(a)

  return a, nil
}

func(b *CommandAdapter) RunTask(a *Application, task string) (*Job, *StatusError) {
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

// Move a file.
func (b *CommandAdapter) MoveFile(a *Application, f *File, dest string) *StatusError {
  if err := a.Move(f, dest); err != nil {
    return CommandError(http.StatusInternalServerError, err.Error())
  }
  return nil
}

// Create a file from a template.
func (b *CommandAdapter) CreateFileTemplate(a *Application, url string, template *ApplicationTemplate) (*File, *StatusError) {
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

// Create a new file and publish it, the file cannot already exist on disc.
func (b *CommandAdapter) CreateFile(a *Application, url string, content []byte) (*File, *StatusError) {
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

// Update file content.
func (b *CommandAdapter) UpdateFile(a *Application, f *File, content []byte) (*File, *StatusError) {
  if err := a.Update(f, content); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }
  return f, nil
}

// Delete a file.
func (b *CommandAdapter) DeleteFile(a *Application, url string) (*File, *StatusError) {
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
