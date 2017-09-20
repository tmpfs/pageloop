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
func (s *HostService) List(argv VoidArgs, reply *Host) *StatusError {
  reply.Containers = s.Host.Containers
  return nil
}
