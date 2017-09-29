package service

import(
  //"fmt"
  "strings"
  "net/http"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/util"
)

var(
  ServicesMetaInfo map[string]*ServiceMeta
)

type ServiceMeta struct {
  Description string `json:"description"`
}

type ServiceRequest struct {
  Service string `json:"service"`
}

type ServiceMethodRequest struct {
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
    InjectServiceMeta(srv, s.Router)
  }

  reply.Reply = m
  return nil
}

// Get a service.
func (s *RpcServices) Read(req *ServiceRequest, reply *ServiceReply) *StatusError {
  m := s.Services.Map()
  if srv, err := LookupService(m, req.Service); err != nil {
    return err
  } else {
    InjectServiceMeta(srv, s.Router)
    reply.Reply = srv
  }
  return nil
}

// Get a service method.
func (s *RpcServices) ReadMethod(req *ServiceMethodRequest, reply *ServiceReply) *StatusError {
  m := s.Services.Map()
  if method, err := LookupServiceMethod(m, req.Service, req.Method); err != nil {
    return err
  } else {
    InjectMethodMeta(method, s.Router)
    reply.Reply = method
  }
  return nil
}

// Get the number of times a service method has been called.
func (s *RpcServices) ReadMethodCalls(req *ServiceMethodRequest, reply *ServiceReply) *StatusError {
  m := s.Services.Map()
  if method, err := LookupServiceMethod(m, req.Service, req.Method); err != nil {
    return err
  } else {
    reply.Reply = method.Calls
  }
  return nil
}

// Inject service meta and route information.
func InjectServiceMeta(srv *ServiceInfo, router *Router) {
  for _, method := range srv.Methods {
    InjectMethodMeta(method, router)
  }
}

func InjectMethodMeta(method *ServiceMethodInfo, router *Router) {
  method.UserMeta = ServicesMetaInfo[method.ServiceMethod]
  // TODO: handle multiple matching routes with different verbs etc
  r := router.GetAll(method.ServiceMethod)
  if r != nil {
    method.UserData = r
  }
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

func init () {

  ServicesMetaInfo = make(map[string]*ServiceMeta)

  describe := func (name string, desc string) {
    ServicesMetaInfo[name] = &ServiceMeta{Description: desc}
  }

  describe("Core.Meta", `Get server meta information.`)
  describe("Core.Stats", `Get server statistics.`)
  describe("Service.List", `List available services.`)
  describe("Service.Read", `Get service information.`)
  describe("Service.ReadMethod", `Get service method information.`)
  describe("Service.ReadMethodCalls", `Get the number of calls for a service method.`)
  describe("Template.List", `List application templates.`)
  describe("Job.List", `Get active jobs.`)
  describe("Job.Read", `Get an active job.`)
  describe("Job.Delete", `Delete an active job.`)
  describe("Host.List", `List application containers.`)
  describe("Container.Read", `Get container information.`)
  describe("Container.CreateApp", `Create a new application.`)
  describe("Application.Read", `Get an application.`)
  describe("Application.Delete", `Delete an application.`)
  describe("Application.ReadFiles", `Get the files list for an application.`)
  describe("Application.ReadPages", `Get the pages list for an application.`)
  describe("Application.DeleteFiles", `Delete files from an application.`)
  describe("Application.RunTask", `Run an application build task.`)
  describe("File.Read", `Get file information.`)
  describe("File.ReadPage", `Get page information.`)
  describe("File.Create", `Create a new file.`)
  describe("File.Save", `Save file content.`)
  describe("File.Delete", `Delete a file.`)
  describe("File.ReadSource", `Get the contents of a file.`)
  describe("File.ReadSourceRaw", `Get the raw contents of a file.`)
  describe("File.Move", `Move a file.`)
  describe("File.CreateTemplate", `Create a file from a template.`)
  describe("Archive.Export", `Export a zip archive.`)
}
