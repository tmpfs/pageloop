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

// Request is a header written before every RPC call. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Request struct {
	ServiceMethod string   // format: "Service.Method"
	Seq           uint64   // sequence number chosen by client
	next          *Request // for free list in Server
}

// Response is a header written before every RPC return. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Response struct {
	ServiceMethod string    // echoes that of the Request
	Seq           uint64    // echoes that of the request
	Error         string    // error, if any.
	next          *Response // for free list in Server
}

type Server struct {
  serviceMap    map[string]*service
  mu          sync.RWMutex // protects the serviceMap
  reqLock     sync.Mutex // protects freeReq
  freeReq     *Request
  respLock    sync.Mutex // protects freeResp
  freeResp    *Response
}

// Register a service with the server.
func (server *Server) Register(rcvr interface{}) error {
  if server.serviceMap == nil {
    server.serviceMap = make(map[string]*service)
  }
  s := new(service)
  s.rcvr = reflect.ValueOf(rcvr)
  s.typ = reflect.TypeOf(rcvr)
  s.name = reflect.Indirect(s.rcvr).Type().Name()
  s.method = suitableMethods(s.typ)
  fmt.Printf("%s\n", s.name)
  fmt.Printf("%#v\n", s.method)
  server.serviceMap[s.name] = s
  return nil
}

// Find a method by name in dot notation (Service.Method).
func (server *Server) Method(name string) (service *service, mtype *methodType, err error) {
  dot := strings.LastIndex(name, ".")
  if dot < 0 {
    err = fmt.Errorf("rpc: service/method request ill-formed: %s", name)
    return
  }

  serviceName := name[:dot]
  methodName := name[dot+1:]

  // Look up the request.
  server.mu.RLock()
  service = server.serviceMap[serviceName]
  server.mu.RUnlock()
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
