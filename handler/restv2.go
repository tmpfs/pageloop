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
)

// Handles requests for application data.
type RestV2Handler struct {
  Adapter *CommandAdapter
	Container *Container
}

// Configure the service. Adds a rest handler for the API URL to
// the passed servemux.
func RestServiceV2(mux *http.ServeMux, adapter *CommandAdapter) http.Handler {
  handler := RestV2Handler{Adapter: adapter}
  mux.Handle(API_URL + "v2/", http.StripPrefix(API_URL + "v2/", handler))
	return handler
}

// Handle REST API endpoint requests.
func (h RestV2Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  h.doServeHttp(res, req)
}

// Primary handler, decoupled from ServeHTTP so we can return from the function.
func (h RestV2Handler) doServeHttp(res http.ResponseWriter, req *http.Request) (int, error) {
	if !utils.IsMethodAllowed(req.Method, RestAllowedMethods) {
    return utils.Errorj(res, CommandError(http.StatusMethodNotAllowed, ""))
	}

  // Parse out an action from the request
  if act, err := h.Adapter.CommandAction(req.Method, req.URL); err != nil {
    return utils.Errorj(res, err)
  } else {
    // Got an action - execute it
    if result, err := h.Adapter.Execute(act); err != nil {
      return utils.Errorj(res, err)
    } else {
      return utils.Json(res, result.Status, result.Data)
    }
    return utils.Errorj(res, CommandError(http.StatusNotFound, ""))
  }
  return utils.Errorj(res, CommandError(http.StatusNotFound, ""))
}
