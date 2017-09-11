package adapter

// @deprecated - all deprecated from v1

import (
  //"fmt"
  //"net/http"
  . "github.com/tmpfs/pageloop/model"
  //. "github.com/tmpfs/pageloop/util"
)

func (b *CommandAdapter) ListApplications(c *Container) []*Application {
  return c.Apps
}
