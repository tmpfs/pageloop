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
	RestAllowedMethods []string = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
)

type TaskJobComplete struct {
  Job *Job
}

func (t *TaskJobComplete) Done(err error, cmd *exec.Cmd, raw string) {
  // TODO: send reply to the client over websocket
  Jobs.Stop(t.Job)
  println("Task job completed: " + t.Job.Name)
  fmt.Printf("%#v\n", t.Job)
}

type UrlList []string

type RestService struct {
	Url string
	Root *PageLoop
}

func NewRestService(root *PageLoop, mux *http.ServeMux) *RestService {
  rest := &RestService{Root: root, Url: API_URL}
  url := API_URL
	mux.Handle(url, http.StripPrefix(url, RestHandler{Root: root}))
	return rest
}

// Handles requests for application data.
type RestHandler struct {
	Root *PageLoop
	Container *model.Container
}

func (h RestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  h.doServeHttp(res, req)
}

// Enapcaulates request information for application API endpoints.
type ApplicationRequestHandler struct {
	Root *PageLoop
  // The container context for the application.
  Container *model.Container
  // Reference to the underlying application, will be nil if not found.
  App * model.Application
  // Request URL path
  Path string
  // Application name.
  Name string
  // Action identifier: files|pages|tasks etc.
  Action string
  // Item to operate on, the remaining part of the path, eg: /docs/index.html
  Item string
}

// Parse the request path and assign fields to the request handler.
func (a *ApplicationRequestHandler) Parse(req *http.Request) {
	a.Path = req.URL.Path
	// Check if an app exists when referenced as /api/apps/{name}
	// and extract path parts.
	if a.Path != "" {
		parts := strings.Split(strings.TrimSuffix(a.Path, "/"), "/")
		a.Name = parts[0]
		if len(parts) > 1 {
			a.Action = parts[1]
		}
		if len(parts) > 2 {
			//item = parts[2]
			a.Item = "/" + strings.Join(parts[2:], "/")
      // Respect input trailing slash used to indicate
      // operations on a directory
      if strings.HasSuffix(a.Path, "/") {
        a.Item += "/"
      }
		}
		a.App = a.Container.GetByName(a.Name)
	}
}

