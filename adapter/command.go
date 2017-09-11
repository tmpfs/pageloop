// Package adapter provides a command adapter for interfacing
// network requests with the underlying model.
package adapter

import (
  "fmt"
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

