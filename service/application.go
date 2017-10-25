package service

import(
  "fmt"
  "net/http"
  "strings"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

// Handler for asynchronous background tasks.
type TaskJobComplete struct {}

func (tj *TaskJobComplete) Done(err error, job *Job) {
  // TODO: send reply to the client over websocket
  fmt.Printf("[job:%d] completed %s\n", job.Number, job.Id)
  Jobs.Stop(job)
}

type ApplicationRequest struct {
  // Application name (id)
  Name string `json:"name"`

  // Container name
  Container string `json:"container"`

  // Application display name
  DisplayName string `json:"display"`

  // Application description
  Description string `json:"description"`

  // Mark this application as a template
  IsTemplate bool `json:"is-template,omitempty"`

  // Name of a task to find
  Task string `json:"task,omitempty"`

  // List used for batch operations
  Batch *UrlList `json:"batch,omitempty"`

	// A source template for this application
	Template *ApplicationTemplate `json:"template,omitempty"`
}

func (req *ApplicationRequest) ToApplication(container *Container) *Application {
  return &Application{
    Container: container,
    Name: req.Name,
    DisplayName: req.DisplayName,
    Description: req.Description,
    IsTemplate: req.IsTemplate}
}

type ApplicationReferenceRequest struct {
  // A reference to an application in the form: file://pageloop.com/{container}/{application}
  Ref string `json:"ref,omitempty"`
}

type ApplicationBatchRequest struct {
  // A reference to an application in the form: file://pageloop.com/{container}/{application}
  Ref string `json:"ref,omitempty"`

  // List used for batch operations
  Batch *UrlList `json:"batch,omitempty"`
}

type ApplicationTaskRequest struct {
  // A reference to an application in the form: file://pageloop.com/{container}/{application}
  Ref string `json:"ref,omitempty"`

  // Name of a task to find
  Task string `json:"task,omitempty"`
}

type AppService struct {
  Host *Host

  // Reference to the mountpoint manager
  Mountpoints *MountpointManager
}

// Read an application.
func (s *AppService) Read(req *ApplicationReferenceRequest, reply *ServiceReply) *StatusError {
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, err := ref.FindApplication(s.Host); err != nil {
    return err
  } else {
    reply.Reply = app
  }
  return nil
}

// Read the files for an application.
func (s *AppService) ReadFiles(req *ApplicationReferenceRequest, reply *ServiceReply) *StatusError {
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, err := ref.FindApplication(s.Host); err != nil {
    return err
  } else {
    reply.Reply = app.Files
  }
  return nil
}

// Read the pages for an application.
func (s *AppService) ReadPages(req *ApplicationReferenceRequest, reply *ServiceReply) *StatusError {
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, err := ref.FindApplication(s.Host); err != nil {
    return err
  } else {
    reply.Reply = app.Pages
  }
  return nil
}

// Delete an application.
func (s *AppService) Delete(req *ApplicationReferenceRequest, reply *ServiceReply) *StatusError {
  //println("application delete called")
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if container, app, err := ref.FindApplication(s.Host); err != nil {
    //println("returning err")
    return err
  } else {
    if app.Protected {
      return CommandError(http.StatusForbidden, "Cannot delete protected application")
    }

    // Stop serving files for the application
    s.Mountpoints.UnmountApplication(app)

    // Delete the mountpoint
    if err := s.Mountpoints.DeleteApplicationMountpoint(app.Url); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }

    // Delete the files
    if err := app.DeleteApplicationFiles(); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }

    // Delete the in-memory application
    container.Del(app)
  }
  return nil
}

// Batch delete files.
func (s *AppService) DeleteFiles(req *ApplicationBatchRequest, reply *ServiceReply) *StatusError {
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, err := ref.FindApplication(s.Host); err != nil {
    return err
  } else {
    var file *File
    var files []*File
    for _, url := range *req.Batch {
      file  = app.Urls[url]
      if file == nil {
        return CommandError(http.StatusNotFound, "File not found for url %s", url)
      }

      if err := app.Del(file); err != nil {
        return CommandError(http.StatusInternalServerError, err.Error())
      }

      files = append(files, file)
    }
    reply.Reply = files
  }
  return nil
}

// Run an application build task.
func(s *AppService) RunTask(req *ApplicationTaskRequest, reply *ServiceReply) *StatusError {
  task := req.Task
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, err := ref.FindApplication(s.Host); err != nil {
    return err
  } else {
    task = strings.TrimPrefix(task, SLASH)
    // No build configuration of missing build task
    if !app.HasBuilder() {
      return CommandError(
        http.StatusNotFound, "Application %s does not have a build configuration (needs build.yml)", app.Name)
    }

    if app.Builder.Tasks[task] == "" {
      return CommandError(
        http.StatusNotFound, "Build configuration task %s not found", task)
    }

    // Run the task and get a job
    if job, err := app.Builder.Run(task, &TaskJobComplete{}); err != nil {
      // Send conflict if job already running, this is a bit flaky is Run()
      // starts returning errors for other reasons :(
      return CommandError(http.StatusConflict, err.Error())
    } else {
      // Accepted for processing
      fmt.Printf("[job:%d] started %s\n", job.Number, job.Id)

      reply.Reply = job
      reply.Status = http.StatusAccepted
    }
  }
  return nil
}

func LookupApplication(host *Host, req *ApplicationRequest) (*Container, *Application, *StatusError) {
  c := host.GetByName(req.Container)
  if c == nil {
    return nil, nil, CommandError(http.StatusNotFound, "Container %s not found", req.Container)
  }
  a := c.GetByName(req.Name)
  if a == nil {
    return nil, nil, CommandError(http.StatusNotFound, "Application %s not found", req.Name)
  }
  return c, a, nil
}
