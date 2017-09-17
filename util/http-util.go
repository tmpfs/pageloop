package util

import (
	"errors"
  "strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
)

const(
	JSON_MIME = "application/json; charset=utf-8"
	SLASH = "/"
)

// Utilities for the REST API endpoints.
type HttpUtil struct {}

// Send an error response to the client as JSON.
func (h HttpUtil) Errorj(res http.ResponseWriter, ex *StatusError) (int, error) {
  message := ex.Error()
  var m map[string] interface{} = make(map[string] interface{})
  m["code"] = ex.Status
  m["message"] = http.StatusText(ex.Status)
  if message != m["message"] {
    m["error"] = message
  }
  if data, err := json.Marshal(m); err != nil {
    return -1, err
  } else {
	  return h.Write(res, ex.Status, data)
  }
}

// Read in a request body.
func (h HttpUtil) ReadBody(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	return ioutil.ReadAll(req.Body)
}

func (h HttpUtil) ReadJson(req *http.Request, input interface{}) *StatusError {
  if content, err := h.ReadBody(req); err != nil {
    return CommandError(http.StatusInternalServerError, err.Error())
  } else {
    if err = json.Unmarshal(content, input); err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }
  }
  return nil
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
