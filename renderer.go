package blocks

import (
  "os"
  "log"
  "html/template"
)

func (app *Application) Render(page *Page) Application {
  data := page.file.data
  tpl := template.New(page.file.Relative)
  tpl, err := tpl.Parse(string(data))
  if err != nil {
    log.Fatal(err)
  }
  //log.Println(app.Pages[0].UserData)
  tpl.Execute(os.Stdout, page.UserData)

  return *app
}

