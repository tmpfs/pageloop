// Package service exposes methods that may be used over the network.
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
