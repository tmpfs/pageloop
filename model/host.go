package model

import(
  "fmt"
)

// References an existing mounted application (and optionally specific file)
// used for the intialization of application files from templates.
type ApplicationTemplate struct {
	Container string
	Application string
	File string
}

// Contains a slice of containers.
type Host struct {
	Containers []*Container `json:"containers"`
}

// Create a new host.
func NewHost() *Host {
	h := new(Host)
	//h.Containers = make([] *Container)
	return h
}

// Add a container.
func (h *Host) Add(c *Container) {
	h.Containers = append(h.Containers, c)
}

// Get a container by name.
func (h *Host) GetByName(name string) *Container {
	for _, container := range h.Containers {
		if container.Name == name {
			return container
		}
	}
	return nil
}

// Find an application from an application template reference.
func (h *Host) LookupTemplate(t *ApplicationTemplate) (*Application, error) {
  container := h.GetByName(t.Container)
  if container == nil {
    return nil, fmt.Errorf("Template container %s does not exist", t.Container)
  }
  app := container.GetByName(t.Application)
  if app == nil {
    return nil, fmt.Errorf("Template application %s does not exist", t.Application)
  }

  return app, nil
}

// Find an application file from an application template reference.
func (h *Host) LookupTemplateFile(t *ApplicationTemplate) (*File, error) {
  var err error
  var app *Application
  if app, err = h.LookupTemplate(t); err != nil {
    return nil, err
  }
  url := t.File
  return app.Urls[url], nil
}

// Generate HTML markup for a directory listing.
func (h *Host) DirectoryListing (file *File) ([]byte, error) {
  // Build the template data
  d := file.DirectoryListing()
  // Get the directory listing template file
  c := h.GetByName("template")
  a :=  c.GetByName("listing")
  f := a.Urls["/index.html"]
  p := f.Page()
  // Parse and execute the template
  if tpl, err := p.ParseTemplate(file.Path, f.Source(false), p.DefaultFuncMap(), false); err != nil {
    return nil, err
  } else {
    if output, err := p.ExecuteTemplate(tpl, d); err != nil {
      return nil, err
    } else {
      return output, nil
    }
  }
  return nil, nil
}
