/* globals Vue CodeMirror document fetch */

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
      url: '',
      pages: null,
      files: null
    }

    this.preview = {
      url: '',
      path: ''
    }
  }

  get containers () {
    return this._containers
  }
}

class EditorApplication {
  constructor (loc) {
    this.loc = loc
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
    // merge properties
    for (let k in app) {
      this.data.app[k] = app[k]
    }
    this.identifier.name = this.loc.container + '/' + app.name
  }

  refresh () {
    this.log(`Loading preview ${this.getPreviewUrl()}`)
    this.preview.path = this.data.app.url
    this.preview.url = this.getPreviewUrl()
  }

  ui (data) {
    let get = this.get

    let switcher = this.switcher = new Vue({
      template: `<div class="switcher" v-bind:class="{hidden: hidden}"></div>`,
      data: function () {
        return {
          hidden: true
        }
      }
    })

    this.identifier = new Vue({
      template: `
        <div class="app-id">
          <a href="#" @click="click">▾ <span class="name">{{name}}</span></a>
        </div>`,
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
          switcher.hidden = !this.show
        }
      }
    })

    this.sidebar = new Vue({
      template: `
        <div class="sidebar">
          <nav class="tabs">
            <a v-bind:class="{selected: currentView === 'pages'}"
              @click="currentView = 'pages'"
              href="#" title="Show pages">Pages</a>
            <a v-bind:class="{selected: currentView === 'files'}"
              @click="currentView = 'files'"
              href="#"  title="Show files">Files</a>
            <a v-bind:class="{selected: currentView === 'components'}"
              @click="currentView = 'components'" href="#"
              title="Show components">Components</a>
          </nav>
          <div class="scroll">
            <component v-bind:is="currentView"></component>
          </div>
        </div>
      `,
      data: function () {
        return {
          pages: [],
          files: [],
          components: [],
          currentView: ''
        }
      },
      methods: {
        loadPages: function (url) {
          return get(url + 'pages/')
            .then((list) => {
              data.app.pages = list
              this.pages = list
            })
        },
        loadFiles: function (url) {
          return get(url + 'files/')
            .then((list) => {
              data.app.files = list
              this.files = list
            })
        }
      },
      components: {
        pages: {
          template: `
            <div class="pages-list">
              <a @click="click(item)" class="page" v-for="item in list">
                <span class="name">{{item.url}}</span>
              </a>
            </div>`,
          data: function () {
            return {
              list: [],
              current: null
            }
          },
          created: function () {
            this.list = this.$parent.pages
          },
          methods: {
            click: function (item) {
              console.log('page clicked')
              console.log(item)
            }
          }
        },
        files: {
          template: `
            <div class="files-list">
              <a @click="click(item)" class="file" v-for="item in list">
                <span class="name">{{item.url}}</span>
              </a>
            </div>`,
          data: function () {
            return {
              list: []
            }
          },
          created: function () {
            this.list = this.$parent.files
          },
          methods: {
            click: function (item) {
              console.log('file clicked')
              console.log(item)
            }
          }
        },
        components: {
          template: `
            <div class="components-list">
            </div>`,
          data: function () {
            return {
              list: []
            }
          },
          created: function () {
            this.list = this.$parent.components
          }
        }
      }
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

    this.editor = new Vue({
      template: `
        <div class="editor">
          <div class="column-header">
            <h2>Editor</h2>
            <div class="column-options">
              <nav class="tabs">
                <a v-bind:class="{selected: currentView === 'source-editor'}"
                  @click="currentView = 'source-editor'"
                  href="#" title="Show source editor">Source</a>
                <a v-bind:class="{selected: currentView === 'visual-editor'}"
                  @click="currentView = 'visual-editor'"
                  href="#"  title="Show visual editor">Visual</a>
              </nav>
            </div>
          </div>
          <div class="scroll">
            <component v-bind:is="currentView"></component>
          </div>
        </div>
      `,
      data: function () {
        return {
          currentView: 'welcome'
        }
      },
      components: {
        welcome: {
          template: `<p>Select a page or file to start editing.</p>`
        },
        'source-editor': {
          template: `<div class="source-editor"></div>`,
          mounted: function () {
            console.log(document.querySelector('.source-editor'))
            CodeMirror(document.querySelector('.source-editor'), {
              value: 'function myScript(){return 100;}\n',
              mode: 'javascript'
            })
          }
        },
        'visual-editor': {
          template: `<div class="visual-editor"></div>`
        }
      }
    })

    Vue.component('editor-main', {
      template: `
        <div class="content-main">
          <div class="switcher hidden"></div>
          <div class="content">
            <div class="sidebar"></div>
            <div class="editor"></div>
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
    this.logger.$mount('footer .log')
    this.identifier.$mount('.app-id')
    this.switcher.$mount('.switcher')
    this.sidebar.$mount('.sidebar')
    this.editor.$mount('.editor')
    this.preview.$mount('.preview')

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

  load (loc, data) {
    let url = `/api/${loc.container}/${loc.application}/`
    this.log(`Loading app data from ${url}`)
    return this.get(url)
      .then((app) => {
        this.setApplication(app)
        this.log(`Loading pages for ${this.data.app.name}`)
        return this.sidebar.loadPages(url)
          .then(this.sidebar.loadFiles(url))
      })
      .then(() => {
        // Load the preview
        this.refresh()
      })
      .catch((err) => { this.log(err) })
  }

  init (loc) {
    loc = loc || this.loc
    this.ui(this.data)
    this.log('Interface created')
    this.load(loc, this.data)
      .then(() => {
        this.sidebar.currentView = 'pages'
        console.log(this.data)
        this.log('Done')
      })
  }
}

let app = new EditorApplication(new LocationParser().parse())
app.init()
