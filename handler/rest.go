// Exposes a REST API to the application
package handler

import (
  //"fmt"
  "mime"
	"net/http"
  "strings"
  "path/filepath"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
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
  Adapter *CommandAdapter
	Container *Container
}

// Configure the service. Adds a rest handler for the API URL to
// the passed servemux.
func RestService(mux *http.ServeMux, adapter *CommandAdapter) http.Handler {
  handler := RestHandler{Adapter: adapter}
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

  // Parse out an action from the request
  if act, err := h.Adapter.HttpAction(req.Method, req.URL); err != nil {
    return utils.Errorj(res, err)
  } else {
    // Attempt to match the action
    if mapping, err := h.Adapter.Find(act); err != nil {
      return utils.Errorj(res, err)
    } else {

      def := mapping.CommandDefinition

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
        // Return the result to the client
        return utils.Json(res, result.Status, result.Data)
      }
    }
    return utils.Errorj(res, CommandError(http.StatusNotFound, ""))
  }
  return utils.Errorj(res, CommandError(http.StatusNotFound, ""))
}
