package model

import (
	"fmt"
  "errors"
  "regexp"
)

var(
	NamePattern string = `^[a-zA-Z0-9]+[-a-zA-Z0-9]*$`
	NamePatternRe = regexp.MustCompile(NamePattern)
)

func ValidName(name string) bool {
  return NamePatternRe.MatchString(name)
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

// Contains a slice of applications.
type Container struct {
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Apps []*Application `json:"apps"`

	// A protected container makes all it's applications protected
	Protected bool `json:"protected,omitempty"`
}

// Create a new container.
func NewContainer(name string, description string, protected bool) *Container {
	return &Container{Name: name, Description: description, Protected: protected}
}

// Add an application to the container, the application must
// have the Name field set and it must not already exist in
// the container list.
//
// Application names may only contain lowercase, uppercase, hyphens
// and digits. They may not begin with a hyphen.
func (c *Container) Add(app *Application) error {
	if app.Name == "" {
		return errors.New("Application name is required to add to container")
	}

	if !NamePatternRe.MatchString(app.Name) {
		return fmt.Errorf("Application name %s must match pattern %s", app.Name, NamePattern)
	}

	var exists *Application = c.GetByName(app.Name)
	if exists != nil {
		return fmt.Errorf("Application exists with name %s", app.Name)
	}

	app.Protected = c.Protected
	app.Container = c
  app.ContainerName = c.Name

	c.Apps = append(c.Apps, app)
	return nil
}

// Remove an application from the container.
func (c *Container) Del(app *Application) error {
	for i, a := range c.Apps {
		if app == a {
      before := c.Apps[0:i]
			after := c.Apps[i+1:]
			c.Apps = append(before, after...)
		}
	}
	return nil
}

// Get an application by name.
func (c *Container) GetByName(name string) *Application {
	for _, app := range c.Apps {
		if app.Name == name {
			return app
		}
	}
	return nil
}

