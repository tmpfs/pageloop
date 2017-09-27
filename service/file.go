package service

import(
  "fmt"
  "strings"
  "net/http"
  "net/url"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type FileRef struct {
  Container string
  Application string
  Url string
}

type FileRequest struct {
  Name string `json:"name"`

  // This is a slash separated path to the source file
  // relative to the application base, it will start with
  // a leading slash. Directories will always have a trailing
  // slash.
  Url string `json:"url"`

  // Destination for file move operations
  Destination string `json:"destination,omitempty"`

	// A source template for this file
	Template *ApplicationTemplate `json:"template,omitempty"`

  // A reference to a file in the form: file://{container}/{application}#{url}
  Ref string `json:"ref,omitempty"`

  // An input value for the file content, passed in when creating or
  // updating files that are not binary
  Value string `json:"value,omitempty"`

  // Value specified as a byte slice, used when creating files from template.
  Bytes []byte
}

func (req *FileRequest) ToFile () *File {
  return &File{Name: req.Name, Url: req.Url}
}

type FileService struct {
  Host *Host
}

// Read a file.
func (s *FileService) Read(file *FileRequest, reply *ServiceReply) *StatusError {
  if _, _, f, err := LookupFile(s.Host, file, false); err != nil {
    return err
  } else {
    reply.Reply = f
  }
  return nil
}

// Move a file.
func (s *FileService) Move(file *FileRequest, reply *ServiceReply) *StatusError {
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
func (s *FileService) ReadSource(file *FileRequest, reply *ServiceReply) *StatusError {
  if _, _, f, err := LookupFile(s.Host, file, false); err != nil {
    return err
  } else {
    reply.Reply = f.Source(false)
  }
  return nil
}

// Read raw file content (includes frontmatter).
func (s *FileService) ReadSourceRaw(file *FileRequest, reply *ServiceReply) *StatusError {
  if _, _, f, err := LookupFile(s.Host, file, false); err != nil {
    return err
  } else {
    reply.Reply = f.Source(true)
  }
  return nil
}

// Save file content.
func (s *FileService) Save(req *FileRequest, reply *ServiceReply) *StatusError {
  if _, app, file, err := LookupFile(s.Host, req, false); err != nil {
    return err
  } else {

    // File content from string value
    if req.Value != "" {
      file.Bytes([]byte(req.Value))
    }

    var content []byte = file.Source(false)

    if err := app.Update(file, content); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }

    if file.Page() != nil {
      reply.Reply = file.Page()
      return nil
    }

    reply.Reply = file
  }
  return nil
}

// Create a new file and publish it, the file cannot already exist on disc.
func (s *FileService) Create(req *FileRequest, reply *ServiceReply) *StatusError {
  if _, app, _, err := LookupFile(s.Host, req, true); err != nil {
    return err
  } else {
    var exists *File = app.Urls[req.Url]

    if exists != nil {
      return CommandError(http.StatusConflict,"File already exists %s", req.Url)
    }
    if app.ExistsConflict(req.Url) {
      return CommandError(http.StatusConflict,"File already exists, publish conflict on %s", req.Url)
    }

    content := req.Bytes
    if req.Value != "" {
      content = []byte(req.Value)
    }

    if file, err := app.Create(req.Url, content); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    } else {
      reply.Reply = file
      reply.Status = http.StatusCreated
    }
  }

  return nil
}

// Create a file from a template.
func (s *FileService) CreateFileTemplate(req *FileRequest, reply *ServiceReply) *StatusError {
  template := req.Template
  if template == nil {
    return CommandError(http.StatusBadRequest, "No template given")
  }
  if _, _, _, err := LookupFile(s.Host, req, true); err != nil {
    return err
  } else {
    if tpl, err := s.Host.LookupTemplateFile(template); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    } else {
      if tpl == nil {
        return CommandError(http.StatusNotFound, "Template file %s does not exist", template.File)
      }
      req.Bytes = tpl.Source(true)
      return s.Create(req, reply)
    }
  }

  return nil
}

func ParseFileUrl(uri string) (ref *FileRef, err error) {
  var u *url.URL
  if u, err = url.Parse(uri); err != nil {
    return
  }
  parts := strings.Split(u.Path, "/")
  if len(parts) != 2 {
    err = fmt.Errorf("Invalid file reference %s", uri)
    return
  }
  ref = &FileRef{Container: u.Host, Application: parts[1], Url: u.Fragment}
  return
}

func LookupFile(host *Host, req *FileRequest, appOnly bool) (*Container, *Application, *File, *StatusError) {
  f := req.ToFile()
  // Parse file URL references
  if ref, err := ParseFileUrl(req.Ref); err != nil {
    return nil, nil, nil, CommandError(http.StatusInternalServerError, err.Error())
  } else {
    f.Owner = &Application{Name: ref.Application, Container: &Container{Name: ref.Container}}
    f.Url = ref.Url
  }

  /*
  if f.Owner == nil {
    return nil, nil, nil, CommandError(
      http.StatusNotFound, "File %s missing owner application (detached file)", f.Url)
  }

  if f.Owner.Container == nil {
    return nil, nil, nil, CommandError(
      http.StatusNotFound, "Application %s missing container (detached app)", f.Owner.Name)
  }
  */

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
