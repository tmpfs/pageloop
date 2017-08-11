// Handler for the REST API.

package pageloop

import (
	//"fmt"
	"log"
	//"errors"
	"net/http"
	"encoding/json"
)

const(
	JSON_MIME = "application/json"
	HTML_MIME = "text/html"

)

var (
	JSON_404 []byte = []byte(`{status:404}`)
)

type RestHandler struct {
	// Reference to the data structures
	Loop *PageLoop
}

// The default server handler, defers to a multiplexer.
func (h RestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	var data []byte

	var write func(code int, data []byte) (int, error)
	write = func(code int, data []byte) (int, error) {
		res.WriteHeader(code)
		return res.Write(data)
	}

	url := req.URL
	path := url.Path
	res.Header().Set("Content-Type", JSON_MIME)

	//fmt.Printf("%#v\n", req)
	//fmt.Printf("%#v\n", path)

	switch req.Method {
		case http.MethodGet:
			// Api root (/)
			if path == "" {
				if data, err = json.Marshal(h.Loop); err == nil {
					println(string(data))
					write(http.StatusOK, data)
					return
				}
			//} else if path == "apps" {
			
			}
		default:
			write(http.StatusMethodNotAllowed, nil)
			return
	}

	if err != nil {
		log.Println(err)
		// TODO: 500 internal error
		//write(http.StatusInternalServerError, byte[](`{error: "` + string(err) + `"}`))
		return
	}

	// TODO: log the error from (int, error) return value
	write(http.StatusNotFound, JSON_404)
}
