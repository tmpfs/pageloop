package service

import(
  //"fmt"
  "net/http"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type ContainerService struct {
  Host *Host
}

// Read a container.
func (s *ContainerService) Read(container string, reply *ServiceReply) *StatusError {
  c := s.Host.GetByName(container)
  if c == nil {
    return CommandError(http.StatusNotFound, "Container %s not found", container)
  }
  reply.Reply = c
  return nil
}
