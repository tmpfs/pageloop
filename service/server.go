package service

import(
  "fmt"
  "reflect"
  "sync"
  "strings"
)

type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
	ReplyType  reflect.Type
	numCalls   uint
}

type service struct {
  name string
  rcvr reflect.Value
  typ reflect.Type
  method map[string]*methodType // registered methods
}

type ServiceManager struct {
  services map[string]*service
}

// Register a service with the server.
func (server *ServiceManager) Register(rcvr interface{}) error {
  if server.services == nil {
    server.services = make(map[string]*service)
  }
  s := new(service)
  s.rcvr = reflect.ValueOf(rcvr)
  s.typ = reflect.TypeOf(rcvr)
  s.name = reflect.Indirect(s.rcvr).Type().Name()
  s.method = suitableMethods(s.typ, true)
  fmt.Printf("%s\n", s.name)
  fmt.Printf("%#v\n", s.method)
  server.services[s.name] = s
  return nil
}

// Find a method by name in dot notation (Service.Method).
func (server *ServiceManager) Method(name string) (service *service, mtype *methodType, err error) {
  dot := strings.LastIndex(name, ".")
  if dot < 0 {
    err = fmt.Errorf("rpc: service/method request ill-formed: %s", name)
    return
  }

  serviceName := name[:dot]
  methodName := name[dot+1:]

  // Look up the request.
  //server.mu.RLock()
  service = server.services[serviceName]
  //server.mu.RUnlock()
  if service == nil {
    err = fmt.Errorf("rpc: can't find service %s", serviceName)
    return
  }
  mtype = service.method[methodName]
  if mtype == nil {
    err = fmt.Errorf("rpc: can't find method %s", methodName)
  }
  return
}