// Handle GET requests.
func (a *ApplicationRequestHandler) Get(res http.ResponseWriter, req *http.Request) (int, error) {
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
func (a *ApplicationRequestHandler) Delete(res http.ResponseWriter, req *http.Request) (int, error) {
  app := a.App

	// DELETE /api/{container}/{name} - Delete an application
  if a.Name != "" && a.Action == "" {
    if app.Protected {
      return HttpUtils.Error(res, http.StatusForbidden, nil, fmt.Errorf("Cannot delete protected application"))
    }

    // Stop serving files for the application
    a.Root.UnmountApplication(app)

    // Delete the mountpoint
    if err := a.Root.DeleteApplicationMountpoint(app); err != nil {
      return HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
    }

    // Delete the files
    if err := a.Root.DeleteApplicationFiles(app); err != nil {
      return HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
    }

    // Delete the in-memory application
    a.Container.Del(app)

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

func (a *ApplicationRequestHandler) deleteFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) *model.File {
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

func (a *ApplicationRequestHandler) Post(res http.ResponseWriter, req *http.Request) (int, error) {
  // POST /api/{container}/{app}/files/{url}
  if a.Name != "" && a.Action == FILES && a.Item != "" {
    if file := a.postFile(a.Item, a.App, res, req); file != nil {
      return HttpUtils.Json(res, http.StatusOK, file)
    }
  }
  return HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
}

// Update the content of a file.
func (a *ApplicationRequestHandler) postFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) *model.File {
	var err error
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
        HttpUtils.Error(res, http.StatusBadRequest, nil,
          fmt.Errorf("Cannot move file, source and destination are equal: %s", url))
        return nil
      }

      if err = app.Move(file, loc); err != nil {
        HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
        return nil
      }
      return file
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

func (a *ApplicationRequestHandler) Put(res http.ResponseWriter, req *http.Request) (int, error) {
  return -1, nil
}


// Handles application information (files, pages etc.)
func (h RestHandler) doServeHttp(res http.ResponseWriter, req *http.Request) (int, error) {
	var err error
	var data []byte
  var file *model.File

	if !HttpUtils.IsMethodAllowed(req.Method, RestAllowedMethods) {
		return HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
	}

	url := req.URL
	path := url.Path

  // TODO: only allow this in Dev mode?
  res.Header().Set("Access-Control-Allow-Origin", "*")

	// List host containers
	if path == "" {
		if req.Method != http.MethodGet {
			return HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
		}
    return HttpUtils.Json(res, http.StatusOK, h.Root.Host.Containers)
  // List available application templates
	} else if path == TEMPLATES {
    apps := adapter.ListApplicationTemplates()
    return HttpUtils.Json(res, http.StatusOK, apps)
  }

	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) > 0 {
		var c *model.Container = h.Root.Host.GetByName(parts[0])
		// Container not found
		if c == nil {
			return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
		}

    h.Container = c

		// Proxy to the app handler
		// Using http.StripPrefix() here does not invoke
		// the underlying handler???
		// handler := RestHandler{Root: h.Root, Container: c}
		req.URL.Path = strings.TrimPrefix(req.URL.Path, parts[0])
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/")
	}

  info := &ApplicationRequestHandler{Root: h.Root, Container: h.Container}
  info.Parse(req)

  path = info.Path
  app := info.App
  name := info.Name
  action := info.Action
  item := info.Item

  // Container level endpoints
	switch req.Method {
		case http.MethodPut:
		  // PUT /api/{container}/
			if path == "" {
				var input *model.Application = &model.Application{}
				_, err = HttpUtils.ValidateRequest(SchemaAppNew, input, req)
				if err != nil {
					return HttpUtils.Error(res, http.StatusBadRequest, nil, err)
				}

        existing := h.Container.GetByName(input.Name)
        if existing != nil {
					return HttpUtils.Error(res, http.StatusPreconditionFailed, nil, fmt.Errorf("Application %s already exists", input.Name))
        }

        input.Url = input.MountpointUrl(h.Container)

        // mountpoint exists
        exists := h.Root.HasMountpoint(input.Url)
        if exists {
					return HttpUtils.Error(res, http.StatusPreconditionFailed, nil, fmt.Errorf("Mountpoint URL %s already exists", input.Url))
        }

        var mountpoint *Mountpoint

        // Create and save a mountpoint for the application.
        if mountpoint, err = h.Root.CreateMountpoint(input); err != nil {
					return HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
        }

        if input.Template != nil {
          var source *model.Application

          // Find the template app/ directory
          if source, err = h.Root.LookupTemplate(input.Template); err != nil {
            return HttpUtils.Error(res, http.StatusBadRequest, nil, err);
          }

          // Copy template source files
          if err = h.Root.CopyApplicationTemplate(input, source); err != nil {
            return HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
          }
        }

        var app *model.Application

        // Load and publish the app source files
        if app, err = h.Root.LoadMountpoint(*mountpoint, h.Container); err != nil {
          return HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
        }

        // Mount the application
        h.Root.MountApplication(app)

				return HttpUtils.Created(res, OK)
			}
  }

  // Application must exist
  if info.App == nil {
    return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
  }

	switch req.Method {
		case http.MethodGet:
      return info.Get(res, req)
		case http.MethodDelete:
      return info.Delete(res, req)
		case http.MethodPut:
      if path != "" {
				// PUT /api/{container}/{app}/files/{url}
				if name != "" && action == FILES && item != "" {
					if file = h.putFile(item, app, res, req); file != nil {
            return HttpUtils.Json(res, http.StatusCreated, file)
          }
				} else if (name != "" && action == TASKS && item != "") {
          taskName := strings.TrimPrefix(item, "/")
          taskName = strings.TrimSuffix(taskName, "/")

          // No build configuration of missing build task
          if !app.HasBuilder() || app.Builder.Tasks[taskName] == "" {
            return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
          }

          fullName := fmt.Sprintf("%s/%s:%s", app.Container.Name, app.Name, taskName)

          if Jobs.GetRunningJob(fullName) != nil {
            return HttpUtils.Error(res, http.StatusConflict, nil, fmt.Errorf("Job %s is already running", fullName))
          }

          // Set up a new job for the task
          job := Jobs.NewJob(fullName)
          Jobs.Start(job)

          println("run task job: " + fullName)

          // Run the task
          app.Builder.Run(taskName, app, &TaskJobComplete{Job: job})

          // Accepted for processing

          fmt.Printf("%#v\n", job)

          // TODO: send job information to the client
          return HttpUtils.Write(res, http.StatusAccepted, OK)
        }

				return HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
			}
		case http.MethodPost:
      return info.Post(res, req)
	}

	if err != nil {
		return HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
	}

	if data != nil {
		return HttpUtils.Ok(res, data)
  }

	return HttpUtils.Error(res, http.StatusNotFound, nil, nil)
}

// Create a new file for an application
func (h RestHandler) putFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) *model.File {
	var err error

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

    if file, err = h.Root.LookupTemplateFile(input); err != nil {
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
