package pageloop

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"encoding/json"
  "testing"
)

var url string = "http://localhost:3579"
var api string = url + "/api"
var rpcUrl string = url + "/rpc/"
var appApi string = api + "/apps"

var server *http.Server

// Test call to listen without a server
func TestListenError(t *testing.T) {
	var err error
  loop := &PageLoop{}
	err = loop.Listen(nil)
	if err == nil {
		t.Fatal("Expected error response from call to listen without server")
	}

	conf := ServerConfig{Addr: ":443", Dev: true}
	if server, err = loop.NewServer(conf); err != nil {
		t.Fatal(err)
	}

	defer server.Close()

	var c chan error = make(chan error)
	go func(ch chan<- error) { err = loop.Listen(server); if err != nil {ch <-err}; close(ch)} (c)
	err = <-c
	if err == nil {
		t.Fatal("Expected error response from call to listen with port under 1024")
	}
}

// Start a mock server running for subsequent tests.
func TestStartServer(t *testing.T) {
	var err error
  var apps []Mountpoint
	apps = append(apps, Mountpoint{UrlPath: "/app/mock-app/", Path: "test/fixtures/mock-app"})
  loop := &PageLoop{}
	conf := ServerConfig{Mountpoints: apps, Addr: ":3579", Dev: true}
	if server, err = loop.NewServer(conf); err != nil {
		t.Fatal(err)
	}
	go loop.Listen(server)
}

// Test GET to the home application
func TestMainPages(t *testing.T) {
	var err error
	var resp *http.Response

	// GET /
	if resp, _, err = get(url + "/"); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, HTML_MIME)

	// GET /-/source/
	if resp, _, err = get(url + "/-/source/"); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, HTML_MIME)

	// GET /app/mock-app/
	if resp, _, err = get(url + "/app/mock-app/"); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, HTML_MIME)

	// GET /app/mock-app/-/source/
	if resp, _, err = get(url + "/app/mock-app/-/source/"); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, HTML_MIME)
}

// Test GET for 404 responses
func TestNotFound(t *testing.T) {
	var err error
	var resp *http.Response
	// GET /not-found/
	if resp, _, err = get(url + "/not-found/"); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusNotFound)

	// GET /api/not-found/
	if resp, _, err = get(api + "/not-found/"); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusNotFound)
}

// Test RPC server 
func TestRpcService(t *testing.T) {
	var err error
	var resp *http.Response
	var body []byte
	var doc []byte

	doc = []byte(`{"id": 0, "method": "app.List", "params": [{"index": 0}]}`)
	if resp, body, err = post(rpcUrl, JSON_MIME, doc); err != nil {
		t.Fatal(err)
	}

	println(string(body))
	println(resp)
}

