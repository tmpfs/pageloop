package core

import(
  "expvar"
  "time"
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
  StartTime time.Time
  Http *expvar.Map
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
  Stats.Http.Add("bytes-in", 0)
  Stats.Http.Add("bytes-out", 0)

  expvar.Publish("uptime", expvar.Func(uptime))
}
