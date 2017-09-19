package adapter

import(
  "fmt"
  "log"
  "reflect"
  "sync"
  "strings"
  "unicode"
  "unicode/utf8"
)

var rpc *ServiceManager = &ServiceManager{}
// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

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

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			if reportErr {
				log.Println("method", mname, "has wrong number of ins:", mtype.NumIn())
			}
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				log.Println(mname, "argument type not exported:", argType)
			}
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			if reportErr {
				log.Println("method", mname, "reply type not a pointer:", replyType)
			}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				log.Println("method", mname, "reply type not exported:", replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if reportErr {
				log.Println("method", mname, "has wrong number of outs:", mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			if reportErr {
				log.Println("method", mname, "returns", returnType.String(), "not error")
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}

// Private

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func init () {
  rpc.Register(new(Core))

  s, m, _ := rpc.Method("Core.Meta")
  fmt.Printf("%#v\n", s)
  fmt.Printf("%#v\n", m)
}
