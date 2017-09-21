// Exposes a REST API to the application
package handler

import (
  "fmt"
  "mime"
	"net/http"
  "strconv"
  "strings"
  "path/filepath"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/rpc"
  . "github.com/tmpfs/pageloop/service"
  . "github.com/tmpfs/pageloop/util"
)

var(
  utils = HttpUtil{}
  // Allowed methods.
	RestAllowedMethods []string = []string{
    http.MethodGet,
    http.MethodPost,
    http.MethodPut,
    http.MethodDelete,
    http.MethodOptions}
	SchemaAppNew = MustAsset("schema/app-new.json")
)

// Type that extracts values from an HTTP request
// and builds the arguments that should be passed
// to the corresponding rpc service method.
type HttpArguments struct {
  req *http.Request
  *Parameters
}

// Determine the arguments to use for an rpc method call.
//
// The name argument should be the qualified Service.Method name.
func (args *HttpArguments) Get(name string, req *http.Request) (argv interface{}, err *StatusError) {

  switch name {
    case "Container.CreateApp":
      c := &Container{Name: args.Parameters.Context}
      var app *Application = &Application{Container: c}
      if _, err := utils.ValidateRequest(SchemaAppNew, app, req); err != nil {
        return nil, CommandError(http.StatusBadRequest, err.Error())
      }
      argv = app
    case "Container.Read":
      argv = &Container{Name: args.Parameters.Context}
    case "Application.ReadFiles":
      fallthrough
    case "Application.ReadPages":
      fallthrough
    case "Application.Delete":
      fallthrough
    case "Application.Read":
      argv = &Application{
        Name: args.Parameters.Target,
        ContainerName: args.Parameters.Context}
    case "Application.DeleteFiles":
      var list UrlList = make(UrlList, 0)
      if err := utils.ReadJson(req, &list); err != nil {
        return nil, CommandError(http.StatusInternalServerError, err.Error())
      }
      argv = &Application{
        Name: args.Parameters.Target,
        ContainerName: args.Parameters.Context,
        Batch: &list}
    case "Application.RunTask":
      argv = &Application{
        Name: args.Parameters.Target,
        ContainerName: args.Parameters.Context,
        Task: args.Parameters.Item}
    case "File.ReadSource":
      fallthrough
    case "File.ReadSourceRaw":
      fallthrough
    case "File.Move":
      c := &Container{Name: args.Parameters.Context}
      a := &Application{Name: args.Parameters.Target, Container: c}
      argv = &File{Owner: a, Url: args.Parameters.Item, Destination: req.Header.Get("Location")}
    case "File.Save":
      c := &Container{Name: args.Parameters.Context}
      a := &Application{Name: args.Parameters.Target, Container: c}
      f := &File{Owner: a, Url: args.Parameters.Item}
      if content, err := utils.ReadBody(req); err != nil {
        return nil, CommandError(http.StatusInternalServerError, err.Error())
      } else {
        f.Bytes(content)
      }
      argv = f
  }
  return argv, nil
}

func NewHttpArguments(req *http.Request) *HttpArguments {
  args := &HttpArguments{req: req, Parameters: &Parameters{}}
  args.Parameters.Parse(req.URL.Path)
  return args
}

type Parameters struct {
  // input path
  Path string
  // Slice of parameter parts
  Parts []string
  // The operation type, cannot be a wildcard.
  Type string `json:"type"`
  // Context for the operation. May be a container reference, job number etc.
  Context string `json:"context"`
  // Target for the operation, typically an application.
  Target string `json:"target"`
  // A filter operation for the request.
  Filter string `json:"filter"`
  // An item, may contain slashes.
  Item string `json:"item"`
}

