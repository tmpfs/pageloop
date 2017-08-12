package pageloop

import (
	//"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
  "testing"
)

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

func TestRestService(t *testing.T) {
	Start(t)

	// GET /api/
	resp, err := http.Get("http://localhost:3579/api/")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	var res map[string] interface{} = make(map[string] interface{})
	json.Unmarshal(body, &res)

	var apps []interface{}
	var app map [string] interface{}
	var name string
	var ok bool

	if apps, ok = res["apps"].([]interface{}); !ok {
		t.Error("Unexpected type for apps list")
	}

	//println(apps[0])

	if app, ok = apps[0].(map [string] interface{}); !ok {
		t.Error("Unexpected type for app")
	}

	if name , ok = app["name"].(string); !ok {
		t.Error("Unexpected type for name")
	}


	println(name)

	//println(ok)
	//fmt.Printf("%#v\n", res)
}
