/* globals document fetch */
class EditorApplication {
  constructor () {
    this.container = null
    this.application = null

    this.app = null
    this.doc = document
    this.files = null
    this.pages = null
    this.el = {
      appid: this.doc.querySelector('.app-id'),
      sidebar: this.doc.querySelector('.sidebar'),
      editor: this.doc.querySelector('.editor'),
      preview: this.doc.querySelector('.preview'),
      live: this.doc.querySelector('.live')
    }
  }

  text (txt) {
    return this.doc.createTextNode(txt)
  }

  get (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  init () {
    let pth = this.doc.location.pathname
    let parts = pth.replace(/\/+$/, '').split('/')
    this.application = parts.pop()
    this.container = parts.pop()

    this.el.appid.appendChild(this.text(`${this.container} / ${this.application}`))

    let url = `/api/${this.container}/${this.application}/`
    this.get(url)
      .then((app) => {
        this.app = app
        console.log(app)
        return this.get(url + 'pages/')
          .then((pages) => {
            console.log(pages)
            this.pages = pages
            return this.get(url + 'files/')
              .then((files) => {
                console.log(files)
                this.files = files
              })
          })
      })
      .then(() => {
        console.log('all init done')
      })
      .catch((err) => { throw err })
  }
}

let app = new EditorApplication()
app.init()
