// Package handler provides HTTP and network handlers.
package handler

import (
  "fmt"
  //"mime"
  "strings"
  "strconv"
	"net/http"
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

// Get the service method arguments for the given request and matched route.
func Argv(route *Route, req *http.Request, res http.ResponseWriter) (argv interface{}, err *StatusError) {
  name := route.ServiceMethod
  switch name {
    case "Container.CreateApp":
      var app *ApplicationRequest =
        &ApplicationRequest{Container: route.Parameters.Context}

      // TODO: move validation to service logic
      if _, err := utils.ValidateRequest(SchemaAppNew, app, req); err != nil {
        return nil, CommandError(http.StatusBadRequest, err.Error())
      }
      argv = app
    case "Container.Read":
      argv = &Container{Name: route.Parameters.Context}
    case "Archive.Import":
      fallthrough
    case "Archive.Export":
      app := &ApplicationRequest{
        Name: route.Parameters.Target,
        Container: route.Parameters.Context}

      name := app.Name

      // Switch archive type
      archiveType := ArchiveFull
      filter := strings.TrimSuffix(route.Parameters.Item, SLASH)
      filter = strings.TrimPrefix(filter, SLASH)
      if filter == "source" {
        archiveType = ArchiveSource
        name += "-" + filter
      } else if filter == "public" {
        archiveType = ArchivePublic
        name += "-" + filter
      }

      name += ".zip"
      res.Header().Set("Content-Disposition", "attachment; filename=" + name)

      argv = &ArchiveRequest{Application: app, Writer: res, Type: archiveType, Name: name}
    case "Application.ReadFiles":
      fallthrough
    case "Application.ReadPages":
      fallthrough
    case "Application.Delete":
      fallthrough
    case "Application.Read":
      argv = &ApplicationRequest{
        Name: route.Parameters.Target,
        Container: route.Parameters.Context}
    case "Application.DeleteFiles":
      var list UrlList = make(UrlList, 0)
      if err := utils.ReadJson(req, &list); err != nil {
        return nil, CommandError(http.StatusInternalServerError, err.Error())
      }
      argv = &ApplicationRequest{
        Name: route.Parameters.Target,
        Container: route.Parameters.Context,
        Batch: &list}
    case "Application.RunTask":
      argv = &ApplicationRequest{
        Name: route.Parameters.Target,
        Container: route.Parameters.Context,
        Task: route.Parameters.Item}
    case "File.ReadSource":
      fallthrough
    case "File.ReadSourceRaw":
      fallthrough
    case "File.Read":
      /*
      c := &Container{Name: route.Parameters.Context}
      a := &Application{Name: route.Parameters.Target, Container: c}
      f := &File{Owner: a, Url: route.Parameters.Item}
      */
      // argv = f
      ref := fmt.Sprintf(
        "file://%s/%s#%s",
        route.Parameters.Context,
        route.Parameters.Target,
        route.Parameters.Item)
      argv = &FileRequest{Ref: ref}
    case "File.Move":
      /*
      c := &Container{Name: route.Parameters.Context}
      a := &Application{Name: route.Parameters.Target, Container: c}
      argv = &File{Owner: a, Url: route.Parameters.Item, Destination: req.Header.Get("Location")}
      */
      ref := fmt.Sprintf(
        "file://%s/%s#%s",
        route.Parameters.Context,
        route.Parameters.Target,
        route.Parameters.Item)
      argv = &FileRequest{Ref: ref, Destination: req.Header.Get("Location")}
    case "File.CreateTemplate":
      /*
      c := &Container{Name: route.Parameters.Context}
      a := &Application{Name: route.Parameters.Target, Container: c}
      f := &File{Owner: a, Url: route.Parameters.Item}
      f.Template = &ApplicationTemplate{}
      if err := utils.ReadJson(req, f.Template); err != nil {
        return nil, err
      }
      argv = f
      */
      ref := fmt.Sprintf(
        "file://%s/%s#%s",
        route.Parameters.Context,
        route.Parameters.Target,
        route.Parameters.Item)
      f := &FileRequest{Ref: ref, Template: &ApplicationTemplate{}}
      if err := utils.ReadJson(req, f.Template); err != nil {
        return nil, err
      }
      argv = f
    case "File.Create":
      fallthrough
    case "File.Save":
      /*
      c := &Container{Name: route.Parameters.Context}
      a := &Application{Name: route.Parameters.Target, Container: c}
      f := &File{Owner: a, Url: route.Parameters.Item}
      if content, err := utils.ReadBody(req); err != nil {
        return nil, CommandError(http.StatusInternalServerError, err.Error())
      } else {
        f.Bytes(content)
      }
      argv = f
      */
      ref := fmt.Sprintf(
        "file://%s/%s#%s",
        route.Parameters.Context,
        route.Parameters.Target,
        route.Parameters.Item)
      f := &FileRequest{Ref: ref}
      if content, err := utils.ReadBody(req); err != nil {
        return nil, CommandError(http.StatusInternalServerError, err.Error())
      } else {
        f.Bytes = content
      }
      argv = f
    case "Service.Read":
      argv = &ServiceRequest{Service: route.Parameters.Context}
    case "Service.ReadMethodCalls":
      fallthrough
    case "Service.ReadMethod":
      argv = &ServiceMethodRequest{
        Service: route.Parameters.Context,
        Method: route.Parameters.Target}
  }
  return argv, nil
}


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

  if route, err := DefaultRouter.Find(req); err != nil {
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
        if argv, err := Argv(route, req, res); err != nil {
          return utils.Errorj(res, err)
        } else {

          // Got some arguments to use for the request
          if argv != nil {
            rpcreq.Argv(argv)
          }

          // Call the service function
          Stats.Rpc.Add("calls", 1)
          if reply, err := h.Services.Call(rpcreq); err != nil {
            Stats.Rpc.Add("errors", 1)
            // Send status error if we can
            if err, ok := reply.Error.(*StatusError); ok {
              return utils.Errorj(res, err)
            // Otherwise handle as plain error
            } else {
              return utils.Errorj(
                res, CommandError(http.StatusInternalServerError, err.Error()))
            }
          } else {
            // NOTE: we don't need to test reply.Error as the error is always returned

            // Success send the response to the client
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
            if route.ResponseType == ResponseTypeNone {
              // Service method wrote the response body
              return 0, nil
            } else if route.ResponseType == ResponseTypeByte {
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

  // Default response is not found
  return utils.Errorj(
    res, CommandError(http.StatusNotFound, ""))
}
