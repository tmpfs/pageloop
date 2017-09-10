package handler

import (
  //"log"
	"fmt"
	"errors"
	"net/http"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
  . "github.com/tmpfs/pageloop/core"
	. "github.com/tmpfs/pageloop/model"
)

// TODO: operate using the command adapter!!!

func RpcHandler(mux *http.ServeMux, host *Host) http.Handler {
	// RPC endpoint
	endpoint := rpc.NewServer()
	// Do not specify charset on MIME type here
	endpoint.RegisterCodec(json.NewCodec(), "application/json")

	hostService := new(HostService)
	hostService.Host = host

	endpoint.RegisterService(hostService, "host")

	app := new(AppService)
  app.Host = host

	endpoint.RegisterService(app, "app")

	mux.Handle(RPC_URL, endpoint)

	return endpoint
}

type HostService struct {
  Host *Host
}

type HostListArgs struct {
	Index int `json:"index"`
	Len int `json:="length"`
}

type HostListReply struct {
	Containers []*Container `json:"containers"`
}

// Get a slice of the host container list.
//
// If length is zero it is set to the number of applications so
// pass index zero without a length to list all applications.
func (h *HostService) List(r *http.Request, args *HostListArgs, reply *HostListReply) error {
	var host *Host = h.Host
	if args.Len == 0 {
		args.Len = len(host.Containers) - args.Index
	}
	reply.Containers = host.Containers[args.Index:args.Len]

	return nil
}

type AppListArgs struct {
	GroupId string `json:"gid"`
	Index int `json:"index"`
	Len int `json:="length"`
}

type AppListReply struct {
	Apps []*Application `json:"apps"`
}

type AppService struct {
  Host *Host
}

// Get a slice of the application list for a container.
//
// If length is zero it is set to the number of applications so
// pass index zero without a length to list all applications.
func (h *AppService) List(r *http.Request, args *AppListArgs, reply *AppListReply) error {
	var container *Container
	if container = h.Host.GetByName(args.GroupId); container == nil {
		return errors.New(fmt.Sprintf("No container found for %s", args.GroupId))
	}
	if args.Len == 0 {
		args.Len = len(container.Apps) - args.Index
	}
	reply.Apps = container.Apps[args.Index:args.Len]
	return nil
}

type AppGetArgs struct {
	GroupId string `json:"gid"`
	Name string `json:"name"`
}

type AppGetReply struct {
	App *Application `json:"app"`
}

// Get an application by name.
func (h *AppService) Get(r *http.Request, args *AppGetArgs, reply *AppGetReply) error {
	var container *Container
	if container = h.Host.GetByName(args.GroupId); container == nil {
		return errors.New(fmt.Sprintf("No container found for %s", args.GroupId))
	}
	reply.App = container.GetByName(args.Name)
	return nil
}
