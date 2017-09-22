// Package rpc provides a transport agnostic rpc service manager.
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
package rpc

import(
  "fmt"
  "reflect"
  "sync"
  "strings"
)

type methodType struct {
	sync.Mutex  // protects counters
	method      reflect.Method
	numCalls    uint
	ArgType     reflect.Type
	ReplyType   reflect.Type
}

type service struct {
  name string
  rcvr reflect.Value
  typ reflect.Type
  method map[string]*methodType // registered methods
}

// Represents the arguments that may be passed to a method invocation.
//
// You can manually set the Argv value by calling Argv() on a Request.
type RequestArguments struct {
  Argv reflect.Value
  Replyv reflect.Value
}

// Request is a header written before every RPC call. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Request struct {
	ServiceMethod string        // format: "Service.Method"
	Seq           uint64        // sequence number chosen by client
  Arguments     *RequestArguments
  service       *service
  methodType    *methodType
	next          *Request      // for free list in ServiceMap
}

// Response is a header written before every RPC return. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Response struct {
	ServiceMethod string        // echoes that of the Request
	Seq           uint64        // echoes that of the request
	Error         error         // error, if any
  Reply         interface{}   // method invocation reply argument
	next          *Response     // for free list in ServiceMap
}

type ServiceMap struct {
  serviceMap  map[string]*service
  mu          sync.RWMutex    // protects the serviceMap
  reqLock     sync.Mutex      // protects freeReq
  freeReq     *Request
  respLock    sync.Mutex      // protects freeResp
  freeResp    *Response
}

type ServiceInfo struct {
  Name string `json:"name"`
  Methods []*ServiceMethodInfo `json:"methods"`
}

type ServiceMethodInfo struct {
  // Fully qualified method name
  ServiceMethod string `json:"method"`
  // Method name
  Name string `json:"name"`
  // Number of calls for this method
  Calls uint `json:"calls"`
  // Type for the argument (first argument)
  ArgType string `json:"arg"`
  // Type for the reply value (second argument)
  ReplyType string `json:"reply"`
  // Placeholder for user data associated with the service (route information)
  UserData interface{} `json:"info,omitempty"`
}

// Set argv for a service method call request.
//
// This is for the case when you need to build the arguments
// manually from various input sources such as an HTTP request
// where input can come from URL parameters, query string, body data etc.
//
// You need to be sure you pass the correct type here otherwise you will
// get a runtime panic when you attempt to call the underlying method.
func (req *Request) Argv(args interface{}) {
  req.Arguments.Argv = reflect.ValueOf(args)
}

// Register a service and panic on error.
func (server *ServiceMap) MustRegister(rcvr interface{}, name string) {
  if err := server.Register(rcvr, name); err != nil {
    panic(err)
  }
}

// Get a map of public service information.
func (server *ServiceMap) Map() map[string]*ServiceInfo {
  m := make(map[string]*ServiceInfo)
  for key, srv := range server.serviceMap {
    info := &ServiceInfo{Name: key}

    for _, mt := range srv.method {
      mi := &ServiceMethodInfo{
        ServiceMethod: key + "." + mt.method.Name,
        Calls: mt.numCalls,
        ArgType: mt.ArgType.String(),
        ReplyType: mt.ReplyType.String(),
        Name: mt.method.Name}
      info.Methods = append(info.Methods, mi)
    }
    key = strings.ToLower(key)
    m[key] = info
  }
  return m
}

// Register a service with the server.
//
// If name is the empty string the service name is inferred from
// the type name of rcvr.
func (server *ServiceMap) Register(rcvr interface{}, name string) error {
  if server.serviceMap == nil {
    server.serviceMap = make(map[string]*service)
  }
  s := new(service)
  s.rcvr = reflect.ValueOf(rcvr)
  s.typ = reflect.TypeOf(rcvr)
  if name == "" {
    name = reflect.Indirect(s.rcvr).Type().Name()
  }
  s.name = name
  if method, err := suitableMethods(s.typ); err != nil {
    return err
  } else {
    s.method = method
    server.serviceMap[s.name] = s
  }
  return nil
}

// Get a method call request.
func (server *ServiceMap) Request(name string, seq uint64) (req *Request, err error) {
  var service *service
  var mtype *methodType

  if service, mtype, err = server.method(name); err != nil {
    return
  }

  args := server.arguments(mtype)
  req = &Request{
    ServiceMethod: name,
    Seq: seq,
    Arguments: args,
    service: service,
    methodType: mtype}
  return
}

// Call the method using the given request.
//
// Returns a Response propagated with the Reply argument from the
// method invocation and an Error if the method returned an error.
func (server *ServiceMap) Call(req *Request) (res *Response, err error) {
  var reply interface{}
  res = &Response{ServiceMethod: req.ServiceMethod, Seq: req.Seq}
  reply, err = server.call(req)
  if err != nil {
    res.Error = err
  } else {
    res.Reply = reply
  }
  return
}

// Determine if a service method exists by name using dot notation (Service.Method).
func (server *ServiceMap) HasMethod(name string) bool {
  if _, _, err := server.method(name); err != nil {
    return false
  }
  return true
}

// Private

// Call the service method represented by the request.
//
// Returns the method call reply argument and an error if set.
func (server *ServiceMap) call(req *Request) (reply interface{}, err error) {
  mtype := req.methodType
	mtype.Lock()
  mtype.numCalls++
  mtype.Unlock()
  function := mtype.method.Func

  // Invoke the method, providing a new value for the reply.
  returnValues := function.Call([]reflect.Value{req.service.rcvr, req.Arguments.Argv, req.Arguments.Replyv})

  // Set up the reply to return
  reply = req.Arguments.Replyv.Interface()

  if !returnValues[0].IsNil() {
    // The return value for the method should be an error implementation.
    errResponse := returnValues[0].Interface()
    if errResponse != nil {
      err = errResponse.(error)
    }
  }

  return
}

// Initialize the method arguments.
func (server *ServiceMap) arguments (mtype *methodType) *RequestArguments {
  var argv, replyv reflect.Value
  // Decode the argument value.
  if mtype.ArgType.Kind() == reflect.Ptr {
    argv = reflect.New(mtype.ArgType.Elem())
  } else {
    argv = reflect.New(mtype.ArgType).Elem()
  }
  // argv guaranteed to be a pointer now.
  replyv = reflect.New(mtype.ReplyType.Elem())
  return &RequestArguments{Argv: argv, Replyv: replyv}
}

// Split a fully qualified service name, error if there is no dot in the name.
func (server *ServiceMap) split(name string) (service string, method string, err error) {
  dot := strings.LastIndex(name, ".")
  if dot < 0 {
    err = fmt.Errorf("rpc: service/method request ill-formed: %s", name)
    return
  }
  service = name[:dot]
  method = name[dot+1:]
  return
}
// Find a method by name in dot notation (Service.Method).
func (server *ServiceMap) method(name string) (service *service, mtype *methodType, err error) {
  var serviceName string
  var methodName string
  if serviceName, methodName, err = server.split(name); err != nil {
    return
  }

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
