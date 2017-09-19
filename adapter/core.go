package adapter

type Core struct {}
type MetaArgs struct {}
type MetaReply struct {}

func (c *Core) Meta(args MetaArgs, reply *MetaReply) error {
  return nil
}
