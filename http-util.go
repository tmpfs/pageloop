package pageloop

import (
	"errors"
  "strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
  "github.com/tmpfs/pageloop/model"
	"github.com/xeipuuv/gojsonschema"
)

// Utilities for the REST API endpoints.
type HttpUtil struct {}

// Send an error response to the client.
func (h HttpUtil) Error(res http.ResponseWriter, code int, data []byte, exception error) (int, error) {
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
	return h.Write(res, code, data)
}

// Read in a request body.
func (h HttpUtil) ReadBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	return ioutil.ReadAll(req.Body)
}

// Validate a client request.
//
// Reads in the request body data, unmarshals to JSON and
// validates the result against the given schema.
func (h HttpUtil) ValidateRequest(schema []byte, input interface{}, req *http.Request) (*gojsonschema.Result, error) {
	var err error
	var body []byte
	var result *gojsonschema.Result
	body, err = h.ReadBody(req)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &input); err != nil {
		return nil, err
	}

	if result, err = h.Validate(schema, body); result != nil {
		if !result.Valid() {
			return nil, errors.New(result.Errors()[0].String())
		}
	}

	return result, nil
}

// Validate client request data.
func (h HttpUtil) Validate(schema []byte, input []byte) (*gojsonschema.Result, error) {
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(input)
	return gojsonschema.Validate(schemaLoader, documentLoader)
}

// Send an OK response to the client.
func (h HttpUtil) Ok(res http.ResponseWriter, data []byte) (int, error) {
	return h.Write(res, http.StatusOK, data)
}

// Send an OK response to the client with a file.
func (h HttpUtil) OkFile(status int, res http.ResponseWriter, f *model.File) (int, error) {
  var data []byte
  var err error
  var target interface{} = f
  if f.Page() != nil {
    target = f.Page()
  }
  if data, err = json.Marshal(target); err != nil {
    return -1, err
  }
  top := []byte(`{"ok":true,"file":`)
  tail := []byte(`}`)
  data = append(top, data...)
  data = append(data, tail...)
	return h.Write(res, status, data)
}

// Send a created response to the client, typically in reply to a PUT.
func (h HttpUtil) Created(res http.ResponseWriter, data []byte) (int, error) {
	return h.Write(res, http.StatusCreated, data)
}

// Write a JSON document to the response from the given doc object.
func (h HttpUtil) Json(res http.ResponseWriter, status int, doc interface{}) (int, error) {
  var data []byte
  var err error
  if data, err = json.Marshal(doc); err != nil {
    return -1, err
  }
	res.Header().Set("Content-Type", JSON_MIME)
	res.Header().Set("Content-Length", strconv.Itoa(len(data)))
  return h.Write(res, status, data)
}

// Write to the HTTP response.
func (h HttpUtil) Write(res http.ResponseWriter, status int, data []byte) (int, error) {
	res.WriteHeader(status)
	return res.Write(data)
}

// Determine if a method exists in a list of allowed methods.
func (h HttpUtil) IsMethodAllowed(method string, methods []string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}
