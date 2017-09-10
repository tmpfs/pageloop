// Exposes a REST API to the application
package handler

import (
	"regexp"
	"strings"
	"net/http"
  "mime"
  "path/filepath"
	"encoding/json"
  . "github.com/tmpfs/pageloop/adapter"
  . "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

const(
	// App actions
	TASKS = "tasks"
	FILES = "files"
	PAGES = "pages"
  TEMPLATES = "templates"
)

// TODO: remove all calls to utils.Error() - prefer ErrorJson()

var(
  utils = HttpUtil{}
	SchemaAppNew = MustAsset("schema/app-new.json")
	CharsetStrip = regexp.MustCompile(`;.*$`)

  // TODO: CORS for OPTIONS requests
  // Allowed methods.
	RestAllowedMethods []string = []string{
    http.MethodGet,
    http.MethodPost,
    http.MethodPut,
    http.MethodDelete,
    http.MethodOptions}
)

// List of URLs used for bulk file operations.
type UrlList []string

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

// Enapcaulates request information for application API endpoints.
type RequestHandler struct {
  Adapter *CommandAdapter
  // The container context for the application.
  Container *Container
  // Reference to the underlying application, will be nil if not found.
  App * Application
  // Request URL path
  Path string
  // Path parts
  Parts []string
  // Container name or top-level action.
  BaseName string
  // Application name.
  Name string
  // Action identifier: files|pages|tasks etc.
  Action string
  // Item to operate on, the remaining part of the path, eg: /docs/index.html
  Item string
}

// Parse the request path and assign fields to the request handler.
func (a *RequestHandler) Parse(req *http.Request) {
	a.Path = req.URL.Path
	if a.Path != "" {
		a.Parts = strings.Split(strings.TrimSuffix(a.Path, SLASH), SLASH)
		a.BaseName = a.Parts[0]
		if len(a.Parts) > 1 {
		  a.Name = a.Parts[1]
    }
    if len(a.Parts) > 2 {
			a.Action = a.Parts[2]
		}
		if len(a.Parts) > 3 {
			a.Item = SLASH + strings.Join(a.Parts[3:], SLASH)
      // Respect input trailing slash used to indicate
      // operations on a directory
      if strings.HasSuffix(a.Path, SLASH) {
        a.Item += SLASH
      }
		}

    // Try to lookup container / application, both may be nil on 404.
    a.Container = a.Adapter.Host.GetByName(a.BaseName)
    if a.Container != nil {
		  a.App = a.Container.GetByName(a.Name)
    }
	}
}

// Handle GET requests.
func (a *RequestHandler) Get(res http.ResponseWriter, req *http.Request) (int, error) {
  if a.Path == "" {
    // TODO: check this is necessary
    // GET /api/
    return utils.Json(res, http.StatusOK, a.Container.Apps)
  } else {
    if a.Action == "" {
      // GET /api/{container}/{application}
      return utils.Json(res, http.StatusOK, a.App)
    } else {
      switch a.Action {
        case FILES:
          if a.Item == "" {
            // GET /api/{container}/{application}/files
            return utils.Json(res, http.StatusOK, a.App.Files)
          } else {
            // GET /api/{container}/{application}/files/{url}
            if file := a.App.GetFileByUrl(a.Item); file == nil {
              return utils.Error(res, http.StatusNotFound, nil, nil)
            } else {
              return utils.Json(res, http.StatusOK, file)
            }
          }
        case PAGES:
          if a.Item == "" {
            // GET /api/{container}/{application}/pages
            return utils.Json(res, http.StatusOK, a.App.Pages)
          } else {
            // GET /api/{container}/{application}/pages/{url}
            if page := a.App.GetPageByUrl(a.Item); page == nil {
              return utils.Error(res, http.StatusNotFound, nil, nil)
            } else {
              return utils.Json(res, http.StatusOK, page)
            }
          }
        default:
          return utils.Error(res, http.StatusNotFound, nil, nil)
      }
    }
  }
  return utils.Error(res, http.StatusNotFound, nil, nil)
}

// Handle DELETE requests.
func (a *RequestHandler) Delete(res http.ResponseWriter, req *http.Request) (int, error) {
  var files []*File

	// DELETE /api/{container}/{name} - Delete an application
  if a.Name != "" && a.Action == "" {
    if app, err := a.Adapter.DeleteApplication(a.Container, a.App); err != nil {
      return utils.ErrorJson(res, err)
    } else {
      return utils.Json(res, http.StatusOK, app)
    }
  // DELETE /api/{container}/{app}/files/ - Bulk file deletion
  } else if a.Action == FILES && a.Item == "" {
    var urls UrlList
    if content, err := utils.ReadBody(req); err != nil {
      return utils.Error(res, http.StatusInternalServerError, nil, err)
    } else {
      if err = json.Unmarshal(content, &urls); err != nil {
        return utils.Error(res, http.StatusInternalServerError, nil, err)
      }

      for _, url := range urls {
        if file, err := a.Adapter.DeleteFile(a.App, url); err != nil {
          return utils.ErrorJson(res, err)
        } else {
          files = append(files, file)
        }
      }
      // If we made it this far all files were deleted
      return utils.Json(res, http.StatusOK, files)
    }
  // DELETE /api/{container}/{app}/files/{url} - Delete a single file
  } else if a.Action == FILES && a.Item != "" {
    if file, err := a.Adapter.DeleteFile(a.App, a.Item); err != nil {
      return utils.ErrorJson(res, err)
    } else {
      files = append(files, file)
      return utils.Json(res, http.StatusOK, files)
    }
  }
  return utils.ErrorJson(res, CommandError(http.StatusNotFound, ""))
}

func (a *RequestHandler) Post(res http.ResponseWriter, req *http.Request) (int, error) {
  // POST /api/{container}/{app}/files/{url}
  if a.Name != "" && a.Action == FILES && a.Item != "" {
    if file, err := a.PostFile(res, req); err != nil {
      return utils.ErrorJson(res, err)
    } else {
      return utils.Json(res, http.StatusOK, file)
    }
  }
  return utils.Error(res, http.StatusNotFound, nil, nil)
}

// Update the content of a file.
func (a *RequestHandler) PostFile(res http.ResponseWriter, req *http.Request) (*File, *StatusError) {
	var err error
  app := a.App
  url := a.Item

	loc := req.Header.Get("Location")
	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

  if loc == "" {
    // No content type header
    if ct == "" {
      return nil, CommandError(http.StatusBadRequest, "Content type header is required")
    }

    // No content length header
    if cl == "" {
      return nil, CommandError(http.StatusBadRequest, "Content length header is required")
    }
  }

	var file *File = app.Urls[url]
	if file != nil {
    // Handle moving the file with Location header
    if loc != "" {
      if url == loc {
        return nil, CommandError(
          http.StatusBadRequest, "Cannot move file, source and destination are equal: %s", url)
      }

      if err := a.Adapter.MoveFile(app, file, loc); err != nil {
        return nil, err
      }
    // Update file content
    } else {
      // Strip charset for mime comparison
      ct = CharsetStrip.ReplaceAllString(ct, "")
      ft := CharsetStrip.ReplaceAllString(file.Mime, "")
      if ft != ct {
        return nil, CommandError(
          http.StatusBadRequest, "Mismatched MIME types attempting to update file")
      }

      // TODO: fix empty reply when there is no request body
      // TODO: stream request body to disc
      var content []byte
      if content, err = utils.ReadBody(req); err == nil {
        // Update the application model
        if _, err := a.Adapter.UpdateFile(app, file, content); err != nil {
          return nil, err
        }
      }
    }
	}

  return file, nil
}

func (a *RequestHandler) Put(res http.ResponseWriter, req *http.Request) (int, error) {
  if a.Path != "" {
    // PUT /api/{container}/{application}/files/{url} - Create a new file.
    if a.Action == FILES && a.Item != "" {
      if file, err := a.PutFile(res, req); err != nil {
        return utils.ErrorJson(res, err)
      } else {
        return utils.Json(res, http.StatusCreated, file)
      }
    // PUT /api/{container}/{application}/tasks/ - Run a build task.
    } else if (a.Action == TASKS && a.Item != "") {
      taskName := strings.TrimPrefix(a.Item, SLASH)
      taskName = strings.TrimSuffix(taskName, SLASH)

      if job, err := a.Adapter.RunTask(a.App, taskName); err != nil {
        return utils.ErrorJson(res, err)
      } else {
        // TODO: send job information to the client - should a Job+Task pair
        return utils.Json(res, http.StatusAccepted, job)
      }
    }

    return utils.Error(res, http.StatusMethodNotAllowed, nil, nil)
  }
  return utils.ErrorJson(res, CommandError(http.StatusNotFound, ""))
}

func (a *RequestHandler) PutApplication(res http.ResponseWriter, req *http.Request) (int, error) {
  var input *Application = &Application{}
  if _, err := utils.ValidateRequest(SchemaAppNew, input, req); err != nil {
    return utils.ErrorJson(res, CommandError(http.StatusBadRequest, err.Error()))
  }
  if app, err := a.Adapter.CreateApplication(a.Container, input); err != nil {
    return utils.ErrorJson(res, err)
  } else {
    // Mount the application
    MountApplication(a.Adapter.Mountpoints.MountpointMap, a.Adapter.Host, app)
    return utils.Json(res, http.StatusCreated, app)
  }

  return utils.ErrorJson(res, CommandError(http.StatusNotFound, ""))
}

// Create a new file for an application
func (a *RequestHandler) PutFile(res http.ResponseWriter, req *http.Request) (*File, *StatusError) {
	var err error
	var content []byte

	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

	if ct == "" {
    ct = mime.TypeByExtension(filepath.Ext(req.URL.Path))
	}

	// No content length header
	if cl == "" {
		return nil, CommandError(http.StatusBadRequest, "Content length header is required")
	}

  isDir := strings.HasSuffix(a.Item, SLASH)

  // TODO: stream request body to disc
  // Read in as file content upload
  if content, err = utils.ReadBody(req); err != nil {
    return nil, CommandError(http.StatusInternalServerError, err.Error())
  }

  // Create from a template
  if !isDir && ct == JSON_MIME {
    tplref := &ApplicationTemplate{}
    if err = json.Unmarshal(content, tplref); err != nil {
      return nil, CommandError(http.StatusInternalServerError, err.Error())
    }
    return a.Adapter.CreateFileTemplate(a.App, a.Item, tplref)
  }

  // Create from request body content
  return a.Adapter.CreateFile(a.App, a.Item, content)
}

// Handle REST API endpoint requests.
func (h RestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  h.doServeHttp(res, req)
}

// Primary handler, decoupled from ServeHTTP so we can return from the function.
func (h RestHandler) doServeHttp(res http.ResponseWriter, req *http.Request) (int, error) {
	if !utils.IsMethodAllowed(req.Method, RestAllowedMethods) {
    return utils.ErrorJson(res, CommandError(http.StatusMethodNotAllowed, ""))
	}

  info := &RequestHandler{Adapter: h.Adapter}
  info.Parse(req)

  if (info.Path == "") {
    // GET / - List host containers.
    if req.Method == http.MethodGet {
      return utils.Json(res, http.StatusOK, h.Adapter.ListContainers())
    }
  } else {
    // GET /templates - List available application templates.
    if req.Method == http.MethodGet && info.Path == TEMPLATES {
      return utils.Json(res, http.StatusOK, h.Adapter.ListApplicationTemplates())
    }

    // METHOD /{container} - 404 if container not found.
		if info.Container == nil {
			return utils.ErrorJson(res, CommandError(http.StatusNotFound, ""))
		}
	}

  // Container level endpoints
	switch req.Method {
		case http.MethodPut:
		  // PUT /api/{container}/ - Create a new application.
			if info.Container != nil && info.Name == "" {
        return info.PutApplication(res, req)
			}
  }

  // Application must exist
  if info.App == nil {
    return utils.Error(res, http.StatusNotFound, nil, nil)
  }

  // Application level endpoints
	switch req.Method {
		case http.MethodGet:
      return info.Get(res, req)
		case http.MethodDelete:
      return info.Delete(res, req)
		case http.MethodPut:
      return info.Put(res, req)
		case http.MethodPost:
      return info.Post(res, req)
	}

  return utils.ErrorJson(res, CommandError(http.StatusNotFound, ""))
}
