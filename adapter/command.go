// Package adapter provides a command adapter for interfacing
// network requests with the underlying model.
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