func (act *Parameters) Parse(path string) {
  act.Path = path
  if act.Path != "" {
    path := strings.TrimPrefix(act.Path, SLASH)
    path = strings.TrimSuffix(path, SLASH)
    act.Parts = strings.Split(path, SLASH)
    act.Type = act.Parts[0]
    if len(act.Parts) > 1 {
      act.Context = act.Parts[1]
    }
    if len(act.Parts) > 2 {
      act.Target = act.Parts[2]
    }
    if len(act.Parts) > 3 {
      act.Filter = act.Parts[3]
    }
    if len(act.Parts) > 4 {
      act.Item = SLASH + strings.Join(act.Parts[4:], SLASH)
      // Respect input trailing slash used to indicate
      // operations on a directory
      if strings.HasSuffix(act.Path, SLASH) {
        act.Item += SLASH
      }
    }
  }

  // So that trailing slash with no URL will match
  // the filter
  if act.Item == SLASH {
    act.Item = ""
  }
}

// Handles requests for application data.
type RestHandler struct {
  Services *ServiceMap

  // Deprecated
  Adapter *CommandAdapter
}

// Configure the service. Adds a rest handler for the API URL to
// the passed servemux.
func RestService(mux *http.ServeMux, adapter *CommandAdapter, services *ServiceMap) http.Handler {
  handler := RestHandler{Adapter: adapter, Services: services}
  mux.Handle(API_URL, http.StripPrefix(API_URL, handler))
	return handler
}

// Handle REST API endpoint requests.
func (h RestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  h.doServeHttp(res, req)
}

// Primary handler, decoupled from ServeHTTP so we can return from the function.
func (h RestHandler) doServeHttp(res http.ResponseWriter, req *http.Request) (int, error) {
	if !utils.IsMethodAllowed(req.Method, RestAllowedMethods) {
    return utils.Errorj(res, CommandError(http.StatusMethodNotAllowed, ""))
	}

	ct := req.Header.Get("Content-Type")

	if ct == "" {
    ct = mime.TypeByExtension(filepath.Ext(req.URL.Path))
	}

	methodSeq := req.Header.Get("X-Method-Seq")

  // Keep uploads working using old api
  if methodSeq == "" {
    return h.doDeprecatedHttp(res, req)
  }

  // TODO: automatically find matching route and validate request is well formed

  // Check sequence number
  if seq, err := strconv.ParseUint(methodSeq, 10, 64); err != nil {
    return utils.Errorj(
      res, CommandError(
        http.StatusBadRequest, "Invalid sequence number: %s", err.Error()))
  // Got a valid sequence number
  } else {
	  serviceMethod := req.Header.Get("X-Method-Name")

    println("REST method name: " + serviceMethod)
    //println("REST Service seq: " + methodSeq)
    //println("REST Service accept: " + accept)

    hasServiceMethod := h.Services.HasMethod(serviceMethod)

    // TODO: send 404 on no service method once refactor completed
    if hasServiceMethod {
      println("REST service method name (start invoke): " + serviceMethod)
      if rpcreq, err := h.Services.Request(serviceMethod, seq); err != nil {
        return utils.Errorj(
          res, CommandError(http.StatusInternalServerError, err.Error()))
      } else {

        // Build rpc arguments
        args := NewHttpArguments(req)
        if argv, err := args.Get(serviceMethod, req); err != nil {
          return utils.Errorj(res, err)
        } else {

          fmt.Printf("%#v\n", args.Parameters)
          // Got some arguments to use for the request
          if argv != nil {
            rpcreq.Argv(argv)
          }

          if reply, err := h.Services.Call(rpcreq); err != nil {
            return utils.Errorj(
              res, CommandError(http.StatusInternalServerError, err.Error()))
          } else {
            if reply.Error != nil {
              // Send status error if we can
              if err, ok := reply.Error.(*StatusError); ok {
                return utils.Errorj(res, err)
              // Otherwise handle as plain error
              } else {
                return utils.Errorj(
                  res, CommandError(http.StatusInternalServerError, err.Error()))
              }
            } else {
              var replyData = reply.Reply
              status := http.StatusOK

              if result, ok := replyData.(*ServiceReply); ok {
                replyData = result.Reply
                if result.Status != 0 {
                  status = result.Status
                }
              }

              //fmt.Printf("%#v\n", reply)
              //fmt.Printf("status: %d\n", status)

              if serviceMethod == "Container.CreateApp" {
                // Mount the application, needs to be done here due to some funky
                // package cyclic references
                if app, ok := replyData.(*Application); ok {
                  MountApplication(h.Adapter.Mountpoints.MountpointMap, h.Adapter.Host, app)
                }
              }

              accept := req.Header.Get("Accept")
              // Client is asking for binary response
              if accept == "application/octet-stream" {
                // If the method result is a slice of bytes send it back
                if content, ok := replyData.([]byte); ok {
                  return utils.Write(res, status, content)
                }
              }
              return utils.Json(res, status, replyData)
            }
          }
        }
      }
    }
  }

  return h.doDeprecatedHttp(res, req)
}


