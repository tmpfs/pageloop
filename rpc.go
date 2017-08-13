package pageloop

import (
  //"log"
	//"errors"
	"net/http"
  //"path/filepath"
	//"regexp"
  //"time"
  //"github.com/tmpfs/pageloop/model"
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

	endpoint.RegisterService(new(HelloService), "hello")

	mux.Handle("/rpc/", endpoint)

	return service
}
type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

type HelloService struct {}

func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	reply.Message = "Hello, " + args.Who + "!"
	return nil
}
