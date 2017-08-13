package pageloop

import (
  //"log"
	//"errors"
	"net/http"
  //"path/filepath"
	//"regexp"
  //"time"
	"github.com/tmpfs/pageloop/model"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
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

	app := new(AppService)
	app.Root = root

	endpoint.RegisterService(app, "app")

	mux.Handle("/rpc/", endpoint)

	return service
}

type AppService struct {
	Root *PageLoop
}

type AppListArgs struct {
	Index int `json:"index"`
	Len int `json:="length"`
}

type AppListReply struct {
	Apps []*model.Application `json:"apps"`
}

// Get a slice of the application list for a container.
//
// If length is zero it is set to the number of applications so 
// pass index zero without a length to list all applications.
func (h *AppService) List(r *http.Request, args *AppListArgs, reply *AppListReply) error {
	var container *model.Container = h.Root.Container
	if args.Len == 0 {
		args.Len = len(container.Apps) - args.Index
	}
	reply.Apps = container.Apps[args.Index:args.Len]
	return nil
}

type AppGetArgs struct {
	Name string `json:"name"`
}

type AppGetReply struct {
	App *model.Application `json:"app"`
}

// Get an application by name.
func (h *AppService) Get(r *http.Request, args *AppGetArgs, reply *AppGetReply) error {
	var container *model.Container = h.Root.Container
	reply.App = container.GetByName(args.Name)
	return nil
}