func (h RestHandler) doDeprecatedHttp(res http.ResponseWriter, req *http.Request) (int, error) {
	if !utils.IsMethodAllowed(req.Method, RestAllowedMethods) {
    return utils.Errorj(res, CommandError(http.StatusMethodNotAllowed, ""))
	}

	ct := req.Header.Get("Content-Type")

	if ct == "" {
    ct = mime.TypeByExtension(filepath.Ext(req.URL.Path))
	}

  // Parse out an action from the request
  if act, err := h.Adapter.HttpAction(req.Method, req.URL); err != nil {
    return utils.Errorj(res, err)
  } else {

    // Attempt to match the action
    if mapping, err := h.Adapter.Find(act); err != nil {
      return utils.Errorj(res, err)
    } else {

      def := mapping.CommandDefinition

      // println(def.MethodName)

      // TODO: use proper RPC arguments interface
      if def.MethodName == "CreateApp" {
        var input *Application = &Application{}
        if _, err := utils.ValidateRequest(SchemaAppNew, input, req); err != nil {
          return utils.Errorj(res, CommandError(http.StatusBadRequest, err.Error()))
        }
        act.Push(input)
      } else if def.MethodName == "DeleteFiles" {
        var input UrlList = make(UrlList, 0)
        if err := utils.ReadJson(req, &input); err != nil {
          return utils.Errorj(res, err)
        }
        act.Push(input)
      } else if def.MethodName == "CreateFile" {
          isDir := strings.HasSuffix(act.Item, SLASH)
          // Create from a template
          if !isDir && ct == JSON_MIME {
            ref := &ApplicationTemplate{}
            if err := utils.ReadJson(req, ref); err != nil {
              return utils.Errorj(res, err)
            }
            act = h.Adapter.Mutate(act, "CreateFileTemplate")
            act.Push(ref)
          // File content bytes creation
          } else {
            if content, err := utils.ReadBody(req); err != nil {
              return utils.Errorj(res, CommandError(http.StatusInternalServerError, err.Error()))
            } else {
              act.Push(content)
            }
          }
      } else if def.MethodName == "UpdateFile" {
          location := req.Header.Get("Location")
          if location != "" {
            act = h.Adapter.Mutate(act, "MoveFile")
            act.Push(location)
          } else {
            if content, err := utils.ReadBody(req); err != nil {
              return utils.Errorj(res, CommandError(http.StatusInternalServerError, err.Error()))
            } else {
              act.Push(content)
            }
          }
        }
      // Invoke the command
      if result, err := h.Adapter.Execute(act); err != nil {
        return utils.Errorj(res, err)
      } else {
        if def.MethodName == "CreateApp" {
          // Mount the application, needs to be done here due to some funky
          // package cyclic references
          if app, ok := result.Data.(*Application); ok {
            MountApplication(h.Adapter.Mountpoints.MountpointMap, h.Adapter.Host, app)
          }
        }

        accept := req.Header.Get("Accept")
        // Client is asking for binary response
        if accept == "application/octet-stream" {
          // If the method result is a slice of bytes send it back
          if content, ok := result.Data.([]byte); ok {
            return utils.Write(res, result.Status, content)
          }
        }

        // Return the result to the client
        return utils.Json(res, result.Status, result.Data)
      }
    }
    return utils.Errorj(res, CommandError(http.StatusNotFound, ""))
  }
  return utils.Errorj(res, CommandError(http.StatusNotFound, ""))
}
