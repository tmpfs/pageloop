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
      switcher: this.doc.querySelector('.switcher'),
      previewUrl: this.doc.querySelector('.preview-url'),
      sidebar: this.doc.querySelector('.sidebar'),
      pagesList: this.doc.querySelector('.pages-list'),
      filesList: this.doc.querySelector('.files-list'),
      componentsList: this.doc.querySelector('.components-list'),
      editor: this.doc.querySelector('.editor'),
      preview: this.doc.querySelector('.preview'),
      live: this.doc.querySelector('.live'),
      log: this.doc.querySelector('.log')
    }
  }

  template (selector) {
    var tpl = document.querySelector('template')
    tpl = tpl && tpl.content ? tpl.content : tpl
    tpl = tpl.querySelector(selector)
    tpl = tpl.cloneNode(true)
    return tpl
  }

  text (txt) {
    return this.doc.createTextNode(txt)
  }

  element (name, attr) {
    let el = this.doc.createElement(name)
    for (let k in attr) {
      el.setAttribute(k, attr[k])
    }
    return el
  }

  removeChildren (node) {
    node.innerHTML = ''
  }

  get (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  getPreviewUrl () {
    return document.location.origin + this.app.url
  }

  refresh () {
    let u = this.getPreviewUrl()
    this.el.previewUrl.innerText = `${this.app.url}`
    this.el.previewUrl.setAttribute('href', u)
    this.el.live.addEventListener('load', () => {
      this.log(`Preview loaded ${u}`)
    })
    this.el.live.setAttribute('src', u)
  }

  log (msg) {
    let p = this.element('p')
    let err = (msg instanceof Error)
    if (err) {
      p.setAttribute('class', 'error')
      msg = '! ' + err
    } else {
      p.setAttribute('class', 'info')
      msg = '# ' + msg
    }
    p.appendChild(this.text(msg))
    this.removeChildren(this.el.log)
    this.el.log.appendChild(p)
  }

  render (data, fn) {
    let out = ''
    data.forEach((item) => {
      out += fn(item)
    })
    return out
  }

  getPageTemplate (item) {
    return `
        <div class="page"><span class="name" data-file="${item.url}" title="Open ${item.name}">${item.url}</span></div>
      `
  }

  getFileTemplate (item) {
    return `
        <div class="file"><span class="name" data-file="${item.url}" title="Open ${item.name}">${item.url}</span></div>
      `
  }

  getSwitcherContainer (c) {
    let o = `<div class="container"><h3>${c.name}</h3><ul>`
    c.apps.forEach((a) => {
      let href = '/apps/edit/' + c.name + '/' + a.name + '/'
      o += `
          <li>
            <a href="${href}">${a.name}</a>
          </li>
        `
    })
    o += `</ul></div>`
    return o
  }

  ui () {
    // Sidebar tabs
    this.doc.querySelectorAll('.tab > a').forEach((n) => {
      console.log(n)
      n.addEventListener('click', (e) => {
        e.preventDefault()
        let el = e.currentTarget
        let target = this.doc.querySelector(el.getAttribute('data-target'))

        this.doc.querySelectorAll('.sidebar > .scroll > *').forEach((n) => {
          n.classList.add('hidden')
        })
        target.classList.remove('hidden')
        this.doc.querySelector('.tab > a.selected').classList.remove('selected')
        el.classList.add('selected')
      })
    })

    // Application switcher
    this.doc.querySelector('.app-id a').addEventListener('click', (e) => {
      e.preventDefault()
      let n = e.currentTarget
      !n.isOpen ? this.select() : this.el.switcher.classList.add('hidden')
      n.isOpen = !n.isOpen
    })
  }

  // Select an application.
  select () {
    let url = `/api/`
    this.log(`Loading container data from ${url}`)
    this.get(url)
      .then((containers) => {
        this.removeChildren(this.el.switcher)
        this.el.switcher.classList.remove('hidden')
        console.log(containers)
        containers.forEach((c) => {
          let cn = this.getSwitcherContainer(c)
          this.el.switcher.innerHTML += cn
        })
      })
      .catch((err) => { this.log(err) })
  }

  init () {
    let pth = this.doc.location.pathname
    let parts = pth.replace(/\/+$/, '').split('/')
    this.application = parts.pop()
    this.container = parts.pop()

    this.el.appid.querySelector('.name').innerText =
      `${this.container} / ${this.application}`

    let url = `/api/${this.container}/${this.application}/`
    this.log(`Loading app data from ${url}`)
    this.get(url)
      .then((app) => {
        this.app = app
        this.log(`Loading pages for ${this.app.name}`)
        return this.get(url + 'pages/')
          .then((pages) => {
            this.pages = pages
            this.el.pagesList.innerHTML = this.render(this.pages, this.getPageTemplate)
            this.log(`Loading files for ${this.app.name}`)
            return this.get(url + 'files/')
              .then((files) => {
                this.files = files
                this.el.filesList.innerHTML = this.render(this.files, this.getFileTemplate)
              })
          })
      })
      .then(() => {
        this.log(`Loading preview ${this.getPreviewUrl()}`)

        // Load the iframe preview
        this.refresh()

        // Initialize the UI event handling
        this.ui()
      })
      .catch((err) => { this.log(err) })
  }
}

let app = new EditorApplication()
app.init()
