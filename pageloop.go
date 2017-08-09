// System for hyperfast HTML document editing.
//
// Stores HTML documents on the server as in-memory DOM 
// documents that can be modified on the client. The client 
// provides an editor view and a preview of the rendered 
// page loaded in an iframe.
package pageloop

import(
  . "github.com/tmpfs/pageloop/model"
)

func LoadApplication(path string, app *Application) *Application {
  
}

// Load an application using the given loader implementation, 
// if a nil loader is given the default file system loader is used.
func (app *Application) Load(path string, loader ApplicationLoader) Application {
  if loader == nil {
    loader = FileSystemLoader{}
  }
  loader.LoadApplication(path, app)
  app.Urls = make(map[string] File)
  app.SetComputedFields(path)
  app.Merge()
  return *app
}

