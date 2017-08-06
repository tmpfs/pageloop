package vdom
  
import(
  "log"
  "bytes"
  "strconv"
  "golang.org/x/net/html"
)

type Vdom struct {
  Document *html.Node
  Map map[string] *html.Node
}

type Settings struct {
  IdAttribute string 
}

func GetSettings() Settings {
  return Settings{IdAttribute: "data-id"}
}

func Parse(b []byte, settings Settings) (*Vdom, error) {
  r := bytes.NewBuffer(b)
  doc, err := html.Parse(r)
  if err != nil {
    return nil, err
  }

  dom := Vdom{Document: doc, Map: make(map[string] *html.Node)}
  var f func(n *html.Node, ids []int)
  f = func(n *html.Node, ids []int) {
    var id string
    var i int = 0
    log.Println("ids length", len(ids))
    log.Println(ids)
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      if c.Type == html.ElementNode {
        list := append(ids, i)
        id = ""
        log.Println(c.Data)
        for index, num := range list {
          log.Println(string(num))
          if index > 0 {
            id += "."
          }
          id += strconv.Itoa(num)
        }
        log.Printf("id: %s", id)
        dom.Map[id] = c
        attr := html.Attribute{Key: settings.IdAttribute, Val: id}
        c.Attr = append(c.Attr, attr)
        log.Println(c.Attr)
        f(c, list)
        i++
      }
    }
  }

  var ids []int
  f(doc, ids)
  return &dom, nil
}
