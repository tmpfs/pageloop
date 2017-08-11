// Handler for the REST API.

package pageloop

import (
	"fmt"
	//"log"
	//"errors"
	"net/http"
  //"path/filepath"
	//"regexp"
  //"time"
  //"github.com/tmpfs/pageloop/model"
)

const(
	JSON_MIME = "application/json"
	HTML_MIME = "text/html"

)

var (
	JSON_404 []byte = []byte(`{status: 404}`)
	JSON_API []byte = []byte(`{version: "1.0"}`)
)

type RestHandler struct {
	// Reference to the data structures
	Loop *PageLoop
}

// The default server handler, defers to a multiplexer.
func (h RestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	var write func(code int, data []byte) (int, error)
	write = func(code int, data []byte) (int, error) {
		res.WriteHeader(code)
		return res.Write(data)
	}

	url := req.URL
	path := url.Path
	res.Header().Set("Content-Type", JSON_MIME)

	fmt.Printf("%#v\n", req)
	fmt.Printf("%#v\n", path)

	println(req.Method)

	// Api root (/)
	switch req.Method {
		case http.MethodGet:
			if path == "" {
				println("writing home")
				write(http.StatusOK, JSON_API)
				return
			} else if path == "apps" {
			
			}
		default:
			write(http.StatusMethodNotAllowed, nil)
			return
	}

	// TODO: log the error from (int, error) return value
	write(404, JSON_404)
}
