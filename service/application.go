package service

import(
  //"fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/util"
)

type AppService struct {
  Host *Host
}

func (s *AppService) Read(argv *ApplicationReference, reply *ServiceReply) *StatusError {
  c := s.Host.GetByName(argv.Container)
  if c == nil {
    return CommandError(http.StatusNotFound, "Container %s not found", argv.Container)
  }

  app := c.GetByName(argv.Application)
  if app == nil {
    return CommandError(http.StatusNotFound, "Application %s not found", argv.Application)
  }

  reply.Reply = app
  return nil
}
