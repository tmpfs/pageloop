package rpc

import(
  "fmt"
  "reflect"
  "unicode"
  "unicode/utf8"
)

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfString = reflect.TypeOf((*string)(nil)).Elem()

// suitableMethods returns suitable Rpc methods of typ, it will return error
// on the first error condition encountered.
func suitableMethods(typ reflect.Type) (map[string]*methodType, error) {
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
			return nil, fmt.Errorf("method %s has wrong number of ins: %d", mname, mtype.NumIn())
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			return nil, fmt.Errorf("%s argument type not exported: %v", mname, argType)
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("method %s reply type not a pointer: %v", mname, replyType)
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			return nil, fmt.Errorf("method %s reply type not exported: %v", mname, replyType)
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			return nil, fmt.Errorf("method %s has wrong number of outs: %d", mname, mtype.NumOut())
		}

    returnType := mtype.Out(0)
    if returnType != typeOfError {
		  // The return type of the method must implement the error interface if it is not a plain error.

      // Check for Error() function
      if returnErrorMethod, ok := returnType.MethodByName("Error"); !ok {
        return nil, fmt.Errorf(
          "method %s return value %s does not have an Error() method", mname, returnType.String())
      } else {
        rmtype := returnErrorMethod.Type
        // Error() should just have a receiver argument
        if rmtype.NumIn() != 1 {
          return nil, fmt.Errorf(
            "method %s return type %s Error() method has wrong number of ins: %d",
            mname, returnType.String(), rmtype.NumIn())
        }

        // Error() method needs one out (string)
        if rmtype.NumOut() != 1 {
          return nil, fmt.Errorf(
            "method %s return type %s Error() method has wrong number of outs: %d",
            mname, returnType.String(), rmtype.NumOut())
        }


        // Check that Error() return type is string
        errorReturnType := rmtype.Out(0)
        if errorReturnType != typeOfString {
          return nil, fmt.Errorf(
            "method %s return type %s Error() method does not return string", mname, returnType.String())
        }
      }
    }
    methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}

	return methods, nil
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
