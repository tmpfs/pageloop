package service

import(
  //"fmt"
)

// Type for service methods that do not accept any arguments.
type VoidArgs struct {}

type Core struct {
  Name string
  Version string
}

type MetaReply struct {
  Name string `json:"name"`
  Version string `json:"version"`
}

func (s *Core) Meta(args VoidArgs, reply *MetaReply) error {
  reply.Name = s.Name
  reply.Version = s.Version
  return nil
}
