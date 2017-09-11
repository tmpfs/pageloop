package adapter

import (
  "net/url"
  "strings"
  . "github.com/tmpfs/pageloop/util"
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

func (act *Action) ActionOnly() bool {
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
      act.Action = act.Parts[3]
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
}

func (act *Action) ContextMatch(in *Action) bool {
  return (act.Wildcard(act.Context) || act.Context == in.Context)
}

func (act *Action) TargetMatch(in *Action) bool {
  return (act.Wildcard(act.Target) || act.Target == in.Target)
}

func (act *Action) ActionMatch(in *Action) bool {
  return (act.Wildcard(act.Action) || act.Action == in.Action)
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
    if act.ActionOnly() && in.ActionOnly() && act.ContextMatch(in) && act.ActionMatch(in) {
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
