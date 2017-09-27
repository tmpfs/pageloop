package service

import(
   //"fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type ContainerService struct {
  // Reference to the host
  Host *Host

  // Reference to the mountpoint manager
  Mountpoints *MountpointManager
}

type ContainerRequest struct {
  Name string `json:"name"`
}

// Read a container.
func (s *ContainerService) Read(container *ContainerRequest, reply *ServiceReply) *StatusError {
  if c, err := LookupContainer(s.Host, container); err != nil {
    return err
  } else {
    reply.Reply = c
  }
  return nil
}

// Create application.
func (s *ContainerService) CreateApp(req *ApplicationRequest, reply *ServiceReply) *StatusError {
  var cr *ContainerRequest = &ContainerRequest{Name: req.Container}
  if container, err := LookupContainer(s.Host, cr); err != nil {
    return err
  } else {
    // TODO: do not allow creating apps on non-user containers!
    if container.Protected {
      return CommandError(http.StatusForbidden, "Cannot create applications in a protected container.")
    }

    app := req.ToApplication(container)

    existing := container.GetByName(app.Name)
    if existing != nil {
      return CommandError(http.StatusPreconditionFailed, "Application %s already exists", app.Name)
    }

    // Get mountpoint URL.
    app.Url = app.MountpointUrl(container)

    // Mountpoint exists.
    exists := s.Mountpoints.HasMountpoint(app.Url)
    if exists {
      return CommandError(http.StatusPreconditionFailed, "Mountpoint URL %s already exists", app.Url)
    }

    // Create and save a mountpoint for the application.
    if mountpoint, err := s.Mountpoints.CreateMountpoint(app); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    } else {
      // Handle creating from a template.
      if req.Template != nil {
        // Find the template application.
        if source, err := s.Host.LookupTemplate(req.Template); err != nil {
          return CommandError(http.StatusBadRequest, err.Error())
        } else {
          // Copy template source files.
          if err := app.CopyApplicationTemplate(source); err != nil {
            return CommandError(http.StatusInternalServerError, err.Error())
          }
        }
      }

      // Load and publish the app source files, note that we get a new application back
      // after loading the mountpoint.
      if app, err = s.Mountpoints.LoadMountpoint(*mountpoint, container); err != nil {
        return CommandError(http.StatusInternalServerError, err.Error())
      } else {
        // Reply with the new application reference
        reply.Reply = app
        reply.Status = http.StatusCreated
      }
    }
  }

  return nil
}

func LookupContainer(host *Host, container *ContainerRequest) (*Container, *StatusError) {
  c := host.GetByName(container.Name)
  if c == nil {
    return nil, CommandError(http.StatusNotFound, "Container %s not found", container.Name)
  }
  return c, nil
}
