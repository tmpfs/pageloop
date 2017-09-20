package service

import(
  //"fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/util"
)

type AppService struct {
  *ContainerService
  Host *Host
}

// Read an application.
func (s *AppService) Read(app *Application, reply *ServiceReply) *StatusError {
  if app, err := s.lookup(app); err != nil {
    return err
  } else {
    reply.Reply = app
  }
  return nil
}

// Read the files for an application.
func (s *AppService) ReadFiles(app *Application, reply *ServiceReply) *StatusError {
  if app, err := s.lookup(app); err != nil {
    return err
  } else {
    reply.Reply = app.Files
  }
  return nil
}

// Read the pages for an application.
func (s *AppService) ReadPages(app *Application, reply *ServiceReply) *StatusError {
  if app, err := s.lookup(app); err != nil {
    return err
  } else {
    reply.Reply = app.Pages
  }
  return nil
}


// Private

func (s *AppService) lookup(app *Application) (*Application, *StatusError) {
  c := s.Host.GetByName(app.ContainerName)
  if c == nil {
    return nil, CommandError(http.StatusNotFound, "Container %s not found", app.ContainerName)
  }
  app = c.GetByName(app.Name)
  if app == nil {
    return nil, CommandError(http.StatusNotFound, "Application %s not found", app.Name)
  }
  return app, nil
}
