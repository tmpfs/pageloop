// Exposes a REST API to the application model.
package pageloop

import (
  "fmt"
  "os/exec"
	"regexp"
	"strings"
	"net/http"
  "mime"
  "path/filepath"
	"encoding/json"
  "github.com/tmpfs/pageloop/model"
)

const(
	API_URL = "/api/"
	JSON_MIME = "application/json; charset=utf-8"

	// App actions
	TASKS = "tasks"
	FILES = "files"
	PAGES = "pages"
  TEMPLATES = "templates"
)

var(
  HttpUtils = HttpUtil{}
	OK = []byte(`{"ok": true}`)
	SchemaAppNew = MustAsset("schema/app-new.json")
	CharsetStrip = regexp.MustCompile(`;.*$`)

  // Allowed methods.

  // TODO: CORS for OPTIONS requests
	RestAllowedMethods []string = []string{
    http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
)

// List of URLs used for bulk file operations.
type UrlList []string

// Main rest service.
type RestService struct {
  // The base mountpoint URL for the service.
	Url string
	Root *PageLoop
}

// Handles requests for application data.
type RestHandler struct {
	Root *PageLoop
	Container *model.Container
}

// Handler for asynchronous background tasks.
type TaskJobComplete struct {
  Job *Job
}

// Configure the service. Adds a rest handler for the API URL to
// the passed servemux.
func NewRestService(root *PageLoop, mux *http.ServeMux) *RestService {
  rest := &RestService{Root: root, Url: API_URL}
	mux.Handle(API_URL, http.StripPrefix(API_URL, RestHandler{Root: root}))
	return rest
}

func (t *TaskJobComplete) Done(err error, cmd *exec.Cmd, raw string) {
  // TODO: send reply to the client over websocket
  Jobs.Stop(t.Job)
  println("Task job completed: " + t.Job.Name)
  fmt.Printf("%#v\n", t.Job)
}


// Enapcaulates request information for application API endpoints.
type RequestHandler struct {
	Root *PageLoop
  // The container context for the application.
  Container *model.Container
  // Reference to the underlying application, will be nil if not found.
  App * model.Application
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
	// Check if an app exists when referenced as /api/apps/{name}
	// and extract path parts.
	if a.Path != "" {
		a.Parts = strings.Split(strings.TrimSuffix(a.Path, "/"), "/")
		a.BaseName = a.Parts[0]
		if len(a.Parts) > 1 {
		  a.Name = a.Parts[1]
    }
    if len(a.Parts) > 2 {
			a.Action = a.Parts[2]
		}
		if len(a.Parts) > 3 {
			a.Item = "/" + strings.Join(a.Parts[3:], "/")
      // Respect input trailing slash used to indicate
      // operations on a directory
      if strings.HasSuffix(a.Path, "/") {
        a.Item += "/"
      }
		}

    // Try to lookup container / application, both may be nil on 404.
    a.Container = a.Root.Host.GetByName(a.BaseName)
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
    return HttpUtils.Json(res, http.StatusOK, a.Container.Apps)
  } else {
    if a.Action == "" {
      // GET /api/{container}/{application}
      return HttpUtils.Json(res, http.StatusOK, a.App)
    } else {
      switch a.Action {
        case FILES:
          if a.Item == "" {
            // GET /api/{container}/{application}/files
            return HttpUtils.Json(res, http.StatusOK, a.App.Files)
          } else {
            // GET /api/{container}/{application}/files/{url}
            if file := a.App.GetFileByUrl(a.Item); file == nil {
              return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
            } else {
              return HttpUtils.Json(res, http.StatusOK, file)
            }
          }
        case PAGES:
          if a.Item == "" {
            // GET /api/{container}/{application}/pages
            return HttpUtils.Json(res, http.StatusOK, a.App.Pages)
          } else {
            // GET /api/{container}/{application}/pages/{url}
            if page := a.App.GetPageByUrl(a.Item); page == nil {
              return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
            } else {
              return HttpUtils.Json(res, http.StatusOK, page)
            }
          }
        default:
          return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
      }
    }
  }
  return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
}

// Handle DELETE requests.
func (a *RequestHandler) Delete(res http.ResponseWriter, req *http.Request) (int, error) {
  app := a.App

	// DELETE /api/{container}/{name} - Delete an application
  if a.Name != "" && a.Action == "" {
    if err := adapter.DeleteApplication(a.Container, a.App); err != nil {
      return HttpUtils.ErrorJson(res, err)
    }

    return HttpUtils.Ok(res, OK)
  // DELETE /api/{container}/{app}/files/ - Bulk file deletion
  } else if a.Action == FILES && a.Item == "" {
    var urls UrlList

    if content, err := HttpUtils.ReadBody(req); err != nil {
      return HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
    } else {
      if err = json.Unmarshal(content, &urls); err != nil {
        return HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
      }

      for _, url := range urls {
        if file := a.deleteFile(url, app, res, req); file == nil {
          // If we got a nil file an error occured and the response
          // will already have been sent
          return -1, nil
        }
      }
      // If we made it this far all files were deleted
      return HttpUtils.Ok(res, OK)
    }

  // DELETE /api/{container}/{app}/files/{url} - Delete a single file
  } else if a.Action == FILES && a.Item != "" {
    if file := a.deleteFile(a.Item, app, res, req); file != nil {
      return HttpUtils.Ok(res, OK)
    }
  } else {
    return HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
  }

  return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
}

func (a *RequestHandler) deleteFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) *model.File {
  var err error
  var file *model.File = app.Urls[url]
  if file == nil {
    HttpUtils.Error(res, http.StatusNotFound, nil, nil)
    return nil
  }
  if err = app.Del(file); err != nil {
    HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
    return nil
  }
  return file
}

