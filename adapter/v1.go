package adapter

// @deprecated - all deprecated from v1

import (
  //"fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

// Move a file.

// BROKEN in v1 now!
/*
func (b *CommandAdapter) MoveFile(a *Application, f *File, dest string) *StatusError {
  if err := a.Move(f, dest); err != nil {
    return CommandError(http.StatusInternalServerError, err.Error())
  }
  return nil
}
*/

func (b *CommandAdapter) ListApplications(c *Container) []*Application {
  return c.Apps
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

