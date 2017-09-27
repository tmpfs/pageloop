// Package service exposes methods that may be used over the network.
//
// Service methods must declare two arguments both of which are of pointer
// type and the first argument must be a struct. They should return either
// error or *StatusError.
//
// Currently the fields of the method argument must be simple types and slices
// nested structs are not currently supported.
package service

import(
  //"fmt"
  . "github.com/tmpfs/pageloop/core"
)

// Type for service methods that do not accept any arguments.
type VoidArgs struct {}

// CoreService service.
type CoreService struct {}

// Meta information (/).
func (s *CoreService) Meta(argv *VoidArgs, reply *MetaInfo) error {
  reply.Name = MetaData.Name
  reply.Version = MetaData.Version
  return nil
}

// Stats information (/stats)
func (c *CoreService) Stats(argv *VoidArgs, reply *Statistics) error {
  // Update uptime stats
  Stats.Now()

  // Assign statistics fields for reply
  reply.Uptime = Stats.Uptime
  reply.Http = Stats.Http
  reply.Rpc = Stats.Rpc
  reply.Websocket = Stats.Websocket

  return nil
}
