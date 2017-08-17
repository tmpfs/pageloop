/* globals Vue document fetch */

class LocationParser {
  parse () {
    let pth = document.location.pathname
    let parts = pth.replace(/\/+$/, '').split('/')
    let application = parts.pop()
    let container = parts.pop()
    return {application: application, container: container}
  }
}

class AppDataSource {
  constructor (location) {
    this._location = location
    this._containers = null

    // current application
    this.app = {
      url: ''
    }

    this.preview = {
      url: '',
      path: ''
    }

    // application pages
    this.pages = null

    // application files
    this.files = null
  }

  get containers () {
    return this._containers
  }
}

class EditorApplication {
  constructor () {
    this.location = new LocationParser().parse()
    this.data = new AppDataSource(this.location)
  }

  get (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  getPreviewUrl () {
    return document.location.origin + this.data.app.url
  }

  refresh () {
    this.preview.path = this.data.app.url
    this.preview.url = this.getPreviewUrl()
    /*
    let u = this.getPreviewUrl()
    this.el.previewUrl.innerText = `${this.data.app.url}`
    this.el.previewUrl.setAttribute('href', u)
    this.el.live.addEventListener('load', () => {
      this.log(`Preview loaded ${u}`)
    })
    this.el.live.setAttribute('src', u)
    */
  }

  ui (data) {
    Vue.component('app-id', {
      template: `<div class="app-id"><a href="#switch">▾ <span class="name">{{name}}</span></a></div>`,
      data: function () {
        return {
          name: 'container / application'
        }
      }
    })

    Vue.component('app-sidebar', {
      template: `
        <div class="sidebar">
          <h2 class="tab">
            <a class="pages selected" data-target=".pages-list" href="#pages" title="Show pages">Pages</a>
            <a class="files" href="#files"  data-target=".files-list"title="Show files">Files</a>
            <a class="components" href="#components" data-target=".components-list" title="Show components">Components</a>
          </h2>
          <div class="scroll">
              <div class="pages-list"></div>
              <div class="files-list hidden"></div>
              <div class="components-list hidden"></div>
          </div>
        </div>
      `
    })

    this.preview = new Vue({
      template: `
        <div class="preview">
          <h2>Live Preview ~ <a class="preview-url" :href="url" title="Preview URL">{{path}}</a></h2>
          <iframe :src="url" class="live"></iframe>
        </div>
      `,
      data: function () {
        return {
          path: data.preview.path,
          url: data.preview.url
        }
      }
    })

    Vue.component('app-editor', {
      template: `
        <div class="editor">
          <h2>Editor</h2>
          <div class="scroll">
            <p>Select a page or file to start editing.</p>
          </div>
        </div>
      `
    })

    Vue.component('editor-main', {
      template: `
        <div class="content-main">
          <div class="content">
            <app-sidebar></app-sidebar>
            <app-editor></app-editor>
            <div class="preview"></div>
          </div>
        </div>
      `
    })

    this.logger = new Vue({
      template: `
        <div class="log"><p>{{message}}</p></div>
      `,
      data: function () {
        return {
          message: ''
        }
      }
    })

    let header = new Vue({el: 'header'})
    let main = new Vue({el: 'main', data: data})
    let footer = new Vue({el: 'footer', data: data})

    // mount vues
    this.preview.$mount('.preview')
    this.logger.$mount('footer .log')

    return {header: header, main: main, footer: footer}
  }

  log (msg) {
    let err = (msg instanceof Error)
    if (err) {
      // p.setAttribute('class', 'error')
      msg = '! ' + msg
    } else {
      // p.setAttribute('class', 'info')
      msg = '# ' + msg
    }

    this.logger.message = msg
  }

  load (data) {
    let url = `/api/${this.location.container}/${this.location.application}/`
    this.log(`Loading app data from ${url}`)
    this.get(url)
      .then((app) => {
        this.data.app = app
        this.log(`Loading pages for ${this.data.app.name}`)
        return this.get(url + 'pages/')
          .then((pages) => {
            this.data.pages = pages
            this.log(`Loading files for ${this.data.app.name}`)
            return this.get(url + 'files/')
              .then((files) => {
                this.data.files = files
              })
          })
      })
      .then(() => {
        this.log(`Loading preview ${this.getPreviewUrl()}`)

        // Load the iframe preview
        this.refresh()
      })
      .catch((err) => { this.log(err) })
  }

  init () {
    this.ui(this.data)
    this.log('Interface created')
    this.load(this.data)
  }
}

let app = new EditorApplication()
app.init()

