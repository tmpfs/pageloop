package service

import(
  //"fmt"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/util"
)

type RpcServices struct {
  Services *ServiceMap
  Router *Router
}

// List services.
func (s *RpcServices) List(argv *VoidArgs, reply *ServiceReply) *StatusError {
  m := s.Services.Map()

  // Inject route information
  for _, srv := range m {
    for _, method := range srv.Methods {
      r := s.Router.Get(method.ServiceMethod)
      if r != nil {
        method.UserData = r
      }
    }
  }

  reply.Reply = m
  return nil
}
