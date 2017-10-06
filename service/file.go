package service

import(
  // "fmt"
  // "strings"
  "net/http"
  // "net/url"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type FileRef struct {
  Container string
  Application string
  Url string
}

type FileReferenceRequest struct {
  // A reference to a file in the form: file://pageloop.com/{container}/{application}#{url}
  Ref string `json:"ref,omitempty"`
}

type FileMoveRequest struct {
  // A reference to a file in the form: file://pageloop.com/{container}/{application}#{url}
  Ref string `json:"ref,omitempty"`
  // Destination for file move operations
  Destination string `json:"destination,omitempty"`
}

type FileContentRequest struct {
  // A reference to a file in the form: file://pageloop.com/{container}/{application}#{url}
  Ref string `json:"ref,omitempty"`

  // An input value for the file content, passed in when creating or
  // updating files that are not binary
  Value string `json:"value,omitempty"`

  // Value specified as a byte slice, when receiving POST and PUT requests
  Bytes []byte
}

type FileTemplateRequest struct {
  // A reference to a file in the form: file://pageloop.com/{container}/{application}#{url}
  Ref string `json:"ref,omitempty"`

  // Value specified as a byte slice, used when creating files from templates
  Bytes []byte

	// A source template for this file
	Template *ApplicationTemplate `json:"template,omitempty"`
}

type FileService struct {
  Host *Host
}

// Read a file.
func (s *FileService) Read(req *FileReferenceRequest, reply *ServiceReply) *StatusError {
  if req.Ref == "" {
    return CommandError(http.StatusBadRequest, "No file reference for read operation")
  }
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, _, file, err := ref.FindFile(s.Host); err != nil {
    return err
  } else {
    reply.Reply = file
  }
  return nil
}

// Read a file as a page.
func (s *FileService) ReadPage(req *FileReferenceRequest, reply *ServiceReply) *StatusError {
  if req.Ref == "" {
    return CommandError(http.StatusBadRequest, "No file reference for read page operation")
  }
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, _, file, err := ref.FindFile(s.Host); err != nil {
    return err
  } else {
    if file.Page() == nil {
      return CommandError(http.StatusNotFound, "Page %s not found", ref.Url())
    }
    reply.Reply = file.Page()
  }
  return nil
}

// Delete a file.
func (s *FileService) Delete(req *FileReferenceRequest, reply *ServiceReply) *StatusError {
  if req.Ref == "" {
    return CommandError(http.StatusBadRequest, "No file reference for delete operation")
  }
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, file, err := ref.FindFile(s.Host); err != nil {
    return err
  } else {
    if err := app.Del(file); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }
    reply.Reply = file
  }
  return nil
}

// Move a file.
func (s *FileService) Move(req *FileMoveRequest, reply *ServiceReply) *StatusError {
  if req.Ref == "" {
    return CommandError(http.StatusBadRequest, "No file reference for move operation")
  }
  if req.Destination == "" {
    return CommandError(http.StatusBadRequest, "No destination for move operation")
  }
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, file, err := ref.FindFile(s.Host); err != nil {
    return err
  } else {
    if err := app.Move(file, req.Destination); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }
    reply.Reply = file
  }
  return nil
}

// Read file content.
func (s *FileService) ReadSource(req *FileReferenceRequest, reply *ServiceReply) *StatusError {
  if req.Ref == "" {
    return CommandError(http.StatusBadRequest, "No file reference for read source operation")
  }
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, _, file, err := ref.FindFile(s.Host); err != nil {
    return err
  } else {
    reply.Reply = file.Source(false)
  }
  return nil
}

// Read raw file content (includes frontmatter).
func (s *FileService) ReadSourceRaw(req *FileReferenceRequest, reply *ServiceReply) *StatusError {
  if req.Ref == "" {
    return CommandError(http.StatusBadRequest, "No file reference for read source operation")
  }
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, _, file, err := ref.FindFile(s.Host); err != nil {
    return err
  } else {
    reply.Reply = file.Source(true)
  }
  return nil
}

// Save file content.
func (s *FileService) Save(req *FileContentRequest, reply *ServiceReply) *StatusError {
  if req.Ref == "" {
    return CommandError(http.StatusBadRequest, "No file reference for save operation")
  }
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, file, err := ref.FindFile(s.Host); err != nil {
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
func (s *FileService) Create(req *FileContentRequest, reply *ServiceReply) *StatusError {
  if req.Ref == "" {
    return CommandError(http.StatusBadRequest, "No file reference for create operation")
  }
  ref := &AssetReference{}
  ref.ParseUrl(req.Ref)
  if _, app, err := ref.FindApplication(s.Host); err != nil {
    return err
  } else {
    var exists *File = app.Urls[ref.Url()]
    if exists != nil {
      return CommandError(http.StatusConflict,"File already exists %s", ref.Url())
    }
    if app.ExistsConflict(ref.Url()) {
      return CommandError(http.StatusConflict,"File already exists, publish conflict on %s", ref.Url())
    }

    content := req.Bytes
    if req.Value != "" {
      content = []byte(req.Value)
    }

    if file, err := app.Create(ref.Url(), content); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    } else {
      reply.Reply = file
      reply.Status = http.StatusCreated
    }
  }
  return nil
}

// Create a file from a template.
func (s *FileService) CreateTemplate(req *FileTemplateRequest, reply *ServiceReply) *StatusError {
  template := req.Template
  if template == nil {
    return CommandError(http.StatusBadRequest, "No template given")
  }

  if tpl, err := s.Host.LookupTemplateFile(template); err != nil {
    return CommandError(http.StatusInternalServerError, err.Error())
  } else {
    if tpl == nil {
      return CommandError(http.StatusNotFound, "Template file %s does not exist", template.File)
    }
    creq := &FileContentRequest{Ref: req.Ref}
    creq.Bytes = tpl.Source(true)
    return s.Create(creq, reply)
  }
}
