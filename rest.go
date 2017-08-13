// Exposes a REST API to the application model.

package pageloop

import (
	"errors"
	"io/ioutil"
	"strings"
	"net/http"
	"encoding/json"
  "github.com/tmpfs/pageloop/model"
	"github.com/xeipuuv/gojsonschema"
)

const(
	API_URL = "/api/"
	JSON_MIME = "application/json; charset=utf-8"
	//JSON_MIME_UTF8 = "application/json; charset=utf-8"

	// App actions
	FILES = "files"
	PAGES = "pages"
)

var(
	OK = []byte(`{"ok": true}`)
	SchemaAppNew = MustAsset("schema/app-new.json")
)

type RestService struct {
	Root *PageLoop
}

func NewRestService(root *PageLoop, mux *http.ServeMux) *RestService {
	var rest *RestService = &RestService{Root: root}

	var url string
	url= API_URL
	//println("new rest service")
	mux.Handle(url, http.StripPrefix(url, RestRootHandler{Root: root}))

	for key, c := range root.Host.Containers {
		url = API_URL + key + "/apps/"
		mux.Handle(url, http.StripPrefix(url, RestAppHandler{Root: root, Container: c}))
	}

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

	if req.Method != http.MethodGet {
		ex(res, http.StatusMethodNotAllowed, nil, nil)
		return
	}

	url := req.URL
	path := url.Path

	if path == "" {
		if data, err = json.Marshal(h.Root.Host); err == nil {
			ok(res, data)
			return
		}
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
		// DELETE /api/apps/{name}/
		case http.MethodDelete:
			if name != "" && action == "" {
				h.Container.Del(app)

				// TODO: rewrite mountpoints
				// TODO: persist application mountpoints

				ok(res, OK)
				return
			} else {
				ex(res, http.StatusMethodNotAllowed, nil, nil)
				return
			}
		// PUT /api/apps/
		case http.MethodPut:
			if path == "" {
				var input *model.Application = &model.Application{}
				_, err = validateRequest(SchemaAppNew, input, req)
				if err != nil {
					ex(res, http.StatusBadRequest, nil, err)
					return
				}

				//fmt.Printf("PUT input: %#v\n", input)
				//fmt.Printf("PUT result: %#v\n", result)

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
				ex(res, http.StatusMethodNotAllowed, nil, nil)
				return
			}
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

// Validate a client request.
//
// Reads in the request body data, unmarshals to JSON and 
// validates the result against the given schema.
func validateRequest(schema []byte, input interface{}, req *http.Request) (*gojsonschema.Result, error) {
	var err error
	var body []byte
	var result *gojsonschema.Result
	defer req.Body.Close()
	body, err = ioutil.ReadAll(req.Body)
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
