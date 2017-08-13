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
	mux.Handle("/rpc/", endpoint)

	return service
}
