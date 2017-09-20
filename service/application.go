package service

import(
  //"fmt"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/util"
)

type AppService struct {
  Host *Host
}

type AppReference struct {
  Container string `json:"container"`
  Application string `json:"application"`
}

/*
func (b *CommandAdapter) ReadApplication(c string, a string) (*Container, *Application, *StatusError) {
  if container, err := b.ReadContainer(c); err != nil {
    return nil, nil, err
  } else {
    app :=  container.GetByName(a)
    if app == nil {
      return nil, nil, CommandError(http.StatusNotFound, "Application %s not found", a)
    }
    return container, b.CommandExecute.ReadApplication(app), nil
  }
}
*/

func (s *AppService) Read(argv AppReference, reply *ServiceReply) *StatusError {
  // reply.Reply = app
  return nil
}
