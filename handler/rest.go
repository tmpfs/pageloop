// Exposes a REST API to the application
package handler

import (
	"net/http"
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

  // Parse out an action from the request
  if act, err := h.Adapter.HttpAction(req.Method, req.URL); err != nil {
    return utils.Errorj(res, err)
  } else {
    // Attempt to match the action
    if mapping, err := h.Adapter.Find(act); err != nil {
      return utils.Errorj(res, err)
    } else {

      if mapping.CommandDefinition.MethodName == "CreateApp" {
        var input *Application = &Application{}
        if _, err := utils.ValidateRequest(SchemaAppNew, input, req); err != nil {
          return utils.Errorj(res, CommandError(http.StatusBadRequest, err.Error()))
        }
        act.Push(input)
      } else if mapping.CommandDefinition.MethodName == "DeleteFiles" {
        println("read url list for batch delete")
      }

      // Invoke the command
      if result, err := h.Adapter.Execute(act); err != nil {
        // Route does not match
        return utils.Errorj(res, err)
      } else {
        // Return the result to the client
        return utils.Json(res, result.Status, result.Data)
      }
    }
    return utils.Errorj(res, CommandError(http.StatusNotFound, ""))
  }
  return utils.Errorj(res, CommandError(http.StatusNotFound, ""))
}
