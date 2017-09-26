package service

import(
  //"fmt"
  "strings"
  "net/http"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/util"
)

type ServiceLookupRequest struct {
  Service string `json:"service"`
  Method string `json:"method"`
}

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
func (s *RpcServices) Read(req *ServiceLookupRequest, reply *ServiceReply) *StatusError {
  m := s.Services.Map()
  if srv, err := LookupService(m, req.Service); err != nil {
    return err
  } else {
    reply.Reply = srv
  }
  return nil
}

// Get a service method.
func (s *RpcServices) ReadMethod(req *ServiceLookupRequest, reply *ServiceReply) *StatusError {
  m := s.Services.Map()
  if method, err := LookupServiceMethod(m, req.Service, req.Method); err != nil {
    return err
  } else {
    reply.Reply = method
  }
  return nil
}

// Lookup a service.
func LookupService(mapping map[string]*ServiceInfo, service string) (*ServiceInfo, *StatusError) {
  if mapping[service] == nil {
    return nil, CommandError(http.StatusNotFound, "Service %s not found", service)
  }
  return mapping[service], nil
}

// Lookup a service and method.
func LookupServiceMethod(mapping map[string]*ServiceInfo, service string, method string) (*ServiceMethodInfo, *StatusError) {
  if srv, err := LookupService(mapping, service); err != nil {
    return nil, err
  } else {
    name := strings.ToLower(method)
    for _, info := range srv.Methods {
      if name == strings.ToLower(info.Name) {
        return info, nil
      }
    }
  }
  return nil, CommandError(http.StatusNotFound, "Method %s not found for %s service", method, service)
}
