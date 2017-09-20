package service

type Core struct {}
type MetaArgs struct {}
type MetaReply struct {}

func (c *Core) Meta(args MetaArgs, reply *MetaReply) error {
  return nil
}

/*
type fooReply struct{}

func (c *Core) Foo (args MetaArgs, reply *MetaReply) string {
  return ""
}
*/
