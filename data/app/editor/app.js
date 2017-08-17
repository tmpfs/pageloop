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

  setApplication (app) {
    this.data.app = app
    this.identifier.name = this.location.container + '/' + this.location.application
  }

  refresh () {
    this.preview.path = this.data.app.url
    this.preview.url = this.getPreviewUrl()
  }

  ui (data) {
    this.identifier = new Vue({
      template: `<div class="app-id"><a href="#" @click="click">▾ <span class="name">{{name}}</span></a></div>`,
      data: function () {
        return {
          name: data.app.id,
          show: false
        }
      },
      methods: {
        click: function (e) {
          e.preventDefault()
          this.show = !this.show
          console.log('id click: ' + this.show)
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
        <p class="log" v-bind:class="{error: error}">{{message}}</p>
      `,
      data: function () {
        return {
          message: '',
          error: false
        }
      }
    })

    let header = new Vue({el: 'header'})
    let main = new Vue({el: 'main', data: data})
    let footer = new Vue({el: 'footer', data: data})

    // mount views
    this.preview.$mount('.preview')
    this.logger.$mount('footer .log')
    this.identifier.$mount('.app-id')

    return {header: header, main: main, footer: footer}
  }

  log (msg) {
    let err = (msg instanceof Error)
    if (err) {
      // TODO: add error class
      msg = '! ' + msg
      this.logger.error = true
    } else {
      msg = '# ' + msg
      this.logger.error = false
    }

    this.logger.message = msg
  }

  load (data) {
    let url = `/api/${this.location.container}/${this.location.application}/`
    this.log(`Loading app data from ${url}`)
    return this.get(url)
      .then((app) => {
        this.setApplication(app)
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
      .then(() => {
        this.log('Done')
      })
  }
}

let app = new EditorApplication()
app.init()

