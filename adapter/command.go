// Package adapter provides a command adapter for interfacing
// network requests with the underlying model.
package adapter

import (
  //"fmt"
  "net/url"
  "net/http"
  "reflect"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/util"
)

/*
// Handler for asynchronous background tasks.
type TaskJobComplete struct {}

func (tj *TaskJobComplete) Done(err error, job *Job) {
  // TODO: send reply to the client over websocket
  fmt.Printf("[job:%d] completed %s\n", job.Number, job.Id)
  Jobs.Stop(job)
}
*/

// Abstraction that allows many different interfaces to
// the data model whether it is a string command interpreter,
// REST API endpoints, JSON RPC or any other bridge to the
// outside world.
//
// For simplicity with access over HTTP this implementation always
// returns errors with an associated HTTP status code.
//
// Usage of this package involves creating a new action by calling HttpAction(),
// the resulting action should then be passed to Find() to see if the route matches
// followed by a call to Execute() to invoke the command function.
type CommandAdapter struct {
  *CommandExecute
  Name string
  Version string
  Host *Host
  Mountpoints *MountpointManager
}

func NewCommandAdapter(name string, version string, host *Host, mountpoints *MountpointManager) *CommandAdapter {
  a := &CommandAdapter{Name: name, Version: version, Host: host, Mountpoints: mountpoints}
  a.CommandExecute = &CommandExecute{CommandAdapter: a}
  return a
}

// Get a command action from an HTTP verb and request URL.
func (b *CommandAdapter) HttpAction(verb string, url *url.URL) (*Action, *StatusError) {
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

  a.Parse(url.Path)

  return a, nil
}

// Find the command definition for an incoming action.
//
// Must be invoked on an action before passing the action to Execute.
//
// If a match is found then the action is populated with a route action,
// command definition and arguments list.
//
// If no match is found an error is returned.
func (b *CommandAdapter) Find(act *Action) (*ActionMap, *StatusError) {
  mapping := b.handler(act)

  // No definition found
  if mapping == nil {
    return nil, CommandError(http.StatusNotFound, "")
  }

  b.initArguments(act, mapping)
  return mapping, nil
}

// Find the command definition by service method name.
func (b *CommandAdapter) FindService(method string, act *Action) (*ActionMap, *StatusError) {
  if mapping, ok := Services[method]; ok {
    var m reflect.Method
    receiver := reflect.ValueOf(b)
    t := reflect.TypeOf(b)
    def := mapping.CommandDefinition
    def.Receiver = receiver
    m, _ = t.MethodByName(def.MethodName)
    def.Method = m
    b.initArguments(act, mapping)
    return mapping, nil
  }
  return nil, CommandError(http.StatusNotFound, "")
}

// Execute an action.
//
// You should have already invoked Find() so that the action has been
// assigned a route, command definition and arguments list.
//
// The method associated with the command definition is invoked and it's
// return values are converted to an action result which is returned to
// the caller.
func (b *CommandAdapter) Execute(act *Action) (*ActionResult, *StatusError) {
  if act.Command == nil {
    return nil, CommandError(
      http.StatusInternalServerError, "Action has no command, call find before execution")
  }

  def := act.Command
  args := act.Arguments
  // Call the method
  res := def.Method.Func.Call(args)

  // Setup the result object
  var result *ActionResult = &ActionResult{CommandDefinition: def, Action: act}
  result.Status = result.CommandDefinition.Status

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
  result.Data = retval[def.Index]

  // fmt.Printf("%#v\n", result)

  // Done :)
  return result, result.Error
}

// Mutate a matched action such that it is rewritten to a different method call.
func (b *CommandAdapter) Mutate(act *Action, name string) *Action {
  var m reflect.Method
  receiver := reflect.ValueOf(b)
  t := reflect.TypeOf(b)
  for _, mapping := range Routes {
    def := mapping.CommandDefinition
    if def.MethodName == name {
      def.Receiver = receiver
      m, _ = t.MethodByName(def.MethodName)
      def.Method = m
      b.initArguments(act, mapping)
      return act
    }
  }
  return nil
}

// PRIVATE

func (b *CommandAdapter) initArguments(act *Action, mapping *ActionMap) {
  def := mapping.CommandDefinition

  var args []reflect.Value = make([]reflect.Value, 0)
  // Docs say that a Method does not need the receiver argument
  // but it appears we need it
  args = append(args, def.Receiver)

  // Additional arguments to pass after we add the receiver
  if def.Arguments != nil {
    fn := def.Arguments(b, act)
    args = append(args, fn...)
  }

  act.Route = mapping.Action
  act.Command = mapping.CommandDefinition
  act.Arguments = args
}

// This is the route matching stuff. Takes the incoming action
// and finds the first matching route action.
//
// Returns the mapping containing the route action and command definition,
// if the request action does not match any route nil is returned.
func (b *CommandAdapter) handler(act *Action) *ActionMap {
  var m reflect.Method
  receiver := reflect.ValueOf(b)
  t := reflect.TypeOf(b)
  for _, mapping := range Routes {
    a := mapping.Action
    def := mapping.CommandDefinition
    if a.Match(act) {
      def.Receiver = receiver
      m, _ = t.MethodByName(def.MethodName)
      def.Method = m
      return mapping
    }
  }
  return nil
}