// Test REST API endpoints
func TestRestService(t *testing.T) {
	var err error
	var resp *http.Response
	var body []byte
	var doc []byte

	var res map[string] interface{} = make(map[string] interface{})
	var apps []interface{}
	var list []interface{}
	var app map [string] interface{}
	var file map [string] interface{}
	var page map [string] interface{}
	var name string
	var ok bool


	// GET /api/
	if resp, body, err = get(api); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, JSON_MIME)

	if err = json.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}

	if apps, ok = res["apps"].([]interface{}); !ok {
		t.Error("Unexpected type for apps list")
	}

	if app, ok = apps[0].(map [string] interface{}); !ok {
		t.Error("Unexpected type for app")
	}

	if _, ok = app["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}

	// GET /api/apps/
	if resp, body, err = get(fmt.Sprintf("%s%s", api, "/apps/")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, JSON_MIME)

	if err = json.Unmarshal(body, &list); err != nil {
		t.Fatal(err)
	}

	if app, ok = list[0].(map [string] interface{}); !ok {
		t.Error("Unexpected type for app")
	}

	if name, ok = app["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}

	// GET /api/apps/{name}/
	if resp, body, err = get(fmt.Sprintf("%s%s%s", api, "/apps/", name + "/")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, JSON_MIME)

	app = make(map[string] interface{})
	if err = json.Unmarshal(body, &app); err != nil {
		t.Fatal(err)
	}

	if _, ok = app["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}

	// GET /api/apps/{name}/files/
	if resp, body, err = get(fmt.Sprintf("%s%s%s%s", api, "/apps/", name, "/files/")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, JSON_MIME)

	list = make([]interface{}, 128)
	if err = json.Unmarshal(body, &list); err != nil {
		t.Fatal(err)
	}

	if file, ok = list[0].(map [string] interface{}); !ok {
		t.Error("Unexpected type for file")
	}

	if _, ok = file["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}

	if _, ok = file["url"].(string); !ok {
		t.Error("Unexpected type for url")
	}

	// GET /api/apps/{name}/files/{url}/
	if resp, body, err = get(fmt.Sprintf("%s%s%s%s", api, "/apps/", name, "/files/index.html")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, JSON_MIME)

	file = make(map [string]interface{})
	if err = json.Unmarshal(body, &file); err != nil {
		t.Fatal(err)
	}

	if _, ok = file["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}

	if _, ok = file["url"].(string); !ok {
		t.Error("Unexpected type for url")
	}

	if _, ok = file["size"].(float64); !ok {
		t.Error("Unexpected type for size")
	}

	// GET /api/apps/{name}/pages/
	if resp, body, err = get(fmt.Sprintf("%s%s%s%s", api, "/apps/", name, "/pages/")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, JSON_MIME)

	list = make([]interface{}, 128)
	if err = json.Unmarshal(body, &list); err != nil {
		t.Fatal(err)
	}

	if page, ok = list[0].(map [string] interface{}); !ok {
		t.Error("Unexpected type for page")
	}

	if _, ok = page["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}

	if _, ok = page["url"].(string); !ok {
		t.Error("Unexpected type for url")
	}

	if _, ok = page["size"].(float64); !ok {
		t.Error("Unexpected type for size")
	}

	// GET /api/apps/{name}/pages/{url}
	if resp, body, err = get(fmt.Sprintf("%s%s%s%s", api, "/apps/", name, "/pages/index.html")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)
	assertContentType(resp, t, JSON_MIME)

	page = make(map [string]interface{})
	if err = json.Unmarshal(body, &page); err != nil {
		t.Fatal(err)
	}

	if _, ok = page["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}

	if _, ok = page["url"].(string); !ok {
		t.Error("Unexpected type for url")
	}

	if _, ok = page["size"].(float64); !ok {
		t.Error("Unexpected type for size")
	}

	// PUT /api/apps/ - Created
	doc = []byte(`{"name": "test-app"}`)
	if resp, body, err = put(api + "/apps/", doc); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusCreated)

	// GET /api/apps/test-app/ - OK
	if resp, body, err = get(api + "/apps/test-app/"); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)

	// DELETE /api/apps/test-app/ - OK
	if resp, body, err = del(api + "/apps/test-app/", nil); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusOK)

	// GET /api/apps/test-app/ - Not Found
	if resp, body, err = get(api + "/apps/test-app/"); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusNotFound)

	////
	// Error conditions
	////

	// TRACE /api/ - Method Not Allowed
	doc = []byte(``)
	if resp, body, err = do(appApi + "/", http.MethodTrace, doc); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusMethodNotAllowed)

	// POST /api/ - Method Not Allowed
	doc = []byte(`{}`)
	if resp, body, err = post(api + "/", JSON_MIME, doc); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusMethodNotAllowed)

	// DELETE /api/apps/ - Method Not Allowed
	doc = []byte(`{}`)
	if resp, body, err = del(appApi + "/", nil); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusMethodNotAllowed)

	// PUT /api/mock-app/ - Method Not Allowed
	doc = []byte(`{}`)
	if resp, body, err = put(appApi + "/mock-app/", doc); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusMethodNotAllowed)

	// PUT /api/ (malformed json) - Bad Request
	doc = []byte(`{`)
	if resp, body, err = put(api + "/apps/", doc); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusBadRequest)

	// PUT /api/ (schema validation fail) - Bad Request
	doc = []byte(`{}`)
	if resp, body, err = put(api + "/apps/", doc); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusBadRequest)

	// PUT /api/ (app exists) - Precondition Failed
	doc = []byte(`{"name": "mock-app"}`)
	if resp, body, err = put(api + "/apps/", doc); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusPreconditionFailed)

	// PUT /api/ (invalid app name) - Precondition Failed
	doc = []byte(`{"name": "-app"}`)
	if resp, body, err = put(api + "/apps/", doc); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusPreconditionFailed)

	// GET /api/apps/{name}/invalid-action/ (invalid action) - Not Found
	if resp, body, err = get(fmt.Sprintf("%s%s%s%s", api, "/apps/", name, "/invalid-action/")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusNotFound)

	// GET /api/apps/{name}/files/not-found/ (missing file) - Not Found
	if resp, body, err = get(fmt.Sprintf("%s%s%s%s", api, "/apps/", name, "/files/not-found/")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusNotFound)

	// GET /api/apps/{name}/pages/not-found/ (missing page) - Not Found
	if resp, body, err = get(fmt.Sprintf("%s%s%s%s", api, "/apps/", name, "/pages/not-found/")); err != nil {
		t.Fatal(err)
	}
	assertStatus(resp, t, http.StatusNotFound)
}

// Private helpers

func get(url string) (*http.Response, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, body, nil
}


func del(url string, data []byte) (*http.Response, []byte, error) {
	if data == nil {
		data = make([]byte, 0)
	}
	return do(url, http.MethodDelete, data)
}

func post(url string, contentType string, body []byte) (*http.Response, []byte, error) {
	var buf = new(bytes.Buffer)
	buf.Write(body)
	resp, err := http.Post(url, contentType, buf)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, resbody, nil
}

func put(url string, body []byte) (*http.Response, []byte, error) {
	return do(url, http.MethodPut, body)
}

func do(uri string, method string, body []byte) (*http.Response, []byte, error) {
	var err error
	var buf = new(bytes.Buffer)
	buf.Write(body)
	var req *http.Request
	if req, err = http.NewRequest(method, uri, buf); err != nil {
		return nil, nil, err	
	}
	req.Header.Set("Content-Type", JSON_MIME)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, resbody, nil
}

func assertContentType(resp *http.Response, t *testing.T, mime string) {
	ct := resp.Header.Get("Content-Type")
	if ct != mime {
		t.Errorf("Unexpected response content type %s", ct)
	}
}

func assertStatus(resp *http.Response, t *testing.T, code int) {
	if resp.StatusCode != code {
		t.Errorf("Unexpected status code %d wanted %d", resp.StatusCode, code)
	}
}
