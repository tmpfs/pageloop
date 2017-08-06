package blocks

import (
  "os"
  "log"
  "fmt"
  "html/template"
)

func (app *Application) Render(page *Page) Application {
  fmt.Println("--- render function ---")
  fmt.Println()

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

