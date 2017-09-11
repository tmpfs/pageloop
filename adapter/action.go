package adapter

import (
  "fmt"
  "net/http"
  "net/url"
  "reflect"
  "strings"
  . "github.com/tmpfs/pageloop/util"
)

const(
  // Basic CRUD operations
  OperationCreate = iota
  OperationRead
  OperationUpdate
  OperationDelete
)

var(
  ActionList []*ActionMap
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
  // A filter operation for the request.
  Filter string
  // An item, may contain slashes.
  Item string

  // Populated once find has been called.
  Arguments []reflect.Value
  Route *Action
  Command *CommandDefinition
}

// An action definition defines the receiving command function for an incoming action.
type CommandDefinition struct {
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
  Arguments func(b *CommandAdapter, action *Action) []reflect.Value
}

// Combines action routing information with the command definition.
type ActionMap struct {
  *Action
  *CommandDefinition
}

// An ActionResult is the response returned after command invocation.
type ActionResult struct {
  *Action
  *CommandDefinition
  Data interface{}
  Error *StatusError
  Status int
}

func NewAction(op int, path string) *Action {
  act := &Action{Operation: op}
  act.Parse(path)
  return act
}

func (act *Action) IsRoot() bool {
  return len(act.Parts) == 0
}

func (act *Action) TypeOnly() bool {
  return len(act.Parts) == 1
}

func (act *Action) ContextOnly() bool {
  return len(act.Parts) == 2
}

func (act *Action) TargetOnly() bool {
  return len(act.Parts) == 3
}

func (act *Action) FilterOnly() bool {
  return len(act.Parts) == 4
}

func (act *Action) MatchType(in *Action) bool {
  return act.Type == in.Type
}

func (act *Action) Wildcard(val string) bool {
  return val == "*"
}

func (act *Action) Parse(path string) {
  act.Path = path
  if act.Path != "" {
    path := strings.TrimPrefix(act.Path, SLASH)
    path = strings.TrimSuffix(path, SLASH)
    act.Parts = strings.Split(path, SLASH)
    act.Type = act.Parts[0]
    if len(act.Parts) > 1 {
      act.Context = act.Parts[1]
    }
    if len(act.Parts) > 2 {
      act.Target = act.Parts[2]
    }
    if len(act.Parts) > 3 {
      act.Filter = act.Parts[3]
    }
    if len(act.Parts) > 4 {
      act.Item = SLASH + strings.Join(act.Parts[4:], SLASH)
      // Respect input trailing slash used to indicate
      // operations on a directory
      if strings.HasSuffix(act.Path, SLASH) {
        act.Item += SLASH
      }
    }
  }

  // So that trailing slash with no URL will match
  // the filter
  if act.Item == SLASH {
    act.Item = ""
  }
}

func (act *Action) ContextMatch(in *Action) bool {
  return (act.Wildcard(act.Context) || act.Context == in.Context)
}

func (act *Action) TargetMatch(in *Action) bool {
  return (act.Wildcard(act.Target) || act.Target == in.Target)
}

func (act *Action) FilterMatch(in *Action) bool {
  return (act.Wildcard(act.Filter) || act.Filter == in.Filter)
}

func (act *Action) ItemMatch(in *Action) bool {
  item := strings.TrimPrefix(act.Item, SLASH)
  return (act.Wildcard(item) || act.Item == in.Item)
}

func (act *Action) Match(in *Action) bool {
  if act.Operation != in.Operation {
    return false
  }

  // Root match
  if act.IsRoot() && in.IsRoot() {
    return true
  }

  if act.TypeOnly() && in.TypeOnly() && act.Type == in.Type {
    return true
  }

  // Got a type match
  if act.MatchType(in) {
    // Deal with context only
    if act.ContextOnly() && in.ContextOnly() && act.ContextMatch(in) {
      return true
    }

    // Deal with target only
    if act.TargetOnly() && in.TargetOnly() && act.ContextMatch(in) && act.TargetMatch(in) {
      return true
    }

    // Deal with action only
    if act.FilterOnly() && in.FilterOnly() && act.ContextMatch(in) && act.TargetMatch(in) && act.FilterMatch(in) {
      return true
    }

    // println("testing on item")

    // Final path portion is an item, that is a file or page URL potentially
    // containing the slash character.
    if act.Item != "" && act.ItemMatch(in) {
      return true
    }
  }

  return false
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

  a.Parse(url.Path)

  fmt.Printf("%#v\n", a)

  return a, nil
}

func (b *CommandAdapter) Handler(act *Action) *ActionMap {
  var m reflect.Method
  receiver := reflect.ValueOf(b)
  t := reflect.TypeOf(b)

  fmt.Printf("TEST ON ACTION: %#v\n", act)

  for _, mapping := range ActionList {
    a := mapping.Action
    def := mapping.CommandDefinition
    fmt.Printf("test for match pattern: %#v\n", a.Path)
    fmt.Printf("test for match input: %#v\n", act.Path)
    if a.Match(act) {
      println("got method match: " + def.MethodName)
      def.Receiver = receiver
      m, _ = t.MethodByName(def.MethodName)
      def.Method = m
      def.ArityIn = m.Type.NumIn()
      def.ArityOut = m.Type.NumOut()
      return mapping
    }
  }
  return nil
}

