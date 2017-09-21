package service

import(
  // "fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type FileService struct {
  Host *Host
}

// Move a file.
func (s *FileService) Move(file *File, reply *ServiceReply) *StatusError {
  if file.Destination == "" {
    return CommandError(http.StatusBadRequest, "No destination for move operation")
  }
  if _, app, f, err := s.lookup(file); err != nil {
    return err
  } else {
    if err := app.Move(f, file.Destination); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }
    reply.Reply = f
  }
  return nil
}

// Private

func (s *FileService) lookup(f *File) (*Container, *Application, *File, *StatusError) {
  if f.Owner == nil {
    return nil, nil, nil, CommandError(
      http.StatusNotFound, "File %s missing owner application (detached file)", f.Url)
  }

  if f.Owner.Container == nil {
    return nil, nil, nil, CommandError(
      http.StatusNotFound, "Application %s missing container (detached app)", f.Owner.Name)
  }

  c := s.Host.GetByName(f.Owner.Container.Name)
  if c == nil {
    return nil, nil, nil, CommandError(http.StatusNotFound, "Container %s not found", f.Owner.Container.Name)
  }

  app := c.GetByName(f.Owner.Name)
  if app == nil {
    return nil, nil, nil, CommandError(http.StatusNotFound, "Application %s not found", f.Owner.Name)
  }

  file := app.Urls[f.Url]
  if file == nil {
    return nil, nil, nil, CommandError(http.StatusNotFound, "File %s not found", f.Url)
  }
  return c, app, file, nil
}
