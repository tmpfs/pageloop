package adapter

import (
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

// A command action is a simple representation of a command invocation
// it can be used to execute a command without any object references.
//
// Path references take the form:
//
// /{type}?/{context}?/{target}?/{action}?/{item}?
//
// Where item is a trailer that may include slashes to represent a file URL.
//
// The context part corresponds to a container and the target part corresponds
// to an application.
//
// If a definition maps a part using the wildcard (*) it will match any string.
type Action struct {
  // Source HTTP verb that is translated to an operation constant.
  Verb string
  // A request URL.
  Url *url.URL
  // The path for the request.
  Path string
  // Parsed path parts split on a slash.
  Parts []string
  // The CRUD operation to perform.
  Operation int

  // The operation type, cannot be a wildcard.
  Type string
  // Context for the operation. May be a container reference, job number etc.
  Context string
  // Target for the operation, typically an application.
  Target string
  // A filter operation for the request.
  Filter string
  // An item, may contain slashes.
  Item string

  // Populated once Find has been called.

  // List of arguments to pass to the command function.
  Arguments []reflect.Value
  // The route action that triggered the match on this action.
  Route *Action
  // The command definition used for method invocation.
  Command *CommandDefinition
}

// An action definition defines the receiving command function for an incoming action.
//
// Command functions typically have a return value arity of two:
//
// func() (interface{}, *StatusError)
//
// And the first return value is the one used for the data on the action result unless
// an index has been specified on the command definition in which case it must not be out
// of bounds.
type CommandDefinition struct {
  MethodName string
  // Receiver will be the command adapter.
  Receiver reflect.Value
  // Method is the function to invoke.
  Method reflect.Method
  // HTTP status code to use on success.
  Status int
  // Function called to build the initial argument list.
  // For most invocations these will be sufficient but when creating and
  // updating arguments may need to be added by the caller in which
  // case they should call Push() on the action after a call to Find() and
  // before calling Execute().
  Arguments func(b *CommandAdapter, action *Action) []reflect.Value
  // An index into the command return values to use as the result data.
  Index int
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

// Create a new action for the given operation and path.
func NewAction(op int, path string) *Action {
  act := &Action{Operation: op}
  act.Parse(path)
  return act
}

// Add an argument to the list of arguments that will be passed
// to the command method on invocation.
//
// This wraps the past type as a reflect.Value before appending
// to the arguments slice.
func (act *Action) Push(t interface{}) {
  act.Arguments = append(act.Arguments, reflect.ValueOf(t))
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

// Determine if an incoming action matches this action.
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

    // Final path portion is an item, that is a file or page URL potentially
    // containing the slash character.
    if act.Item != "" && act.ItemMatch(in) {
      return true
    }
  }
  return false
}
