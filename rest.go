// Handler for the REST API.

package pageloop

import (
	//"fmt"
	"log"
	//"errors"
	"strings"
	"net/http"
	"encoding/json"
  "github.com/tmpfs/pageloop/model"
)

const(
	JSON_MIME = "application/json"
	HTML_MIME = "text/html"

	// App actions
	FILES = "files"
	PAGES = "pages"
)

var (
	JSON_404 []byte = []byte(`{status:404}`)
)

type RestService struct {
	Root *PageLoop
}

func (r *RestService) Multiplex(mux *http.ServeMux) {
	var url string
	url= "/api/"
	mux.Handle(url, http.StripPrefix(url, RootHandler{Root: r.Root}))

	url = "/api/app/"
	mux.Handle(url, http.StripPrefix(url, AppHandler{Root: r.Root}))
}

type RootHandler struct {
	Root *PageLoop
}

type AppHandler struct {
	Root *PageLoop
}

// Gets the list of applications.
func (h RootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte

	if req.Method != http.MethodGet {
		ex(res, http.StatusMethodNotAllowed, nil)
		return
	}

	url := req.URL
	path := url.Path

	// Api root (/api/)
	if path == "" {
		if data, err = json.Marshal(h.Root.Container); err == nil {
			ok(res, data)
			return
		}
	}

	if err != nil {
		log.Printf("Internal server error: %s", err.Error())
		ex(res, http.StatusInternalServerError, nil)
		return
	}

	// TODO: log the error from (int, error) return value
	ex(res, http.StatusNotFound, nil)
}

// Handles application information (files, pages etc.)
func (h AppHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte
	var name string
	var action string
	var app *model.Application
	var methods []string = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	if !isMethodAllowed(req.Method, methods) {
		ex(res, http.StatusMethodNotAllowed, nil)
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
		app = h.Root.Container.GetByName(name)
		// Application must exist
		if app == nil {
			ex(res, http.StatusNotFound, nil)
			return
		}
	}

	switch req.Method {
		case http.MethodGet:
			if path == "" {
				if data, err = json.Marshal(h.Root.Container.Apps); err == nil {
					ok(res, data)
					return
				}
			} else {
				if app != nil {
					if action == "" {
						if data, err = json.Marshal(app); err == nil {
							ok(res, data)
							return
						}
					} else {
						switch action {
							case FILES:
								if data, err = json.Marshal(app.Files); err == nil {
									ok(res, data)
									return
								}
							case PAGES:
								if data, err = json.Marshal(app.Pages); err == nil {
									ok(res, data)
									return
								}
						}
					}
				}
			}
	}

	ex(res, http.StatusNotFound, nil)
}

// Send an error response to the client.
func ex(res http.ResponseWriter, code int, data []byte) (int, error) {
	var err error
	if data == nil {
		var m map[string] interface{} = make(map[string] interface{})
		m["code"] = code
		m["message"] = http.StatusText(code)
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

func isMethodAllowed(method string, methods []string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}	
	return false
}