func (a *RequestHandler) Post(res http.ResponseWriter, req *http.Request) (int, error) {
  // POST /api/{container}/{app}/files/{url}
  if a.Name != "" && a.Action == FILES && a.Item != "" {
    if file := a.PostFile(res, req); file != nil {
      return HttpUtils.Json(res, http.StatusOK, file)
    }
  }
  return HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
}

// Update the content of a file.
func (a *RequestHandler) PostFile(res http.ResponseWriter, req *http.Request) *model.File {
	var err error
  app := a.App
  url := a.Item

	loc := req.Header.Get("Location")
	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

  if loc == "" {
    // No content type header
    if ct == "" {
      HttpUtils.Error(res, http.StatusBadRequest, nil, fmt.Errorf("Content type header is required"))
      return nil
    }

    // No content length header
    if cl == "" {
      HttpUtils.Error(res, http.StatusBadRequest, nil, fmt.Errorf("Content length header is required"))
      return nil
    }
  }

	var file *model.File = app.Urls[url]
	if file != nil {
    // Handle moving the file with Location header
    if loc != "" {
      if url == loc {
        HttpUtils.ErrorJson(res,
          CommandError(http.StatusBadRequest, "Cannot move file, source and destination are equal: %s", url))
        return nil
      }

      if err := adapter.MoveFile(app, file, loc); err != nil {
        HttpUtils.ErrorJson(res, err)
        return nil
      }
    // Update file content
    } else {
      // Strip charset for mime comparison
      ct = CharsetStrip.ReplaceAllString(ct, "")
      ft := CharsetStrip.ReplaceAllString(file.Mime, "")
      if ft != ct {
        HttpUtils.Error(res, http.StatusBadRequest, nil, fmt.Errorf("Mismatched MIME types attempting to update file"))
        return nil
      }

      // TODO: fix empty reply when there is no request body
      // TODO: stream request body to disc
      var content []byte
      if content, err = HttpUtils.ReadBody(req); err == nil {
        // Update the application model
        if err = app.Update(file, content); err != nil {
          HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
          return nil
        }
      }
    }
	}

  return file
}

