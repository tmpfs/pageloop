package pageloop

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
  "testing"
)

var url string = "http://localhost:3579"
var api string = url + "/api"
var app string = url + "/app"

func TestStartServer(t *testing.T) {
	var err error
	var server *http.Server
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
	assertHeaders(resp, t, http.StatusOK)

	// GET /-/source/
	if resp, _, err = get(url + "/-/source/"); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t, http.StatusOK)

	// GET /app/mock-app/
	if resp, _, err = get(app + "/mock-app/"); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t, http.StatusOK)

	// GET /app/mock-app/-/source/
	if resp, _, err = get(app + "/mock-app/-/source/"); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t, http.StatusOK)
}

// Test GET for 404 responses
func TestNotFound(t *testing.T) {
	var err error
	var resp *http.Response
	// GET /not-found/
	if resp, _, err = get(url + "/not-found/"); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t, http.StatusNotFound)

	// GET /api/not-found/
	if resp, _, err = get(api + "/not-found/"); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t, http.StatusNotFound)
}

// Test REST API endpoints
func TestRestService(t *testing.T) {
	var err error
	var resp *http.Response
	var body []byte

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
	assertHeaders(resp, t, http.StatusOK)

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
	assertHeaders(resp, t, http.StatusOK)

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
	assertHeaders(resp, t, http.StatusOK)

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
	assertHeaders(resp, t, http.StatusOK)

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

	// GET /api/apps/{name}/pages/
	if resp, body, err = get(fmt.Sprintf("%s%s%s%s", api, "/apps/", name, "/pages/")); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t, http.StatusOK)

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

	/*
	if _, ok = page["size"].(int); !ok {
		t.Error("Unexpected type for size")
	}
	*/
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

func assertHeaders(resp *http.Response, t *testing.T, code int) {
	if resp.StatusCode != code {
		t.Errorf("Unexpected status code %d wanted %d", resp.StatusCode, code)
	}
	/*
	if resp.Header.Get("Content-Type") != JSON_MIME {
		t.Error("Unexpected response content type")
	}
	*/
}

