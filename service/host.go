package service

import(
  //"fmt"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type HostService struct {
  Host *Host
}

// List containers.
func (s *HostService) List(argv *VoidArgs, reply *ServiceReply) *StatusError {
  reply.Reply = s.Host.Containers
  return nil
}
