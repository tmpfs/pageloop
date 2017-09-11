// Package adapter provides a command adapter for interfacing
// network requests with the underlying model.
package adapter

import (
  "fmt"
  "net/url"
  "net/http"
  "reflect"
  "strings"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/util"
)

var(
  ActionMap map[*Action]*ActionDefinition = make(map[*Action]*ActionDefinition)
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
//
// Path references take the form:
//
// /{type}?/{context}?/{target}?/{action}?/{item}?
//
// Where item is a trailer that may includes slashes to represent a file URL.
//
// The context part corresponds to a container and the target part corresponds
// to an application.
//
// If a definition maps a part using the wildcard (*) it will match any string.
type Action struct {
  // Source HTTP verb that is translated to an operation constant
  Verb string
  // A request URL
  Url *url.URL
  // The path for the request
  Path string
  // Parsed path parts split on a slash
  Parts []string
  // The CRUD operation to perform
  Operation int

  // The operation type
  Type string
  // Context for the operation. May be a container reference, job number etc.
  Context string
  // Target for the operation, typically an application.
  Target string
  // An action or filter operation for the request.
  Action string
  // An item, may contain slashes.
  Item string
}

func (act *Action) IsRoot() bool {
  return act.Path == ""
}

func (act *Action) TypeOnly() bool {
  return act.Type != "" && act.Context == "" && act.Target == "" && act.Action == "" && act.Item == ""
}

func (act *Action) ContextOnly() bool {
  return act.Type != "" && act.Context != "" && act.Target == "" && act.Action == "" && act.Item == ""
}

func (act *Action) MatchType(in *Action) bool {
  return act.Type == in.Type
}

func (act *Action) Wildcard(val string) bool {
  return val == "*"
}

func (act *Action) Match(in *Action) bool {
  if act.Operation != in.Operation {
    return false
  }

  // Exact path match
  if act.Path == in.Path {
    return true
  }

  if act.TypeOnly() && in.TypeOnly() && act.Type == in.Type {
    return true
  }

  // Got a type match
  if act.MatchType(in) {
    // Deal with context only
    if act.ContextOnly() && in.ContextOnly() && (act.Wildcard(act.Context) || act.Context == in.Context) {
      return true
    }
  }

  return false
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

// Run an application build task.
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

// Get a command action from an HTTP verb and request URL.
func (b *CommandAdapter) CommandAction(verb string, url *url.URL) (*Action, *StatusError) {
  var a *Action = &Action{Verb: verb, Url: url}
  switch verb {
    case http.MethodPut:
      a.Operation = OperationCreate
    case http.MethodGet:
      a.Operation = OperationRead
    case http.MethodPost:
      a.Operation = OperationUpdate
    case http.MethodDelete:
      a.Operation = OperationDelete
    default:
      return nil, CommandError(http.StatusMethodNotAllowed, "")
  }

  parse := func(path string) {
    a.Path = path
    if a.Path != "" {
      a.Parts = strings.Split(strings.TrimSuffix(a.Path, SLASH), SLASH)
      a.Type = a.Parts[0]
      if len(a.Parts) > 1 {
        a.Context = a.Parts[1]
      }
      if len(a.Parts) > 2 {
        a.Target = a.Parts[2]
      }
      if len(a.Parts) > 3 {
        a.Action = a.Parts[3]
      }
      if len(a.Parts) > 4 {
        a.Item = SLASH + strings.Join(a.Parts[4:], SLASH)
        // Respect input trailing slash used to indicate
        // operations on a directory
        if strings.HasSuffix(a.Path, SLASH) {
          a.Item += SLASH
        }
      }
    }
  }

  parse(url.Path)

  fmt.Printf("%#v\n", a)

  return a, nil
}

type ActionDefinition struct {
  MethodName string
  // Received will be the command adapter
  Receiver reflect.Value
  // Method is the function to invoke
  Method reflect.Method
  // Arity for arguments
  ArityIn int
  // Arity for return value
  ArityOut int
  // HTTP status code to use on success
  Status int
  // Build function invocation arguments
  Arguments func(action *Action) []reflect.Value
}

type ActionResult struct {
  *Action
  *ActionDefinition
  Data interface{}
  Error *StatusError
  Status int
}

func (b *CommandAdapter) Handler(act *Action) (*Action, *ActionDefinition) {
  var m reflect.Method
  receiver := reflect.ValueOf(b)
  t := reflect.TypeOf(b)

  for a, def := range ActionMap {
    // fmt.Printf("test for match: %#v\n", a)
    if a.Match(act) {
      def.Receiver = receiver
      m, _ = t.MethodByName(def.MethodName)
      def.Method = m
      def.ArityIn = m.Type.NumIn()
      def.ArityOut = m.Type.NumOut()
      return a, def
    }
  }
  return nil, nil
}

func (b *CommandAdapter) Execute(act *Action) (*ActionResult, *StatusError) {
  action, def := b.Handler(act)

  // No definition found
  if def == nil {
    return nil, CommandError(http.StatusNotFound, "")
  }

  var args []reflect.Value = make([]reflect.Value, 0)
  // Docs say that a Method does not need the receiver argument
  // but it appears we need it
  args = append(args, def.Receiver)

  // Additional arguments to pass after we add the received
  if def.Arguments != nil {
    fn := def.Arguments(act)
    args = append(args, fn...)
  }

  // fmt.Printf("args:%#v\n", args)

  // TODO: work out correct args

  // Call the method
  res := def.Method.Func.Call(args)

  // Check return value arity
  if len(res) == 0 || len(res) > 2 {
    return nil, CommandError(
      http.StatusInternalServerError, "Invalid command return value arity")
  }

  // Setup the result object
  var result *ActionResult = &ActionResult{ActionDefinition: def, Action: action}
  result.Status = result.ActionDefinition.Status

  // Get the underlying return values and test for error response
  var retval []interface{}
  for _, val := range res {
    v := val.Interface()
    if ex, ok := v.(*StatusError); ok {
      // Mark result with error
      result.Error = ex
    }
    retval = append(retval, v)
  }

  // Assign the method call return value as the result data
  result.Data = retval[0]

  // Done :)
  return result, result.Error
}

func init() {
  // GET /
  ActionMap[&Action{Operation: OperationRead}] =
    &ActionDefinition{
      MethodName: "ListContainers",
      Status: http.StatusOK}
  // GET /templates
  ActionMap[&Action{Operation: OperationRead, Type: "templates"}] =
    &ActionDefinition{
      MethodName: "ListApplicationTemplates",
      Status: http.StatusOK}
  // GET /jobs
  ActionMap[&Action{Operation: OperationRead, Type: "jobs"}] =
    &ActionDefinition{
      MethodName: "ListJobs",
      Status: http.StatusOK}
  // GET /jobs/{id}
  ActionMap[&Action{Operation: OperationRead, Type: "jobs", Context: "*"}] =
    &ActionDefinition{
      MethodName: "ReadJob",
      Arguments: func(action *Action) []reflect.Value {
        var args []reflect.Value
        args = append(args, reflect.ValueOf(action.Context))
        return args
      },
      Status: http.StatusOK}
  // DELETE /jobs/{id}
  ActionMap[&Action{Operation: OperationRead, Type: "jobs", Context: "*"}] =
    &ActionDefinition{
      MethodName: "AbortJob",
      Arguments: func(action *Action) []reflect.Value {
        var args []reflect.Value
        args = append(args, reflect.ValueOf(action.Context))
        return args
      },
      Status: http.StatusOK}
}
