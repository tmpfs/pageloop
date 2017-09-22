package core

import(
  "expvar"
  "time"
  "encoding/json"
)

var(
  Stats *Statistics
)

type Uptime struct {
  Hours float64 `json:"hours"`
  Minutes float64 `json:"minutes"`
  Seconds float64 `json:"seconds"`
  String string `json:"value"`
}

type Statistics struct {
  // Time the statistics started
  StartTime time.Time `json:"-"`
  // Uptime calculation
  Uptime *Uptime `json:"uptime"`
  // HTTP server statistics
  Http *expvar.Map `json:"http"`
  // RPC service statistics
  Rpc *expvar.Map `json:"rpc"`
  // Websocket client connections and stats
  Websocket *expvar.Map `json:"ws"`
}

// Update the uptime and assign it to the statistics.
func (s *Statistics) Now() *Statistics {
  s.Uptime = uptime().(*Uptime)
  return s
}

func (s *Statistics) MarshalJSON() ([]byte, error) {
  o := make(map[string]interface{})
  o["uptime"] = uptime()
  o["http"] = mapToInterface(s.Http)
  o["ws"] = mapToInterface(s.Websocket)
  o["rpc"] = mapToInterface(s.Rpc)
  return json.Marshal(&o)
}

// Private

func mapToInterface (hashmap *expvar.Map) map[string]interface{} {
  vals := make(map[string]interface{})
  hashmap.Do(func(mkv expvar.KeyValue) {
    if i, ok := mkv.Value.(*expvar.Int); ok {
      vals[mkv.Key] = i.Value()
    }
  })
  return vals
}

func uptime () interface{} {
  duration := time.Since(Stats.StartTime)
  up := &Uptime{
    Hours: duration.Hours(),
    Minutes: duration.Minutes(),
    Seconds: duration.Seconds(),
    String: duration.String()}
  return up
}

func init() {
  Stats = &Statistics{StartTime: time.Now()}
  Stats.Http = expvar.NewMap("http")
  Stats.Http.Add("requests", 0)
  Stats.Http.Add("responses", 0)
  Stats.Http.Add("body-in", 0)
  Stats.Http.Add("body-out", 0)

  Stats.Websocket = expvar.NewMap("ws")
  Stats.Websocket.Add("connections", 0)

  Stats.Rpc = expvar.NewMap("rpc")
  Stats.Rpc.Add("calls", 0)
  Stats.Rpc.Add("errors", 0)

  expvar.Publish("uptime", expvar.Func(uptime))
}
