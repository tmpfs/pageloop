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

// Handles requests for application data.
type RestHandler struct {
  Services *ServiceMap
  Adapter *CommandAdapter
	Container *Container
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

  accept := req.Header.Get("Accept")
	methodSeq := req.Header.Get("X-Method-Seq")

  // if seq, ok := methodSeq.(uint64)

  // Check sequence number
  if seq, err := strconv.ParseUint(methodSeq, 10, 64); err != nil {
    return utils.Errorj(
      res, CommandError(
        http.StatusBadRequest, "Invalid sequence number: %s", err.Error()))
  // Got a valid sequence number
  } else {
	  serviceMethod := req.Header.Get("X-Method-Name")

    println("REST Service method name: " + serviceMethod)
    println("REST Service seq: " + methodSeq)
    println("REST Service accept: " + accept)

    hasServiceMethod := h.Services.HasMethod(serviceMethod)

    // TODO: send 404 on no service method once refactor completed

    if hasServiceMethod {
      println("REST Service has service method!!")
      if req, err := h.Services.Request(serviceMethod, seq); err != nil {
        return utils.Errorj(
          res, CommandError(http.StatusInternalServerError, err.Error()))
      } else {

        // TODO: build arguments list

        if reply, err := h.Services.Call(req); err != nil {
          return utils.Errorj(
            res, CommandError(http.StatusInternalServerError, err.Error()))
        } else {
          if reply.Error != nil {
            // TODO: test for comand error response
            return utils.Errorj(
              res, CommandError(http.StatusInternalServerError, err.Error()))
          } else {
            fmt.Printf("%#v\n", reply.Reply)
            // TODO: get correct status code
            return utils.Json(res, http.StatusOK, reply.Reply)
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
