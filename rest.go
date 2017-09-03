// Exposes a REST API to the application model.

package pageloop

import (
  "fmt"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
	"net/http"
  "mime"
  pth "path"
  "path/filepath"
	"encoding/json"
  "github.com/tmpfs/pageloop/model"
	"github.com/xeipuuv/gojsonschema"
)

const(
	API_URL = "/api/"
	JSON_MIME = "application/json; charset=utf-8"

	// App actions
	FILES = "files"
	PAGES = "pages"
)

var(
	OK = []byte(`{"ok": true}`)
	SchemaAppNew = MustAsset("schema/app-new.json")
	CharsetStrip = regexp.MustCompile(`;.*$`)
  TemplateNewFile map[string] []byte
)

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

  res.Header().Set("Access-Control-Allow-Origin", "*")

	// List host containers
	if path == "" {
		if req.Method != http.MethodGet {
			ex(res, http.StatusMethodNotAllowed, nil, nil)
			return
		}

		if data, err = json.Marshal(h.Root.Host.Containers); err == nil {
			ok(res, data)
			return
		}
	}

	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		var c *model.Container = h.Root.Host.GetByName(parts[0])
		if c == nil {
			// Container not found
			ex(res, http.StatusNotFound, nil, nil)
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
		ex(res, http.StatusInternalServerError, nil, nil)
		return
	}

	// TODO: log the error from (int, error) return value
	ex(res, http.StatusNotFound, nil, nil)
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
	var methods []string = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	if !isMethodAllowed(req.Method, methods) {
		ex(res, http.StatusMethodNotAllowed, nil, nil)
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
			ex(res, http.StatusNotFound, nil, nil)
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
								ex(res, http.StatusNotFound, nil, nil)
								return
						}
					}
				}
			}
		// DELETE /api/{container}/{name}/
		case http.MethodDelete:
			if name != "" && action == "" {
				if app.Protected {
					ex(res, http.StatusForbidden, nil, errors.New("Cannot delete protected application"))
					return
				}

        // Stop serving files for the application
        h.Root.UnmountApplication(app)

        // Delete the mountpoint
        if err = h.Root.DeleteApplicationMountpoint(app); err != nil {
					ex(res, http.StatusInternalServerError, nil, err)
					return
        }

        // Delete the files
        if err = h.Root.DeleteApplicationFiles(app); err != nil {
					ex(res, http.StatusInternalServerError, nil, err)
					return
        }

        // Delete the in-memory application
        h.Container.Del(app)

				ok(res, OK)
				return
			} else if action == FILES && item != "" {
				var file *model.File = app.Urls[item]
				if file == nil {
					ex(res, http.StatusNotFound, nil, nil)
					return
				}

				if err = app.Del(file); err != nil {
					ex(res, http.StatusInternalServerError, nil, err)
					return
				}

				ok(res, OK)
				return
			} else {
				ex(res, http.StatusMethodNotAllowed, nil, nil)
				return
			}
		// PUT /api/{container}/
		case http.MethodPut:
			if path == "" {
				var input *model.Application = &model.Application{}
				_, err = validateRequest(SchemaAppNew, input, req)
				if err != nil {
					ex(res, http.StatusBadRequest, nil, err)
					return
				}

        input.Url = pth.Clean(input.Url)

        if !strings.HasSuffix(input.Url, "/") {
          input.Url += "/"
        }

        existing := h.Container.GetByName(input.Name)
        if existing != nil {
					ex(res, http.StatusPreconditionFailed, nil, fmt.Errorf("Application %s already exists", input.Name))
					return
        }

        // mountpoint exists
        exists := h.Root.HasMountpoint(input.Url)
        if exists {
					ex(res, http.StatusPreconditionFailed, nil, fmt.Errorf("Mountpoint URL %s already exists", input.Url))
					return
        }

        var mountpoint *Mountpoint

        // Create and save a mountpoint for the application.
        if mountpoint, err = h.Root.CreateMountpoint(input); err != nil {
					ex(res, http.StatusInternalServerError, nil, err)
					return
        }

        if input.Template != nil {
          var source *model.Application

          // Find the template app/ directory
          if source, err = h.Root.LookupTemplate(input.Template); err != nil {
            ex(res, http.StatusBadRequest, nil, err);
            return
          }

          // Copy template source files
          if err = h.Root.CopyApplicationTemplate(input, source); err != nil {
            ex(res, http.StatusInternalServerError, nil, err)
            return
          }
        }

        var app *model.Application

        // Load and publish the app source files
        if app, err = h.Root.LoadMountpoint(*mountpoint, h.Container); err != nil {
          ex(res, http.StatusInternalServerError, nil, err)
          return
        }

        // Mount the application
        h.Root.MountApplication(app)

				created(res, OK)
				return
			} else {
				// PUT /api/{container}/{app}/files/{url}
				if name != "" && action == FILES && item != "" {
					putFile(item, app, res, req)
					return
				}

				ex(res, http.StatusMethodNotAllowed, nil, nil)
				return
			}
		case http.MethodPost:
			// POST /api/{container}/{app}/files/{url}
			if name != "" && action == FILES && item != "" {
				postFile(item, app, res, req)
				return
			}
			ex(res, http.StatusMethodNotAllowed, nil, nil)
			return
	}

	if err != nil {
		ex(res, http.StatusInternalServerError, nil, err)
		return
	}

	if data != nil {
		ok(res, data)
		return
	}

	ex(res, http.StatusNotFound, nil, nil)
}

