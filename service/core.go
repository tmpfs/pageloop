package service

import(
  //"fmt"
  "expvar"
)

// Type for service methods that do not accept any arguments.
type VoidArgs struct {}

type Core struct {
  Name string
  Version string
}

type MetaReply struct {
  Name string `json:"name"`
  Version string `json:"version"`
}

type StatsReply struct {
  Uptime map[string]interface{} `json:"uptime"`
  Http map[string]interface{} `json:"http"`
  Ws map[string]interface{} `json:"ws"`
}

// Meta information (/).
func (s *Core) Meta(argv VoidArgs, reply *MetaReply) error {
  reply.Name = s.Name
  reply.Version = s.Version
  return nil
}

// Stats information (/stats)
func (c *Core) Stats(argv VoidArgs, reply *StatsReply) error {
  expvar.Do(func(kv expvar.KeyValue) {

    // Ignore built in exposed variables
    if kv.Key == "cmdline" || kv.Key == "memstats" {
      return
    }

    // TODO: handle strings and floats

    //fmt.Printf("%#v\n", kv.Key)
    //fmt.Printf("%#v\n", kv.Value)

    var values map[string]interface{}

    // Handle maps
    if hashmap, ok := kv.Value.(*expvar.Map); ok {
      values = make(map[string]interface{})
      hashmap.Do(func(mkv expvar.KeyValue) {
        if i, ok := mkv.Value.(*expvar.Int); ok {
          values[mkv.Key] = i.Value()
        }
      })
    } else {
      // Handle functions
      if fn, ok := kv.Value.(expvar.Func); ok {
        res := fn.Value()
        if mapped, ok := res.(map[string]interface{}); ok {
          values = mapped
        }
      }
    }

    switch (kv.Key) {
      case "uptime":
        reply.Uptime = values
      case "http":
        reply.Http = values
      case "ws":
        reply.Ws = values
    }
  })
  return nil
}
