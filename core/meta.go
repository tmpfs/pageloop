package core

// Application meta data stuff

const(
  Name string = "pageloop"
  Version string = "1.0"
)

var(
  MetaData *MetaInfo = &MetaInfo{Name: Name, Version: Version}
)

type MetaInfo struct {
  Name string `json:"name"`
  Version string `json:"version"`
}
