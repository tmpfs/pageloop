package handler

import(
  //"fmt"
  "net/http"
  "strings"
  . "github.com/tmpfs/pageloop/util"
)

var(
  router *Router
)

// Parameters represents the parsed URL path parameters.
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

type Route struct {
  // Name of a service method to invoke.
  ServiceMethod string
  // Route definition path or request path
  Path string
  // Request method
  Method string
  // Status code to send when ok
  Status int
  // Parsed path parameters
  *Parameters
  // Condition used to match route
  Condition func(req *http.Request) bool
}

type Router struct {
  get []*Route
  put []*Route
  post []*Route
  del []*Route
  methods map[string]*Route
}

// Adds a route to the router.
func (r *Router) Add(route *Route) *Route {
  route.Parameters = &Parameters{}
  route.Parameters.Parse(route.Path)
  if r.methods == nil {
    r.methods = make(map[string]*Route)
  }
  r.methods[route.ServiceMethod] = route
  switch(route.Method) {
    case http.MethodGet:
      r.get = append(r.get, route)
    case http.MethodPut:
      r.put = append(r.put, route)
    case http.MethodPost:
      r.post = append(r.post, route)
    case http.MethodDelete:
      r.del = append(r.del, route)
  }
  return route
}

func (r *Router) Find(req *http.Request) *Route {
  var res *Route
  list := r.list(req.Method)
  res = r.matchPath(req.URL.Path, list)
  // TODO: shortcut lookup on X-Method-Seq / X-Method-Name
  return res
}

func (r *Router) matchPath(path string, list []*Route) *Route {
  params := &Parameters{}
  params.Parse(path)

  matches := func(mapped *Route, params *Parameters) bool {
    // Root match
    if mapped.Path == "" && (path == "" || path == "/") {
      return true
    }

    // Must be same length
    if len(mapped.Parameters.Parts) != len(params.Parts) {
      return false
    }

    var i int
    var l int = len(mapped.Parameters.Parts)
    var p string

    for i, p = range mapped.Parameters.Parts {
      if p == "*" {
        continue
      } else if (p != params.Parts[i]) {
        return false
      }
    }

    //fmt.Printf("i is: %d\n", i)
    //fmt.Printf("l is: %d\n", l)

    // Everything matched
    if i == (l - 1) {
      return true
    }

    return false
  }

  for _, mapped := range list {
    if matches(mapped, params) {
      r := &Route{
        ServiceMethod: mapped.ServiceMethod,
        Path: path,
        Parameters: params,
        Method: mapped.Method,
        Status: mapped.Status}
      return r
    }
  }

  return nil
}

func (r *Router) list(method string) []*Route {
  switch(method) {
    case http.MethodGet:
      return r.get
    case http.MethodPut:
      return r.put
    case http.MethodPost:
      return r.post
    case http.MethodDelete:
      return r.del
  }
  return nil
}

func init() {
  router = &Router{}

  add := func(route *Route) *Route {
    // router.routes = append(router.routes, route)
    return router.Add(route)
  }

  route := func(service string, path string, method string, status int) *Route {
    r := &Route{ServiceMethod: service, Path: path, Method: method, Status: status}
    return add(r)
  }

  var r *Route

  route("Core.Meta", "", http.MethodGet, http.StatusOK)
  route("Core.Stats", "/stats", http.MethodGet, http.StatusOK)
  route("Template.ReadApplications", "/templates", http.MethodGet, http.StatusOK)
  route("Job.ActiveJob", "/jobs", http.MethodGet, http.StatusOK)
  route("Job.Read", "/jobs/*", http.MethodGet, http.StatusOK)
  route("Job.Delete", "/jobs/*", http.MethodDelete, http.StatusOK)
  route("Host.List", "/apps", http.MethodGet, http.StatusOK)
  route("Container.Read", "/apps/*", http.MethodGet, http.StatusOK)
  route("Container.CreateApp", "/apps/*", http.MethodPut, http.StatusCreated)
  route("Application.Read", "/apps/*/*", http.MethodGet, http.StatusOK)
  route("Application.Delete", "/apps/*/*", http.MethodDelete, http.StatusOK)
  route("Application.ReadFiles", "/apps/*/*/files", http.MethodGet, http.StatusOK)
  route("Application.ReadPages", "/apps/*/*/pages", http.MethodGet, http.StatusOK)
  route("Application.DeleteFiles", "/apps/*/*/files", http.MethodDelete, http.StatusOK)
  route("Application.RunTask", "/apps/*/*/tasks/", http.MethodPut, http.StatusAccepted)
  route("File.Read", "/apps/*/*/files/*", http.MethodGet, http.StatusOK)
  route("File.ReadPage", "/apps/*/*/pages/*", http.MethodGet, http.StatusOK)
  route("File.ReadSource", "/apps/*/*/src/*", http.MethodGet, http.StatusOK)
  route("File.ReadSourceRaw", "/apps/*/*/raw/*", http.MethodGet, http.StatusOK)
  route("File.Create", "/apps/*/*/files/*", http.MethodPut, http.StatusCreated)
  route("File.Save", "/apps/*/*/files/*", http.MethodPost, http.StatusOK)
  route("File.Delete", "/apps/*/*/files/*", http.MethodDelete, http.StatusOK)

  // TODO: conditional on location header
  r = route("File.Move", "/apps/*/*/files/*", http.MethodPost, http.StatusOK)
  r.Condition = func(req *http.Request) bool {
    return req.Header.Get("Location") != ""
  }

  // TODO: conditional on template object
  route("File.CreateTemplate", "/apps/*/*/files/*", http.MethodPut, http.StatusCreated)
}