// Find the command definition for an incoming action.
//
// If a match is found then the action is populated with a route action,
// command definition and arguments list.
//
// If no match is found an error is returned.
func (b *CommandAdapter) Find(act *Action) (*ActionMap, *StatusError) {
  mapping := b.Handler(act)

  // No definition found
  if mapping == nil {
    return nil, CommandError(http.StatusNotFound, "")
  }

  def := mapping.CommandDefinition

  var args []reflect.Value = make([]reflect.Value, 0)
  // Docs say that a Method does not need the receiver argument
  // but it appears we need it
  args = append(args, def.Receiver)

  // Additional arguments to pass after we add the received
  if def.Arguments != nil {
    fn := def.Arguments(b, act)
    args = append(args, fn...)
  }

  act.Route = mapping.Action
  act.Command = mapping.CommandDefinition
  act.Arguments = args

  return mapping, nil
}

// Execute an action.
//
// This attempts to find an action and action definition that match the
// supplied action. It is an error if the action does not match a mapped
// route otherwise the action is executed and the result is returned back
// to the caller as an action result.
func (b *CommandAdapter) Execute(act *Action) (*ActionResult, *StatusError) {
  if act.Command == nil {
    return nil, CommandError(
      http.StatusInternalServerError, "Action has no command, call find before execution")
  }

  def := act.Command
  args := act.Arguments
  // Call the method
  res := def.Method.Func.Call(args)

  // Check return value arity
  if len(res) == 0 || len(res) > 2 {
    return nil, CommandError(
      http.StatusInternalServerError, "Invalid command return value arity")
  }

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
  result.Data = retval[0]

  // Done :)
  return result, result.Error
}

// Initialize the action list with route actions and action definitions.
func init() {
  push := func(action *Action, def *CommandDefinition) {
    ActionList = append(ActionList, &ActionMap{Action: action, CommandDefinition: def})
  }

  contextArg := func(b *CommandAdapter, action *Action) []reflect.Value {
    var args []reflect.Value
    args = append(args, reflect.ValueOf(action.Context))
    return args
  }

  containerArg := func(b *CommandAdapter, action *Action) []reflect.Value {
    var args []reflect.Value
    args = append(args, reflect.ValueOf(action.Context))
    return args
  }

  applicationArg := func(b *CommandAdapter, action *Action) []reflect.Value {
    var args []reflect.Value
    args = append(args, reflect.ValueOf(action.Context), reflect.ValueOf(action.Target))
    return args
  }

  fileArg := func(b *CommandAdapter, action *Action) []reflect.Value {
    var args []reflect.Value
    args = append(
      args,
      reflect.ValueOf(action.Context),
      reflect.ValueOf(action.Target),
      reflect.ValueOf(action.Item))
    return args
  }

  // GET /
  push(NewAction(OperationRead, ""),
    &CommandDefinition{MethodName: "Meta", Status: http.StatusOK})
  // GET /templates
  push(NewAction(OperationRead, "/templates"),
    &CommandDefinition{MethodName: "ListApplicationTemplates", Status: http.StatusOK})
  // GET /jobs
  push(NewAction(OperationRead, "/jobs"),
    &CommandDefinition{MethodName: "ListJobs", Status: http.StatusOK})
  // GET /jobs/{id}
  push(NewAction(OperationRead, "/jobs/*"),
    &CommandDefinition{MethodName: "ReadJob", Status: http.StatusOK, Arguments: contextArg})
  // DELETE /jobs/{id}
  push(NewAction(OperationDelete, "/jobs/*"),
    &CommandDefinition{MethodName: "AbortJob", Status: http.StatusOK, Arguments: contextArg})
  // GET /apps
  push(NewAction(OperationRead, "/apps"),
    &CommandDefinition{MethodName: "ReadHost", Status: http.StatusOK})
  // GET /apps/{container}
  push(NewAction(OperationRead, "/apps/*"),
    &CommandDefinition{MethodName: "ReadContainer", Status: http.StatusOK, Arguments: containerArg})
  // GET /apps/{container}/{application}
  push(NewAction(OperationRead, "/apps/*/*"),
    &CommandDefinition{MethodName: "ReadApplication", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/files
  push(NewAction(OperationRead, "/apps/*/*/files"),
    &CommandDefinition{MethodName: "ReadApplicationFiles", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/pages
  push(NewAction(OperationRead, "/apps/*/*/pages"),
    &CommandDefinition{MethodName: "ReadApplicationPages", Status: http.StatusOK, Arguments: applicationArg})
  // GET /apps/{container}/{application}/files/{url}
  push(NewAction(OperationRead, "/apps/*/*/files/*"),
    &CommandDefinition{MethodName: "ReadFile", Status: http.StatusOK, Arguments: fileArg})
  // GET /apps/{container}/{application}/pages/{url}
  push(NewAction(OperationRead, "/apps/*/*/pages/*"),
    &CommandDefinition{MethodName: "ReadPage", Status: http.StatusOK, Arguments: fileArg})
}
