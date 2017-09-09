// Exposes a REST API to the application model.

package pageloop

import (
  "fmt"
  "os/exec"
	"errors"
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
	var rest *RestService = &RestService{Root: root, Url: API_URL}

	var url string
	url= API_URL
	mux.Handle(url, http.StripPrefix(url, RestRootHandler{Root: root}))

	return rest
}

// Handles requests to the API root.
type RestRootHandler struct {
	Root *PageLoop
}

// Handles requests for application data.
type RestAppHandler struct {
	Root *PageLoop
	Container *model.Container
}

// Gets the list of applications for the API root (/api/).
func (h RestRootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte

	url := req.URL
	path := url.Path

  // TODO: only allow this in Dev mode?
  res.Header().Set("Access-Control-Allow-Origin", "*")

	// List host containers
	if path == "" {
		if req.Method != http.MethodGet {
			HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
			return
		}

		if data, err = json.Marshal(h.Root.Host.Containers); err == nil {
			HttpUtils.Ok(res, data)
			return
		}
  // List available application templates
	} else if path == "templates" {
    // Get built in and user templates
    c := h.Root.Host.GetByName("template")
    u := h.Root.Host.GetByName("user")
    list := append(c.Apps, u.Apps...)
    var apps []*model.Application
    for _, app := range list {
      if app.IsTemplate {
        apps = append(apps, app)
      }
    }

    if content, err := json.Marshal(apps); err != nil {
      HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
      return
    } else {
      HttpUtils.Ok(res, content)
      return
    }
  }

	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		var c *model.Container = h.Root.Host.GetByName(parts[0])
		if c == nil {
			// Container not found
			HttpUtils.Error(res, http.StatusNotFound, nil, nil)
			return
		}

		// Proxy to the app handler

		// Using http.StripPrefix() here does not invoke
		// the underlying handler???
		handler := RestAppHandler{Root: h.Root, Container: c}
		req.URL.Path = strings.TrimPrefix(req.URL.Path, parts[0])
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/")
		handler.ServeHTTP(res, req)
		return
	}

	if err != nil {
		HttpUtils.Error(res, http.StatusInternalServerError, nil, nil)
		return
	}

	// TODO: log the error from (int, error) return value
	HttpUtils.Error(res, http.StatusNotFound, nil, nil)
}

