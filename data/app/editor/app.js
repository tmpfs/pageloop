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
  constructor (loc) {
    this.url = `/api/${loc.container}/${loc.application}/`
    this.loc = loc

    // current application
    this.app = {
      url: '',
      identifier: '',
      owner: loc.container,
      pages: [],
      files: []
    }

    this.preview = {
      url: '',
      path: ''
    }
  }

  json (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  getApplication () {

    console.log('getting app data: ' + this.url)

    return this.json(this.url)
      .then((app) => {
        // merge properties
        for (let k in app) {
          this.app[k] = app[k]
        }
        this.app.identifier = this.app.owner + '/' + this.app.name
        console.log('id: ' + this.app.identifier)
        return this.app
      })
  }

  getPages () {
    let url = this.url + 'pages/'
    return this.json(url)
      .then((list) => {
        this.app.pages = list
        return list
      })
  }

  getFiles () {
    let url = this.url + 'files/'
    return this.json(url)
      .then((list) => {
        this.app.files = list
        return list
      })
  }
}

class EditorApplication {
  constructor (loc) {
    this.loc = loc
    this.data = new AppDataSource(this.loc)
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
    this.log(`Loading preview ${this.getPreviewUrl()}`)
    this.preview.path = this.data.app.url
    this.preview.url = this.getPreviewUrl()
  }

  ui () {
    let data = this.data
    let bus = new Vue()

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
          name: data.app.identifier,
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
          pages: data.app.pages,
          files: data.app.files,
          components: [],
          currentView: ''
        }
      },
      methods: {
        loadPages: function () {
          return data.getPages()
            .then((list) => {
              this.pages = list
              bus.$emit('pages:load', list)
            })
        },
        loadFiles: function (url) {
          return data.getFiles()
            .then((list) => {
              this.files = list
              console.log('FILES LOADED')
              bus.$emit('files:load', list)
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
            bus.$on('pages:load', (list) => {
              this.list = list
            })
          },
          methods: {
            click: function (item) {
              bus.$emit('open:file', item)
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
            bus.$on('files:load', (list) => {
              this.list = list
            })
          },
          methods: {
            click: function (item) {
              bus.$emit('open:file', item)
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
          created: function () {}
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
          currentView: 'welcome',
          defaultOpenView: 'source-editor',
          currentFile: null
        }
      },
      created: function () {
        bus.$on('open:file', (item) => {
          this.open(item)
        })
      },
      methods: {
        open: function (item) {
          console.log('open in editor')
          console.log(item)

          if (this.currentView === 'welcome') {
            this.currentView = this.defaultOpenView
          }

          this.currentFile = item

          const mime = item.mime.replace(/;.*$/, '')

          switch (mime) {
            case 'text/html':
              console.log('got html mime type')
              break
          }
        }
      },
      components: {
        welcome: {
          template: `<p>Select a page or file to start editing.</p>`
        },
        'source-editor': {
          template: `<div class="source-editor"></div>`,
          mounted: function () {
            console.log('mounted')
            console.log(this.$parent.currentFile)

            /*
            let value = ''
            let item = this.$parent.currentFile
            if (item) {
              value = `Loading ${item.url}`
            }

            let url = `/apps/source/${}`
            */

            /*
            CodeMirror(document.querySelector('.source-editor'), {
              value: value,
              mode: 'htmlmixed',
              theme: 'midnight'
            })
            */
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
      msg = '! ' + msg
      this.logger.error = true
    } else {
      msg = '# ' + msg
      this.logger.error = false
    }
    this.logger.message = msg
  }

  load () {
    let data = this.data
    this.log(`Loading app from ${data.url}`)
    return this.data.getApplication()
      .then((app) => {
        this.identifier.name = app.identifier
        this.log(`Loading pages for ${this.data.app.name}`)
      })
      .then(this.sidebar.loadPages())
      .then(this.sidebar.loadFiles())
      .then(() => {
        this.refresh()
      })
      .catch((err) => { this.log(err) })
  }

  init (loc) {
    loc = loc || this.loc
    this.ui()
    this.load()
      .then(() => {
        this.sidebar.currentView = 'pages'
        this.log('Done')

        //
        console.log(this.data)
      })
  }
}

let app = new EditorApplication(new LocationParser().parse())
app.init()
