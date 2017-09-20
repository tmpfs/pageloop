package service

import(
  //"fmt"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type AppService struct {

}

type AppReference struct {
  Container string
  Application string
}

func (s *AppService) Read(argv AppReference, app *Application) *StatusError {
  return nil
}