// Handles application information (files, pages etc.)
func (h RestAppHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte
	//var body []byte
	var name string
	var action string
	// File or Page
	var item string
	var app *model.Application
  var file *model.File

	if !HttpUtils.IsMethodAllowed(req.Method, RestAllowedMethods) {
		HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
		return
	}

	url := req.URL
	path := url.Path

	// Check if an app exists when referenced as /api/apps/{name}
	// and extract path parts.
	if path != "" {
		parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
		name = parts[0]
		if len(parts) > 1 {
			action = parts[1]
		}
		if len(parts) > 2 {
			//item = parts[2]
			item = "/" + strings.Join(parts[2:], "/")

      // Respect input trailing slash used to indicate
      // operations on a directory
      if strings.HasSuffix(path, "/") {
        item += "/"
      }
		}
		app = h.Container.GetByName(name)
		// Application must exist
		if app == nil {
			HttpUtils.Error(res, http.StatusNotFound, nil, nil)
			return
		}
	}

	switch req.Method {
		case http.MethodGet:
			if path == "" {
				// GET /api/apps/
				data, err = json.Marshal(h.Container.Apps)
			} else {
				if app != nil {
					if action == "" {
						// GET /api/apps/{name}
						data, err = json.Marshal(app)
					} else {
						switch action {
							case FILES:
								if item == "" {
									// GET /api/apps/{name}/files
									data, err = json.Marshal(app.Files)
								} else {
									// GET /api/apps/{name}/files/{url}
									file := app.GetFileByUrl(item)
									// Data is nil so we send a 404
									if file == nil {
										break
									}
									data, err = json.Marshal(file)
								}
							case PAGES:
								if item == "" {
									// GET /api/apps/{name}/pages
									data, err = json.Marshal(app.Pages)
								} else {
									// GET /api/apps/{name}/pages/{url}
									page := app.GetPageByUrl(item)
									// Data is nil so we send a 404
									if page == nil {
										break
									}
									data, err = json.Marshal(page)
								}
							default:
								HttpUtils.Error(res, http.StatusNotFound, nil, nil)
								return
						}
					}
				}
			}
		// DELETE /api/{container}/{name}/
		case http.MethodDelete:
			if name != "" && action == "" {
				if app.Protected {
					HttpUtils.Error(res, http.StatusForbidden, nil, errors.New("Cannot delete protected application"))
					return
				}

        // Stop serving files for the application
        h.Root.UnmountApplication(app)

        // Delete the mountpoint
        if err = h.Root.DeleteApplicationMountpoint(app); err != nil {
					HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
					return
        }

        // Delete the files
        if err = h.Root.DeleteApplicationFiles(app); err != nil {
					HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
					return
        }

        // Delete the in-memory application
        h.Container.Del(app)

				HttpUtils.Ok(res, OK)
				return
      // DELETE /api/{container}/{app}/files/ - Bulk file deletion
      } else if action == FILES && item == "" {
        var urls UrlList
        var content []byte

        if content, err = HttpUtils.ReadBody(req); err != nil {
          HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
          return
        }

        if err = json.Unmarshal(content, &urls); err != nil {
          HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
          return
        }

        for _, url := range urls {
          if file = h.deleteFile(url, app, res, req); file == nil {
            // If we got a nil file an error occured and the response
            // will already have been sent
            return
          }
        }

        // If we made it this far all files were deleted
        HttpUtils.Ok(res, OK)
				return
			} else if action == FILES && item != "" {
        if file = h.deleteFile(item, app, res, req); file != nil {
          HttpUtils.Ok(res, OK)
        }
				return
			} else {
				HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
				return
			}
		// PUT /api/{container}/
		case http.MethodPut:
			if path == "" {
				var input *model.Application = &model.Application{}
				_, err = HttpUtils.ValidateRequest(SchemaAppNew, input, req)
				if err != nil {
					HttpUtils.Error(res, http.StatusBadRequest, nil, err)
					return
				}

        existing := h.Container.GetByName(input.Name)
        if existing != nil {
					HttpUtils.Error(res, http.StatusPreconditionFailed, nil, fmt.Errorf("Application %s already exists", input.Name))
					return
        }

        input.Url = input.MountpointUrl(h.Container)

        // mountpoint exists
        exists := h.Root.HasMountpoint(input.Url)
        if exists {
					HttpUtils.Error(res, http.StatusPreconditionFailed, nil, fmt.Errorf("Mountpoint URL %s already exists", input.Url))
					return
        }

        var mountpoint *Mountpoint

        // Create and save a mountpoint for the application.
        if mountpoint, err = h.Root.CreateMountpoint(input); err != nil {
					HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
					return
        }

        if input.Template != nil {
          var source *model.Application

          // Find the template app/ directory
          if source, err = h.Root.LookupTemplate(input.Template); err != nil {
            HttpUtils.Error(res, http.StatusBadRequest, nil, err);
            return
          }

          // Copy template source files
          if err = h.Root.CopyApplicationTemplate(input, source); err != nil {
            HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
            return
          }
        }

        var app *model.Application

        // Load and publish the app source files
        if app, err = h.Root.LoadMountpoint(*mountpoint, h.Container); err != nil {
          HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
          return
        }

        // Mount the application
        h.Root.MountApplication(app)

				HttpUtils.Created(res, OK)
				return
			} else {
				// PUT /api/{container}/{app}/files/{url}
				if name != "" && action == FILES && item != "" {
					if file = h.putFile(item, app, res, req); file != nil {
            HttpUtils.OkFile(http.StatusCreated, res, file)
          }
					return
				} else if (name != "" && action == TASKS && item != "") {
          taskName := strings.TrimPrefix(item, "/")
          taskName = strings.TrimSuffix(taskName, "/")

          // No build configuration of missing build task
          if !app.HasBuilder() || app.Builder.Tasks[taskName] == "" {
            HttpUtils.Error(res, http.StatusNotFound, nil, nil)
            return
          }

          fullName := fmt.Sprintf("%s/%s:%s", app.Container.Name, app.Name, taskName)

          if Jobs.GetRunningJob(fullName) != nil {
            HttpUtils.Error(res, http.StatusConflict, nil, fmt.Errorf("Job %s is already running", fullName))
            return
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
          HttpUtils.Write(res, http.StatusAccepted, OK)

          return
        }

				HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
				return
			}
		case http.MethodPost:
			// POST /api/{container}/{app}/files/{url}
			if name != "" && action == FILES && item != "" {
				if file = h.postFile(item, app, res, req); file != nil {
          HttpUtils.OkFile(http.StatusOK, res, file)
        }
				return
			}
			HttpUtils.Error(res, http.StatusMethodNotAllowed, nil, nil)
			return
	}

	if err != nil {
		HttpUtils.Error(res, http.StatusInternalServerError, nil, err)
		return
	}

	if data != nil {
		HttpUtils.Ok(res, data)
		return
	}

	HttpUtils.Error(res, http.StatusNotFound, nil, nil)
}

func (h RestAppHandler) deleteFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) *model.File {
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

// Create a new file for an application
func (h RestAppHandler) putFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) *model.File {
	var err error

	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

	if ct == "" {
    ct = mime.TypeByExtension(filepath.Ext(req.URL.Path))
	}

	// No content length header
	if cl == "" {
		HttpUtils.Error(res, http.StatusBadRequest, nil, errors.New("Content length header is required"))
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

// Update file content for an application
func (h RestAppHandler) postFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) *model.File {
	var err error
	loc := req.Header.Get("Location")
	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

  if loc == "" {
    // No content type header
    if ct == "" {
      HttpUtils.Error(res, http.StatusBadRequest, nil, errors.New("Content type header is required"))
      return nil
    }

    // No content length header
    if cl == "" {
      HttpUtils.Error(res, http.StatusBadRequest, nil, errors.New("Content length header is required"))
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
      // okFile(http.StatusOK, res, file)

      return file
    // Update file content
    } else {
      // Strip charset for mime comparison
      ct = CharsetStrip.ReplaceAllString(ct, "")
      ft := CharsetStrip.ReplaceAllString(file.Mime, "")
      if ft != ct {
        HttpUtils.Error(res, http.StatusBadRequest, nil, errors.New("Mismatched MIME types attempting to update file"))
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
