package pageloop

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
  "testing"
)

var url string = "http://localhost:3579"
var api string = url + "/api/"

func Start(t *testing.T) *PageLoop {
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
	return loop
}

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

func assertHeaders(resp *http.Response, t *testing.T) {
	if resp.Header.Get("Content-Type") != JSON_MIME {
		t.Error("Unexpected response content type")
	}
}

func TestRestService(t *testing.T) {
	Start(t)

	//sleep(1)

	var err error
	var resp *http.Response
	var body []byte

	// GET /api/
	if resp, body, err = get(api); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t)

	var res map[string] interface{} = make(map[string] interface{})
	err = json.Unmarshal(body, &res)

	var apps []interface{}
	var app map [string] interface{}
	var name string
	var ok bool

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
	if resp, body, err = get(fmt.Sprintf("%s%s", api, "apps/")); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t)

	var list []interface{}
	err = json.Unmarshal(body, &list)

	if app, ok = list[0].(map [string] interface{}); !ok {
		t.Error("Unexpected type for app")
	}

	if name, ok = app["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}

	// GET /api/apps/{name}/
	if resp, body, err = get(fmt.Sprintf("%s%s%s", api, "apps/", name + "/")); err != nil {
		t.Fatal(err)
	}
	assertHeaders(resp, t)

	app = make(map[string] interface{})
	err = json.Unmarshal(body, &app)

	if _, ok = app["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}
}
