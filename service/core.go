package service

type Core struct {}
type MetaArgs struct {}
type MetaReply struct {}

type FooError struct {}

func (f *FooError) Error() string {
  return "Foo error"
}

func (c *Core) Meta(args MetaArgs, reply *MetaReply) error {
  return nil
}

type fooReply struct{}

func (c *Core) Foo (args MetaArgs, reply *MetaReply) *FooError {
  return nil
}
