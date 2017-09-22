// Package handler provides HTTP and network handlers.
package handler

import (
  // "fmt"
  //"mime"
  "strconv"
	"net/http"
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

// Handles requests for application data.
type RestHandler struct {
  Services *ServiceMap
  Host *Host
  Mountpoints *MountpointManager
}

// Configure the service. Adds a rest handler for the API URL to
// the passed servemux.
func RestService(mux *http.ServeMux, services *ServiceMap, host *Host, mountpoints *MountpointManager) http.Handler {
  handler := RestHandler{Services: services, Host: host, Mountpoints: mountpoints}
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

  // Never cache API requests
  res.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")

  if route, err := router.Find(req); err != nil {
    return utils.Errorj(res, err)
  } else {
    // No matching route
    if route == nil {
      return utils.Errorj(
        res, CommandError(http.StatusNotFound, "No route matched for path %s", req.URL.Path))
    }

    // fmt.Printf("route: %#v\n", route)
    hasServiceMethod := h.Services.HasMethod(route.ServiceMethod)

    // Check if the service method is available
    if !hasServiceMethod {
      return utils.Errorj(
        res, CommandError(http.StatusNotFound,
        "No service available for method name %s using path %s", route.ServiceMethod, req.URL.Path))
    } else {
      // println("REST service method name (start invoke): " + route.ServiceMethod)

      // Get a service method call request
      if rpcreq, err := h.Services.Request(route.ServiceMethod, route.Seq); err != nil {
        return utils.Errorj(
          res, CommandError(http.StatusInternalServerError, err.Error()))
      } else {
        // Get rpc arguments
        if argv, err := router.Argv(route, req); err != nil {
          return utils.Errorj(res, err)
        } else {

          // Got some arguments to use for the request
          if argv != nil {
            rpcreq.Argv(argv)
          }

          // Call the service function
          if reply, err := h.Services.Call(rpcreq); err != nil {
            return utils.Errorj(
              res, CommandError(http.StatusInternalServerError, err.Error()))
          } else {
            // Reply with error when available
            if reply.Error != nil {
              // Send status error if we can
              if err, ok := reply.Error.(*StatusError); ok {
                return utils.Errorj(res, err)
              // Otherwise handle as plain error
              } else {
                return utils.Errorj(
                  res, CommandError(http.StatusInternalServerError, err.Error()))
              }
            // Success send the response to the client
            } else {
              var replyData = reply.Reply

              // TODO: use route status and remove from ServiceReply
              status := http.StatusOK

              if result, ok := replyData.(*ServiceReply); ok {
                replyData = result.Reply
                if result.Status != 0 {
                  status = result.Status
                }
              }

              // NOTE: After functions need some thought!
              if route.ServiceMethod == "Container.CreateApp" {
                // Mount the application, needs to be done here due to some funky
                // package cyclic references
                if app, ok := replyData.(*Application); ok {
                  MountApplication(h.Mountpoints.MountpointMap, h.Host, app)
                }
              }

              // Indicate to the client the response type.

              // Allows the client to determine whether a response should
              // be parsed as JSON or not.
              res.Header().Set("X-Response-Type", strconv.Itoa(route.ResponseType))

              // Determine how we should reply to the client
              if route.ResponseType == ResponseTypeByte {
                // TODO: work out correct MIME type from file???

                // If the method result is a slice of bytes send it back
                if content, ok := replyData.([]byte); ok {
                  return utils.Write(res, status, content)
                } else {
                  return utils.Errorj(
                    res, CommandError(
                      http.StatusInternalServerError,
                      "Service method failed to return []byte for binary response type"))
                }
              }

              // Assume JSON response if response type not already handled
              return utils.Json(res, status, replyData)
            }
          }
        }
      }
    }
  }

  // Default response is not found
  return utils.Errorj(
    res, CommandError(http.StatusNotFound, ""))
}
