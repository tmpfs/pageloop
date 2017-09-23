package service

import(
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type FileService struct {
  Host *Host
}

// Read a file.
func (s *FileService) Read(file *File, reply *ServiceReply) *StatusError {
  if _, _, f, err := LookupFile(s.Host, file, false); err != nil {
    return err
  } else {
    reply.Reply = f
  }
  return nil
}

// Move a file.
func (s *FileService) Move(file *File, reply *ServiceReply) *StatusError {
  if file.Destination == "" {
    return CommandError(http.StatusBadRequest, "No destination for move operation")
  }
  if _, app, f, err := LookupFile(s.Host, file, false); err != nil {
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
  if _, _, f, err := LookupFile(s.Host, file, false); err != nil {
    return err
  } else {
    reply.Reply = f.Source(false)
  }
  return nil
}

// Read raw file content (includes frontmatter).
func (s *FileService) ReadSourceRaw(file *File, reply *ServiceReply) *StatusError {
  if _, _, f, err := LookupFile(s.Host, file, false); err != nil {
    return err
  } else {
    reply.Reply = f.Source(true)
  }
  return nil
}

// Save file content.
func (s *FileService) Save(file *File, reply *ServiceReply) *StatusError {
  if _, app, f, err := LookupFile(s.Host, file, false); err != nil {
    return err
  } else {

    // File content from string value
    if file.Value != "" {
      file.Bytes([]byte(file.Value))
      file.Value = ""
    }

    var content []byte = file.Source(false)

    if err := app.Update(f, content); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }

    if f.Page() != nil {
      reply.Reply = f.Page()
      return nil
    }

    reply.Reply = f
  }
  return nil
}

// Create a new file and publish it, the file cannot already exist on disc.
func (s *FileService) Create(file *File, reply *ServiceReply) *StatusError {
  if _, app, _, err := LookupFile(s.Host, file, true); err != nil {
    return err
  } else {
    var exists *File = app.Urls[file.Url]

    if exists != nil {
      return CommandError(http.StatusConflict,"File already exists %s", file.Url)
    }
    if app.ExistsConflict(file.Url) {
      return CommandError(http.StatusConflict,"File already exists, publish conflict on %s", file.Url)
    }

    if file, err := app.Create(file.Url, file.Source(false)); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    } else {
      reply.Reply = file
      reply.Status = http.StatusCreated
    }
  }

  return nil
}

// Create a file from a template.
func (s *FileService) CreateFileTemplate(file *File, reply *ServiceReply) *StatusError {
  template := file.Template
  if template == nil {
    return CommandError(http.StatusBadRequest, "No template given")
  }
  if _, _, _, err := LookupFile(s.Host, file, true); err != nil {
    return err
  } else {
    if tpl, err := s.Host.LookupTemplateFile(template); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    } else {
      if tpl == nil {
        return CommandError(http.StatusNotFound, "Template file %s does not exist", template.File)
      }
      file.Bytes(tpl.Source(true))
      return s.Create(file, reply)
    }
  }

  return nil
}

func LookupFile(host *Host, f *File, appOnly bool) (*Container, *Application, *File, *StatusError) {
  // Parse file URL references
  if f.Ref != "" {
    if ref, err := ParseFileUrl(f.Ref); err != nil {
      return nil, nil, nil, CommandError(http.StatusInternalServerError, err.Error())
    } else {
      f.Owner = &Application{Name: ref.Application, Container: &Container{Name: ref.Container}}
      f.Url = ref.Url
    }
  }

  if f.Owner == nil {
    return nil, nil, nil, CommandError(
      http.StatusNotFound, "File %s missing owner application (detached file)", f.Url)
  }

  if f.Owner.Container == nil {
    return nil, nil, nil, CommandError(
      http.StatusNotFound, "Application %s missing container (detached app)", f.Owner.Name)
  }

  c := host.GetByName(f.Owner.Container.Name)
  if c == nil {
    return nil, nil, nil, CommandError(http.StatusNotFound, "Container %s not found", f.Owner.Container.Name)
  }

  app := c.GetByName(f.Owner.Name)
  if app == nil {
    return nil, nil, nil, CommandError(http.StatusNotFound, "Application %s not found", f.Owner.Name)
  }

  if appOnly {
    return c, app, nil, nil
  }

  file := app.Urls[f.Url]
  if file == nil {
    return nil, nil, nil, CommandError(http.StatusNotFound, "File %s not found", f.Url)
  }
  return c, app, file, nil
}
