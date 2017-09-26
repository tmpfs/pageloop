package service

import(
  //"fmt"
  "net/http"
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

// Get a service.
func (s *RpcServices) Get(name string, reply *ServiceReply) *StatusError {
  m := s.Services.Map()
  if srv, err := LookupService(m, name); err != nil {
    return err
  } else {
    reply.Reply = srv
  }
  return nil
}

func LookupService(mapping map[string]*ServiceInfo, name string) (*ServiceInfo, *StatusError) {
  if mapping[name] == nil {
    return nil, CommandError(http.StatusNotFound, "Service %s not found", name)
  }
  return mapping[name], nil
}
