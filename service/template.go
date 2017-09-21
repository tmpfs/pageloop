package service

import(
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

type TemplateService struct {
  Host *Host
}

// List all system templates and user applications
// that have been marked as a template.
func (s *TemplateService) List(argv *VoidArgs, reply *ServiceReply) *StatusError {
  // Get built in and user templates
  c := s.Host.GetByName("template")
  u := s.Host.GetByName("user")
  list := append(c.Apps, u.Apps...)
  var apps []*Application
  for _, app := range list {
    if app.IsTemplate {
      apps = append(apps, app)
    }
  }
  reply.Reply = apps
  return nil
}
