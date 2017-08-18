// Exposes a REST API to the application model.

package pageloop

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
	"net/http"
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

				h.Container.Del(app)

				// TODO: rewrite mountpoints
				// TODO: persist application mountpoints
				// TODO: unmount application

				ok(res, OK)
				return
			} else if action == FILES && item != "" {
				var file *model.File = app.Urls[item]
				if file == nil {
					ex(res, http.StatusNotFound, nil, nil)
					return
				}

				if err = app.Del(file); err != nil {
					ex(res, http.StatusInternalServerError, nil, nil)
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

				// Add the application to the container.
				if err = h.Container.Add(input); err != nil {
					ex(res, http.StatusPreconditionFailed, nil, err)
					return
				}

				// TODO: create application source file directory - template??
				// TODO: rewrite mountpoints
				// TODO: persist application mountpoints

				created(res, OK)
				return
			} else {
				// PUT /api/{container}/{app}/{url}
				if name != "" && action == FILES && item != "" {
					putFile(item, app, res, req)
					return

				}

				ex(res, http.StatusMethodNotAllowed, nil, nil)
				return
			}
		case http.MethodPost:
			println("got post request")

			// POST /api/{container}/{app}/{url}
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

	var file *model.File = app.Urls[url]
	output := app.GetPathFromUrl(url)

	if file != nil {
		ex(res, http.StatusPreconditionFailed, nil, errors.New("File already exists"))
		return
	}

	// TODO: fix empty reply when there is no request body
	// TODO: stream request body to disc
	var content []byte
	if content, err = readBody(req); err == nil {
		// Update the application model
		if _, err = app.Create(output, content); err != nil {
			ex(res, http.StatusInternalServerError, nil, err)
			return
		}
		created(res, OK)
		return
	}

	/*
	// Be certain the file does not exist on disc
	fh, err := os.Open(output)
	if err != nil {
		if os.IsNotExist(err) {
			// Try to create parent directories
			if err = os.MkdirAll(dir, os.ModeDir | 0755); err != nil {
				ex(res, http.StatusInternalServerError, nil, err)
				return
			}
			// Create the destination file
			if fh, err = os.Create(output); err != nil {
				ex(res, http.StatusInternalServerError, nil, err)
				return
			}

			defer fh.Close()
			var stat os.FileInfo

			if stat, err = fh.Stat(); err != nil {
				ex(res, http.StatusInternalServerError, nil, err)
				return
			}

			mode := stat.Mode()
			if mode.IsDir() {
				ex(res, http.StatusForbidden, nil, errors.New("Attempt to PUT a file to an existing directory"))
				return
			} else if mode.IsRegular() {
				fh, err := os.Create(output)
				if err == nil {
					defer fh.Close()

				}
			}
		}

		ex(res, http.StatusInternalServerError, nil, err)
		return
	}
	defer fh.Close()

	*/

	if err != nil {
		ex(res, http.StatusInternalServerError, nil, err)
		return
	}

	ex(res, http.StatusNotFound, nil, nil)
}

// Update file content for an application
func postFile(url string, app *model.Application, res http.ResponseWriter, req *http.Request) {
	var err error
	ct := req.Header.Get("Content-Type")
	cl := req.Header.Get("Content-Length")

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

	var file *model.File = app.Urls[url]


	if file != nil {
		// Strip charset for mime comparison
		ct = CharsetStrip.ReplaceAllString(ct, "")
		if file.Mime != ct {
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
			ok(res, OK)
			return
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
