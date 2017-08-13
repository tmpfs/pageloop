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
	endpoint.RegisterCodec(json.NewCodec(), JSON_MIME)

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

func (h *AppService) List(r *http.Request, args *AppListArgs, reply *AppListReply) error {
	if args.Len == 0 {
		args.Len = len(h.Root.Container.Apps) - args.Index
	}
	reply.Apps = h.Root.Container.Apps[args.Index:args.Len]
	return nil
}
