// Exposes a REST API to the application model.

package pageloop

import (
	"fmt"
	//"log"
	"errors"
	"io/ioutil"
	"strings"
	"net/http"
	"encoding/json"
  "github.com/tmpfs/pageloop/model"
	"github.com/xeipuuv/gojsonschema"
)

const(
	JSON_MIME = "application/json"

	// App actions
	FILES = "files"
	PAGES = "pages"

	SchemaAppCreate = `{"properties": {"name": {"type": "string"}}, "required": ["name"], "additionalProperties": false}`
)

type RestService struct {
	Root *PageLoop
}

// Configures the REST API handlers.
func (r *RestService) Multiplex(mux *http.ServeMux) {
	var url string
	url= "/api/"
	mux.Handle(url, http.StripPrefix(url, RestRootHandler{Root: r.Root}))

	url = "/api/apps/"
	mux.Handle(url, http.StripPrefix(url, RestAppHandler{Root: r.Root}))
}

// Handles requests to the API root.
type RestRootHandler struct {
	Root *PageLoop
}

// Handles requests for application data.
type RestAppHandler struct {
	Root *PageLoop
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
		if data, err = json.Marshal(h.Root.Container); err == nil {
			ok(res, data)
			return
		}
	}

	if err != nil {
		//log.Printf("Internal server error: %s", err.Error())
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
	var body []byte
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
		app = h.Root.Container.GetByName(name)
		// Application must exist
		if app == nil {
			ex(res, http.StatusNotFound, nil, nil)
			return
		}
	}

	switch req.Method {
		case http.MethodGet:
			// List applications
			if path == "" {
				data, err = json.Marshal(h.Root.Container.Apps)
			// Operate on an app
			} else {
				if app != nil {
					if action == "" {
						data, err = json.Marshal(app)
					} else {
						switch action {
							case FILES:
								if item == "" {
									data, err = json.Marshal(app.Files)
								} else {
									file := app.GetFileByUrl(item)
									// Data is nil so we send a 404
									if file == nil {
										break
									}
									data, err = json.Marshal(file)
								}
							case PAGES:
								if item == "" {
									data, err = json.Marshal(app.Pages)
								} else {
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

		case http.MethodPut:
			// Only allow PUT at /api/apps/
			if path != "" {
				ex(res, http.StatusMethodNotAllowed, nil, nil)
				return
			} else {
				defer req.Body.Close()
				body, err = ioutil.ReadAll(req.Body)
				//body = req.Body.Read()
				println("put to app " + string(body))

				schemaLoader := gojsonschema.NewStringLoader(SchemaAppCreate)
				documentLoader := gojsonschema.NewBytesLoader(body)

				var result *gojsonschema.Result
				if result, err = gojsonschema.Validate(schemaLoader, documentLoader); result != nil {
					if result.Valid() {
						fmt.Printf("The document is valid\n")
					} else {
						errorList := result.Errors()
						ex(res, http.StatusBadRequest, nil, errors.New(errorList[0].String()))
						return
						/*
						for _, desc := range result.Errors() {
							fmt.Printf("- %s\n", desc)
						}
						*/
					}
				}	
			}
		}

	if err != nil {
		//log.Printf("Internal server error: %s", err.Error())
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

// Send an OK response to the client.
func ok(res http.ResponseWriter, data []byte) (int, error) {
	return write(res, http.StatusOK, data)
}

// Write to the HTTP response and set common headers.
func write(res http.ResponseWriter, code int, data []byte) (int, error) {
	res.Header().Set("Content-Type", JSON_MIME)
	res.WriteHeader(code)
	return res.Write(data)
}

// Determine is a method exists in a list of allowed methods.
func isMethodAllowed(method string, methods []string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}	
	return false
}
