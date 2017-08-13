package pageloop

import (
  //"log"
	"fmt"
	"errors"
	"net/http"
	"github.com/tmpfs/pageloop/model"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

const(
	RPC_URL = "/rpc/"
)

type RpcService struct {
	Root *PageLoop
}

func NewRpcService(root *PageLoop, mux *http.ServeMux) *RpcService {
	var service *RpcService = &RpcService{Root: root}

	// RPC endpoint
	endpoint := rpc.NewServer()
	// Do not specify charset on MIME type here
	endpoint.RegisterCodec(json.NewCodec(), "application/json")

	host := new(HostService)
	host.Root = root

	endpoint.RegisterService(host, "host")

	app := new(AppService)
	app.Root = root

	endpoint.RegisterService(app, "app")


	mux.Handle(RPC_URL, endpoint)

	return service
}

type HostService struct {
	Root *PageLoop
}

type HostListArgs struct {
	Index int `json:"index"`
	Len int `json:="length"`
}

type HostListReply struct {
	Hosts map[string] *model.Container `json:"containers"`
}

// Get a slice of the host container list.
//
// If length is zero it is set to the number of applications so 
// pass index zero without a length to list all applications.
func (h *HostService) List(r *http.Request, args *HostListArgs, reply *HostListReply) error {
	reply.Hosts = h.Root.Host.Containers
	return nil
}

type AppListArgs struct {
	GroupId string `json:"gid"`
	Index int `json:"index"`
	Len int `json:="length"`
}

type AppListReply struct {
	Apps []*model.Application `json:"apps"`
}

type AppService struct {
	Root *PageLoop
}

// Get a slice of the application list for a container.
//
// If length is zero it is set to the number of applications so 
// pass index zero without a length to list all applications.
func (h *AppService) List(r *http.Request, args *AppListArgs, reply *AppListReply) error {
	var container *model.Container
	if container = h.Root.Host.Get(args.GroupId); container == nil {
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
	App *model.Application `json:"app"`
}

// Get an application by name.
func (h *AppService) Get(r *http.Request, args *AppGetArgs, reply *AppGetReply) error {
	var container *model.Container
	if container = h.Root.Host.Get(args.GroupId); container == nil {
		return errors.New(fmt.Sprintf("No container found for %s", args.GroupId))
	}
	reply.App = container.GetByName(args.Name)
	return nil
}