// Create a new file for an application
func putFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) {
	var err error

	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

	if ct == "" {
    ct = mime.TypeByExtension(filepath.Ext(req.URL.Path))
	}

	// No content length header
	if cl == "" {
		ex(res, http.StatusBadRequest, nil, errors.New("Content length header is required"))
		return
	}

  isDir := strings.HasSuffix(url, "/")

	var content []byte
  content = TemplateNewFile[ct]
  // Read content from request body if no template available
  // and not operating on a directory
  if !isDir && content == nil {
    // TODO: fix empty reply when there is no request body
    // TODO: stream request body to disc
    if content, err = readBody(req); err != nil {
      ex(res, http.StatusInternalServerError, nil, err)
      return
    }
  }

  // Update the application model
  if _, err = app.Create(url, content); err != nil {
    if err, ok := err.(model.StatusError); ok {
      ex(res, err.Status, nil, err)
      return
    }

    ex(res, http.StatusInternalServerError, nil, err)
    return
  }
  created(res, OK)
  return

	if err != nil {
		ex(res, http.StatusInternalServerError, nil, err)
		return
	}

	ex(res, http.StatusNotFound, nil, nil)
}

// Update file content for an application
func postFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) {
	var err error
	loc := req.Header.Get("Location")
	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

  if loc == "" {
    // No content type header
    if ct == "" {
      ex(res, http.StatusBadRequest, nil, errors.New("Content type header is required"))
      return
    }

    // No content length header
    if cl == "" {
      ex(res, http.StatusBadRequest, nil, errors.New("Content length header is required"))
      return
    }
  }

	var file *model.File = app.Urls[url]
	if file != nil {

    // Handle moving the file with Location header
    if loc != "" {
      if url == loc {
        ex(res, http.StatusBadRequest, nil,
          fmt.Errorf("Cannot move file, source and destination are equal: %s", url))
        return
      }

      if err = app.Move(file, loc); err != nil {
        ex(res, http.StatusInternalServerError, nil, err)
        return
      }
      okFile(res, file)

      return
    // Update file content
    } else {
      // Strip charset for mime comparison
      ct = CharsetStrip.ReplaceAllString(ct, "")
      ft := CharsetStrip.ReplaceAllString(file.Mime, "")
      if ft != ct {
        ex(res, http.StatusBadRequest, nil, errors.New("Mismatched MIME types attempting to update file"))
        return
      }

      // TODO: fix empty reply when there is no request body
      // TODO: stream request body to disc
      var content []byte
      if content, err = readBody(req); err == nil {
        // Update the application model
        if err = app.Update(file, content); err != nil {
          ex(res, http.StatusInternalServerError, nil, err)
          return
        }
        okFile(res, file)
        return
      }

    }
	}

	if err != nil {
		ex(res, http.StatusInternalServerError, nil, err)
		return
	}

	ex(res, http.StatusNotFound, nil, nil)
}


// Send an error response to the client.
func ex(res http.ResponseWriter, code int, data []byte, exception error) (int, error) {
	var err error
	if data == nil {
		var m map[string] interface{} = make(map[string] interface{})
		m["code"] = code
		m["message"] = http.StatusText(code)
		if exception != nil {
			m["error"] = exception.Error()
		}
		if data, err = json.Marshal(m); err != nil {
			return 0, err
		}
	}
	return write(res, code, data)
}

// Private helper functions.

func readBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	return ioutil.ReadAll(req.Body)
}

// Validate a client request.
//
// Reads in the request body data, unmarshals to JSON and
// validates the result against the given schema.
func validateRequest(schema []byte, input interface{}, req *http.Request) (*gojsonschema.Result, error) {
	var err error
	var body []byte
	var result *gojsonschema.Result
	body, err = readBody(req)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &input); err != nil {
		return nil, err
	}

	if result, err = validate(schema, body); result != nil {
		if !result.Valid() {
			return nil, errors.New(result.Errors()[0].String())
		}
	}

	return result, nil
}

// Validate client request data.
func validate(schema []byte, input []byte) (*gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(input)
	return gojsonschema.Validate(schemaLoader, documentLoader)
}

// Send an OK response to the client.
func ok(res http.ResponseWriter, data []byte) (int, error) {
	return write(res, http.StatusOK, data)
}

// Send an OK response to the client with a file.
func okFile(res http.ResponseWriter, f *model.File) (int, error) {
  var data []byte
  var err error
  if data, err = json.Marshal(f); err != nil {
    return -1, err
  }
  top := []byte(`{"ok":true,"file":`)
  tail := []byte(`}`)
  data = append(top, data...)
  data = append(data, tail...)
	return write(res, http.StatusOK, data)
}

// Send a created response to the client, typically in reply to a PUT.
func created(res http.ResponseWriter, data []byte) (int, error) {
	return write(res, http.StatusCreated, data)
}

// Write to the HTTP response and set common headers.
func write(res http.ResponseWriter, code int, data []byte) (int, error) {
	res.Header().Set("Content-Type", JSON_MIME)
	res.WriteHeader(code)
	return res.Write(data)
}

// Determine if a method exists in a list of allowed methods.
func isMethodAllowed(method string, methods []string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

func init() {
  TemplateNewFile = make(map[string] []byte)
  TemplateNewFile["template/markdown+partial"] = MustAsset("app/template/documents/partial-markdown.md")
  TemplateNewFile["template/markdown+standalone"] = MustAsset("app/template/documents/standalone-markdown.md")
  TemplateNewFile["template/html+standalone"] = MustAsset("app/template/documents/standalone-html.html")
  TemplateNewFile["template/html+layout"] = MustAsset("app/template/documents/layout.html")
  TemplateNewFile["template/html+partial"] = MustAsset("app/template/documents/partial-html.html")
}