func (a *RequestHandler) Put(res http.ResponseWriter, req *http.Request) (int, error) {
  if a.Path != "" {
    // PUT /api/{container}/{application}/files/{url} - Create a new file.
    if a.Action == FILES && a.Item != "" {
      if file := a.PutFile(res, req); file != nil {
        return HttpUtils.Json(res, http.StatusCreated, file)
      }
    // PUT /api/{container}/{application}/tasks/ - Run a build task.
    } else if (a.Action == TASKS && a.Item != "") {
      taskName := strings.TrimPrefix(a.Item, "/")
      taskName = strings.TrimSuffix(taskName, "/")

      // No build configuration of missing build task
      if !a.App.HasBuilder() || a.App.Builder.Tasks[taskName] == "" {
        return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
      }

      fullName := fmt.Sprintf("%s/%s:%s", a.App.Container.Name, a.App.Name, taskName)

      if Jobs.GetRunningJob(fullName) != nil {
        return HttpUtils.Error(res, http.StatusConflict, nil, fmt.Errorf("Job %s is already running", fullName))
      }

      // Set up a new job for the task
      job := Jobs.NewJob(fullName)
      Jobs.Start(job)

      println("run task job: " + fullName)

      // Run the task
      a.App.Builder.Run(taskName, a.App, &TaskJobComplete{Job: job})

      // Accepted for processing
      fmt.Printf("%#v\n", job)

      // TODO: send job information to the client
      return HttpUtils.Write(res, http.StatusAccepted, OK)
    }

    return HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
  }
  return -1, nil
}


func (a *RequestHandler) PutApplication(res http.ResponseWriter, req *http.Request) (int, error) {
  var input *model.Application = &model.Application{}
  if _, err := HttpUtils.ValidateRequest(SchemaAppNew, input, req); err != nil {
    return HttpUtils.ErrorJson(res, CommandError(http.StatusBadRequest, err.Error()))
  }
  if err := adapter.CreateApplication(a.Container, input); err != nil {
    return HttpUtils.ErrorJson(res, err)
  }
  return HttpUtils.Created(res, OK)
}

// Create a new file for an application
func (a *RequestHandler) PutFile(res http.ResponseWriter, req *http.Request) *model.File {
	var err error
  app := a.App
  url := a.Item

	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

	if ct == "" {
    ct = mime.TypeByExtension(filepath.Ext(req.URL.Path))
	}

	// No content length header
	if cl == "" {
		HttpUtils.Error(res, http.StatusBadRequest, nil, fmt.Errorf("Content length header is required"))
		return nil
	}

  isDir := strings.HasSuffix(url, "/")

	var content []byte

  // Lookup template file
  if !isDir && ct == JSON_MIME {
    // TODO: stream request body to disc
    if content, err = HttpUtils.ReadBody(req); err != nil {
      HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
      return nil
    }

    input := &model.ApplicationTemplate{}
    if err = json.Unmarshal(content, input); err != nil {
      HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
      return nil
    }

    var file *model.File

    if file, err = a.Root.Host.LookupTemplateFile(input); err != nil {
      HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
      return nil
    }

    if file == nil {
      HttpUtils.Error(res, http.StatusNotFound, nil, fmt.Errorf("Template file %s does not exist", input.File))
      return nil
    }

    content = file.Source(true)
  } else {
    // TODO: stream request body to disc
    // Read in as file content upload
    if content, err = HttpUtils.ReadBody(req); err != nil {
      HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
      return nil
    }
  }

  // Update the application model
  var file *model.File
  if file, err = app.Create(url, content); err != nil {
    if err, ok := err.(model.StatusError); ok {
      HttpUtils.Error(res, err.Status, nil, err)
      return nil
    }

    HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
    return nil
  }

  return file
}

// Handle REST API endpoint requests.
func (h RestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  h.doServeHttp(res, req)
}

// Primary handler, decoupled from ServeHTTP so we can return from the function.
func (h RestHandler) doServeHttp(res http.ResponseWriter, req *http.Request) (int, error) {
	if !HttpUtils.IsMethodAllowed(req.Method, RestAllowedMethods) {
    return HttpUtils.ErrorJson(res, CommandError(http.StatusMethodNotAllowed, ""))
	}

  info := &RequestHandler{Root: h.Root}
  info.Parse(req)

  if (info.Path == "") {
    // GET / - List host containers.
    if req.Method == http.MethodGet {
      return HttpUtils.Json(res, http.StatusOK, adapter.ListContainers())
    }
  } else {
    // GET /templates - List available application templates.
    if req.Method == http.MethodGet && info.Path == TEMPLATES {
      return HttpUtils.Json(res, http.StatusOK, adapter.ListApplicationTemplates())
    }

    // METHOD /{container} - 404 if container not found.
		if info.Container == nil {
			return HttpUtils.ErrorJson(res, CommandError(http.StatusNotFound, ""))
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
    return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
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

	return HttpUtils.ErrorJson(res, CommandError(http.StatusNotFound, ""))
}
