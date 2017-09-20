// Package service provides a transport agnostic rpc service manager.
//
// It borrows from the net/rpc package for the reflection code but does
// not provide any networking code.
//
// Useful when you want to expose the same service API methods over different
// transports. Likely this can be done with the net/rpc package I just couldn't
// figure out how to make it work nicely.
//
// The problem is that we want to expose a REST API and JSON-RPC over websockets.
// However if we use the JSON-RPC package then we have to declare an initial
// *http.Request argument but for websockets we don't have an HTTP request.
//
// Also in the case of JSON-RPC over websockets deserializing to an RPC call is
// simple but in the case of the REST API we need to dynamically build the arguments
// list from the HTTP request information.
//
// This modified rpc implementation allows us to provide services in a consistent
// fashion (no initial *http.Request argument) and allow for arguments to be collated
// in the case of the REST API.
//
// Services must declare methods in exactly the same way as for net/rpc with the
// single distinction that the return value does not need to be error it needs to
// conform to the error interface (provide an Error() function) which allows our
// service methods to return custom error implementations.
//
// This package does not allow setting custom service names they are always inferred
// from the receiver name.
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
  service       *service
  methodType    *methodType
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

// Register a service and panic on error.
func (server *Server) MustRegister(rcvr interface{}) {
  if err := server.Register(rcvr); err != nil {
    panic(err)
  }
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
  if method, err := suitableMethods(s.typ); err != nil {
    return err
  } else {
    s.method = method
    server.serviceMap[s.name] = s
    return nil
  }
}

// Find a method by name in dot notation (Service.Method).
func (server *Server) method(name string) (service *service, mtype *methodType, err error) {
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

// Get a method call request.
func (server *Server) Method(name string, seq uint64) (req *Request, err error) {
  var service *service
  var mtype *methodType

  if service, mtype, err = server.method(name); err != nil {
    return
  }

  req = &Request{ServiceMethod: name, Seq: seq, service: service, methodType: mtype}
  return
}
