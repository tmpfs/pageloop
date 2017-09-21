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

// Read file content.
func (s *FileService) ReadSource(file *File, reply *ServiceReply) *StatusError {
  if _, _, f, err := s.lookup(file); err != nil {
    return err
  } else {
    reply.Reply = f.Source(false)
  }
  return nil
}

// Read raw file content (includes frontmatter).
func (s *FileService) ReadSourceRaw(file *File, reply *ServiceReply) *StatusError {
  if _, _, f, err := s.lookup(file); err != nil {
    return err
  } else {
    reply.Reply = f.Source(true)
  }
  return nil
}

// Save file content.
func (s *FileService) Save(file *File, reply *ServiceReply) *StatusError {
  if _, app, f, err := s.lookup(file); err != nil {
    return err
  } else {

    if err := app.Update(f, file.Source(false)); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }

    if f.Page() != nil {
      reply.Reply = f.Page()
      return nil
    }
    reply.Reply = f
  }
  return nil

  /*
  if app, file, err :=  b.ReadFile(c, a, f); err != nil {
    return nil, err
  } else {
    if file, err := b.CommandExecute.UpdateFile(app, file, content); err != nil {
      return nil, err
    } else {
      if file.Page() != nil {
        return file.Page(), nil
      }
      return file, nil
    }
  }
  */
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
