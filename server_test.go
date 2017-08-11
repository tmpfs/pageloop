package pageloop

import (
	//"net/http"
  "testing"
)

func Start(t *testing.T) *PageLoop {
	var err error
  var apps []Mountpoint
	apps = append(apps, Mountpoint{UrlPath: "/app/mock-app/", Path: "test/fixtures/mock-app"})
  loop := &PageLoop{}
	conf := ServerConfig{Mountpoints: apps, Addr: ":3577", Dev: true}
	if _, err = loop.NewServer(conf); err != nil {
		t.Error(err)
	}
	go loop.Listen()
	return loop
}

func TestServer(t *testing.T) {
	loop := Start(t)
	if loop == nil {
		t.Error("Failed to acquire pageloop entry")
	}
}
