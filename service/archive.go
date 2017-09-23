package service

import(
  // "fmt"
  //"net/http"
  //. "github.com/tmpfs/pageloop/core"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type ArchiveService struct {
  // Reference to the host
  Host *Host
}



// Export a zip archive of application files.
func (s *ArchiveService) Export(app *Application, reply *ServiceReply) *StatusError {
  if _, a, err := LookupApplication(s.Host, app); err != nil {
    return err
  } else {
    println("Implement zip archive logic.")
    reply.Reply = a
  }
  return nil
}
