/* globals Vue Vuex CodeMirror document fetch document history window */

class Router {
  constructor (href, strip) {
    this.defaultHref = href
    this.routes = []
    this.strip = strip
  }

  navigate (href, state) {
    let url = this.url(href)
    history.pushState({href: href, state: state}, '', url)
    this.route(href)
  }

  url (href) {
    return this.pathname + '#' + href
  }

  get pathname () {
    return document.location.pathname
  }

  get hash () {
    let h = document.location.hash.replace(/^#/, '')
    if (this.strip) {
      h = h.replace(/\/$/, '')
    }
    return h
  }

  replace (href, trigger) {
    document.location.replace(this.url(href))
    if (trigger) {
      this.route(href)
    }
  }

  add (ptn, map, fn) {
    if (typeof map === 'function') {
      fn = map
      map = null
    }
    this.routes.push({ptn: ptn, fn: fn, map: map})
  }

  route (href, state) {
    function result (href, route) {
      let o = {
        state: state,
        href: href,
        route: route,
        parts: [],
        map: {}
      }

      let parts = href.replace(/^\//, '').replace(/\/$/, '').split('/')
      o.parts = parts

      if (route.map) {
        route.map.forEach((val, i) => {
          o.map[val] = parts[i]
        })
      }

      return o
    }

    let r, ptn, fn, res
    for (let i = 0; i < this.routes.length; i++) {
      r = this.routes[i]
      ptn = r.ptn
      fn = r.fn
      if (typeof ptn === 'string' && href === ptn) {
        res = fn(result(href, r))
        if (res !== true) {
          break
        }
      } else if (ptn instanceof RegExp && ptn.test(href)) {
        res = fn(result(href, r))
        if (res !== true) {
          break
        }
      }
    }
  }

  start () {
    window.addEventListener('popstate', (e) => {
      if (e.state && e.state.href) {
        this.route(e.state.href, e.state.state)
      } else {
        this.route(this.hash)
      }
    })
    if (!this.hash) {
      if (this.defaultHref) {
        this.replace(this.defaultHref, true)
      }
    } else {
      this.route(this.hash)
    }
  }
}

class Log {
  constructor () {
    this.maximum = 1024
    this.messages = []
  }

  add (message) {
    this.messages.push(message)
    if (this.messages.length > this.maximum) {
      this.messages.shift()
    }
  }

  get last () {
    let m = null
    if (this.messages.length) {
      m = this.messages[this.messages.length - 1]
    }
    return m
  }

  toString () {
    let m = this.last
    if (m instanceof Error) {
      return m.message || ('' + m)
    }
  }
}

class AppDataSource {
  constructor () {
    this.api = '/api/'
    this.containers = []
    this.setApplication('', '')

    this.mainView = ''
    this.sidebarView = ''
    this.editorView = ''
    this.defaultEditorView = 'source-editor'

    this.defaultFile = {content: ''}

    this.previewUrl = ''

    // State for maximized column in edit mode
    this.maximizedColumn = ''

    this.log = new Log()

    this._flash = undefined
  }

  get flash () {
    let f = this._flash
    this._flash = undefined
    return f
  }

  set flash (msg) {
    this._flash = msg
  }

  getAppHref (...args) {
    let p = ['apps', this.container, this.application]
    args = args.map((val) => {
      return val.replace(/^\//, '')
    })
    p.push(...args)
    return p.join('/')
  }

  getIndexFile () {
    let files = this.app.files
    for (let i = 0; i < files.length; i++) {
      // got a published index page whether the source is
      // HTML or markdown
      if (files[i].uri === '/index.html') {
        return files[i]
      }
    }
  }

  setApplication (container, application) {
    this.container = container
    this.application = application
    this.url = `${this.api}${container}/${application}/`
    this.raw = `/apps/raw/${container}/${application}`

    // current application
    this.app = {
      url: '',
      identifier: '',
      owner: container,
      pages: [],
      files: [],
      // current selected file
      current: this.defaultFile
    }
  }

  get current () {
    return this.app.current
  }

  set current (file) {
    if (file) {
      let pages = this.app.pages
      for (let i = 0; i < pages.length; i++) {
        if (pages[i].url === file.url) {
          file.page = pages[i]
          break
        }
      }
    }
    this.app.current = file
  }

  hasApplication () {
    return this.container && this.application
  }

  hasFile () {
    return this.app.current.url !== undefined
  }

  getFile () {
    if (this.hasFile()) {
      return this.app.current
    }
    return null
  }

  isPage (file) {
    return file && file.page !== undefined
  }

  json (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  getContainers () {
    return this.json(this.api)
  }

  getApplication (loc) {
    if (!loc) {
      loc = this.loc
    }
    return this.json(this.url)
  }

  getPages () {
    let url = this.url + 'pages/'
    return this.json(url)
  }

  getFiles () {
    let url = this.url + 'files/'
    return this.json(url)
  }

  getFileContents (pathname) {
    let url = this.raw
    return fetch(url + pathname)
      .catch((err) => err)
  }

  deleteFile (file) {
    let url = this.url + 'files' + file.url
    let opts = {
      method: 'DELETE'
    }
    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  saveFile (file, value) {
    file.content = value

    let url = this.url + 'files' + file.url
    let opts = {
      method: 'POST',
      headers: {
        'Content-Type': file.mime
      },
      body: value
    }
    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  createNewFile (name, template) {
    let url = this.url + 'files' + name
    let opts = {
      method: 'PUT',
      headers: {
        'Content-Type': template
      }
    }

    // Create empty file
    if (template === '') {
      opts.headers['Content-Length'] = 0
      opts.body = ''
    }
    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }
}

class EditorApplication {
  constructor () {
    this.bus = new Vue()
    let data = this.data = new AppDataSource()
    let store = this.store = new Vuex.Store({
      state: this.data,
      mutations: {
        flash: function (state, message) {
          state.flash = message
        },
        log: function (state, message) {
          state.log.add(message)
        },
        containers (state, list) {
          state.containers = list
        },
        app (state, app) {
          // merge properties
          for (let k in app) {
            state.app[k] = app[k]
          }
          state.app.identifier = state.app.owner + ' / ' + state.app.name
        },
        files: function (state, list) {
          state.app.files = list
        },
        pages: function (state, list) {
          state.app.pages = list
        },
        'main-view': function (state, view) {
          state.mainView = view
        },
        'sidebar-view': function (state, view) {
          state.sidebarView = view
        },
        'editor-view': function (state, view) {
          state.editorView = view
        },
        'current-file': function (state, file) {
          state.current = file
        },
        'preview-url': function (state, url) {
          console.log('setting preview url: ' + url)
          state.previewUrl = url
        },
        'reset-current-file': function (state, url) {
          state.current = state.defaultFile
          state.previewUrl = false
        },
        'maximize-column': function (state, info) {
          state.maximizedColumn = info
        }
      },
      actions: {
        'log': function (context, message) {
          context.commit('log', message)
        },
        'navigate': function (context, request) {
          return r.navigate(request.href, request.state)
        },
        'containers': function (context) {
          return data.getContainers()
            .then((list) => {
              context.commit('containers', list)
            })
        },
        'app': function (context) {
          return data.getApplication()
            .then((doc) => {
              context.commit('app', doc)
            })
        },
        'list-files': function (context) {
          return data.getFiles()
            .then((list) => {
              context.commit('files', list)
            })
        },
        'list-pages': function (context) {
          return data.getPages()
            .then((list) => {
              context.commit('pages', list)
            })
        },
        'reload': function (context) {
          return context.dispatch('list-pages')
            .then(() => context.dispatch('list-files'))
        },
        'get-file-contents': function (context, item) {
          return data.getFileContents(item.url)
            .then((res) => {
              // TODO: get blob for binary types
              return res.text()
            })
        },
        'open-file': function (context, file) {
          return context.dispatch('get-file-contents', file)
            .then((content) => {
              file.content = content
              context.commit('current-file', file)
              context.commit('preview-url', file.uri)
              if (context.state.editorView === 'welcome') {
                context.commit('editor-view', context.state.defaultEditorView)
              }
            })
        },
        'go-page': function (context, file) {
          let href = context.state.getAppHref('pages', file.url)
          return context.dispatch('navigate', {href: href, state: file})
        },
        'go-file': function (context, file) {
          let href = context.state.getAppHref('files', file.url)
          return context.dispatch('navigate', {href: href, state: file})
        },
        'preview-refresh': function (context) {
          context.commit('preview-url', context.state.previewUrl)
        },
        'new-file': function (context, {name, template, action}) {
          if (!/^\//.test(name)) {
            name = '/' + name
          }
          return context.state.createNewFile(name, template)
            .then((res) => {
              // Show error response
              if (res.response.status !== 201) {
                let doc = res.document
                let msg = doc.error || doc.message
                msg = `[${res.response.status}] ${msg}`
                return context.dispatch('log', new Error(msg))
              }

              context.dispatch('log', `Created ${name}`)

              context.dispatch('reload')
                .then(() => {
                  // Open the newly created file
                  let files = context.state.app.files
                  for (let i = 0; i < files.length; i++) {
                    if (files[i].url === name) {
                      context.dispatch(action, files[i])
                      break
                    }
                  }
                })
            })
        },
        'delete-file': function (context, file) {
          let list = context.state.app.files
          let ctx = context.state.sidebarView
          if (ctx === 'pages') {
            list = context.state.app.pages
          }
          let index = list.indexOf(file)
          let len = list.length
          return context.state.deleteFile(file)
            .then((res) => {
              let doc = res.document
              if (res.response.status !== 200) {
                let msg = doc.error || doc.message
                msg = `[${res.response.status}] ${msg}`
                return context.dispatch('log', new Error(msg))
              }

              return context.dispatch('reload')
                .then(() => {
                  if (len <= 1) {
                    // TODO: select next nearest file
                    context.commit('reset-current-file')
                    store.commit('editor-view', 'welcome')
                  } else if (index > -1) {
                    let neighbour = list[index - 1] || list[index + 1]
                    let href = context.state.getAppHref(ctx, neighbour.url)
                    return context.dispatch('navigate', {href: href})
                  }
                })
            })
        }
      }
    })

    let r = this.router = new Router('apps', true)
    r.add(/^apps\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+\/(pages|files)\/(.*)$/,
      ['section', 'container', 'application', 'action'],
      (match) => {
        let href = '/' + match.parts.slice(4).join('/')
        let container = match.map.container
        let application = match.map.application
        let action = match.map.action
        let file

        function findAndOpen (href) {
          let arr = data.app.files
          if (action === 'pages') {
            arr = data.app.pages
          }
          for (let i = 0; i < arr.length; i++) {
            if (arr[i].url === href) {
              store.dispatch('open-file', arr[i])
              return arr[i]
            }
          }
        }

        function trigger () {
          file = findAndOpen(href)
          if (!file) {
            // Continue route processing to trigger a 404
            return true
          }
          store.commit('main-view', 'edit')
          store.commit('sidebar-view', action)
          if (!store.state.editorView || store.state.editorView === 'welcome') {
            store.commit('editor-view', store.state.defaultEditorView)
          }
        }

        // Need to load application data
        if (container !== data.container || (container === data.container && application !== data.application)) {
          this.load(match.map.container, match.map.application)
            .then(() => {
              return trigger()
            })
        } else {
          return trigger()
        }
      })
    r.add(/^apps\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+\/(files|pages|components|new|del)$/,
      ['section', 'container', 'application', 'action'],
      (match) => {
        let container = match.map.container
        let application = match.map.application
        let action = match.map.action

        // Need to load application data
        if (container !== data.container || (container === data.container && application !== data.application)) {
          this.load(match.map.container, match.map.application)
            .then(() => {
              store.commit('reset-current-file')
              store.commit('main-view', 'edit')
              store.commit('sidebar-view', action)
              store.commit('editor-view', 'welcome')
            })
        } else {
          store.commit('reset-current-file')
          store.commit('sidebar-view', action)
          store.commit('editor-view', 'welcome')
        }
      })
    r.add(/^apps\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+$/,
      ['section', 'container', 'application'],
      (match) => {
        this.load(match.map.container, match.map.application)
          .then(() => {
            let index = store.state.getIndexFile()
            if (index) {
              let href = match.href + '/files' + index.url
              // Redirect to index page if there is one
              return r.replace(href, true)
            }
            store.commit('reset-current-file')
            store.commit('main-view', 'edit')
            store.commit('editor-view', 'welcome')
          })
      })
    r.add(/^(|apps|docs|edit|settings)$/, ['section'], (match) => {
      let section = match.map.section

      // Request with just the #
      if (section === '') {
        return r.replace('apps', true)
      } else if (section === 'apps') {
        return this.store.dispatch('containers')
          .then(() => {
            store.commit('main-view', section)
          })
      } else if (section === 'edit') {
        if (data.hasApplication()) {
          return r.replace('apps/' + data.container + '/' + data.application, true)
        } else {
          // no app being edited redirect to apps list
          return r.replace('apps', true)
        }
      }
      store.commit('main-view', section)
    })
    r.add(/^404$/, (match) => {
      store.commit('main-view', 'not-found')
    })
    r.add(/.*/, (match) => {
      store.commit('flash', r.hash)
      r.replace('404', true)
    })
  }

  get (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  ui () {
    let data = this.data

    let sidebar = {
      template: `
        <div class="sidebar"
          :class="{maximized: maximized === 'sidebar', minimized: maximized != '' && maximized !== 'sidebar'}">
          <div class="column-header">
            <nav class="tabs">
              <a v-bind:class="{selected: currentView === 'pages'}"
                @click="currentView = 'pages'"
                title="Show pages">Pages</a>
              <a v-bind:class="{selected: currentView === 'files'}"
                @click="currentView = 'files'"
                title="Show files">Files</a>
              <a v-bind:class="{selected: currentView === 'components'}"
                @click="currentView = 'components'"
                title="Show components">Components</a>
            </nav>
          </div>
          <nav class="toolbar">
            <a
              @click="quickDeleteFile"
              v-bind:class="{hidden: !shouldDelete}"
              class="danger"
              title="Confirm file deletion">Delete {{currentFile.name}}</a>
            <a
              @click="shouldDelete = true"
              v-bind:class="{hidden: !canDelete}"
              title="Delete File">➖</a>
            <a
              @click="showNewFileView"
              v-bind:class="{disabled: currentView === 'new'}"
              title="New File">➕</a>
            <a
              @click="maximized = 'sidebar'"
              :class="{hidden: maximized === 'sidebar'}"
              title="Maximize">◩</a>
            <a
              @click="maximized = ''"
              :class="{hidden: maximized !== 'sidebar'}"
              title="Minimize">▣</a>
          </nav>
          <div class="scroll">
            <component v-bind:is="currentView"></component>
          </div>
        </div>
      `,
      data: function () {
        return {
          shouldDelete: false
        }
      },
      computed: {
        maximized: {
          get: function () {
            return this.$store.state.maximizedColumn
          },
          set: function (val) {
            this.$store.commit('maximize-column', val)
          }
        },
        canDelete: function () {
          return this.$store.state.hasFile()
        },
        currentFile: function () {
          return this.$store.state.current
        },
        currentView: {
          get: function () {
            return this.$store.state.sidebarView
          },
          set: function (val) {
            var values = [val]
            var file = this.$store.state.getFile()
            if (file !== null) {
              if (val === 'files') {
                values.push(file.url)
              } else if (val === 'pages' && this.$store.state.isPage(file)) {
                values.push(file.url)
              }
            }
            let href = this.$store.state.getAppHref(...values)
            this.$store.dispatch('navigate', {href: href})
          }
        }
      },
      methods: {
        quickDeleteFile: function (e) {
          e.preventDefault()
          return this.$store.dispatch('delete-file', this.currentFile)
            .then(() => { this.shouldDelete = false })
        },
        showNewFileView: function () {
          this.previousView = this.currentView
          this.currentView = 'new'
        },
        closeNewFileView: function () {
          this.currentView = this.previousView || 'pages'
        }
      },
      components: {
        'new': {
          template: `
            <div class="new-page">
              <section>
                <h3>File Name</h3>
                <form @submit="createNewFile">
                  <input v-model="fileName" type="text" name="name" :value="fileName" />
                  <p class="small">Tip: Use <code>/path/to/file/document.md</code> to create directories when adding new files.</p>

                  <div class="template-select">
                    <h3>Template</h3>
                    <p class="small">Select an optional file template:</p>
                    <ul class="small compact-list">
                      <li>
                        <input type="radio" v-model="template"
                          id="empty-file" name="template" value="" checked />
                        <label for="empty-file">Empty File</label>
                      </li>
                      <li>
                        <input type="radio" @change="extension = '.md'" v-model="template"
                          id="markdown-partial" name="template" value="template/markdown+partial" />
                        <label for="markdown-partial">Markdown Partial</label>
                      </li>
                      <li>
                        <input type="radio" @change="extension = '.md'" v-model="template"
                          id="markdown-standalone" name="template" value="template/markdown+standalone" />
                        <label for="markdown-standalone">Markdown Standalone</label>
                      </li>
                      <li>
                        <input type="radio" v-model="template"
                          id="html-layout" @change="extension = '.html'" name="template" value="template/html+layout" />
                        <label for="html-layout">HTML Layout</label>
                      </li>
                      <li>
                        <input type="radio" v-model="template"
                          id="html-partial" @change ="extension = '.html'" name="template" value="template/html+partial" />
                        <label for="html-partial">HTML Partial</label>
                      </li>
                      <li>
                        <input type="radio" v-model="template"
                          id="html-standalone" @change="extension = '.html'" name="template" value="template/html+standalone" />
                        <label for="html-standalone">HTML Standalone</label>
                      </li>
                    </ul>
                  </div>
                  <nav class="form-actions">
                    <input @click="cancel" type="reset" name="Reset" value="Cancel" />
                    <input type="submit" name="Create" value="Create" />
                  </nav>
                </form>
              </section>
            </div>`,
          data: function () {
            return {
              fileName: '/untitled.md',
              template: '',
              extension: ''
            }
          },
          watch: {
            extension: function (val) {
              this.displayExtension = val
            }
          },
          computed: {
            displayExtension: {
              get: function () {
                return this.extension
              },
              set: function (val) {
                if (val) {
                  let current = this.fileName
                  if (/[^.]+\.[^.]*$/.test(current)) {
                    current = current.replace(/\.[^.]*$/, val)
                    this.fileName = current
                  }
                }
              }
            }
          },
          methods: {
            cancel: function (e) {
              e.preventDefault()
              this.$parent.closeNewFileView()
            },
            createNewFile: function (e) {
              e.preventDefault()
              let action = 'go-file'
              if (this.$parent.previousView === 'pages') {
                action = 'go-page'
              }
              return this.$store.dispatch('new-file', {name: this.fileName, template: this.template, action: action})
                .then(() => {
                  this.$parent.closeNewFileView()
                })
            }
          }
        },
        pages: {
          template: `
            <div class="pages-list">
              <a @click="click(item)" class="page" :class="{selected: currentFile.url === item.url}" v-for="item in list">
                <span class="name">{{item.url}}</span>
              </a>
            </div>`,
          computed: {
            currentFile: function () {
              return this.$store.state.app.current
            },
            list: function () {
              return this.$store.state.app.pages
            }
          },
          methods: {
            click: function (item) {
              return this.$store.dispatch('go-page', item)
            }
          }
        },
        files: {
          template: `
            <div class="files-list">
              <a @click="click(item)" class="file" :class="{selected: currentFile.url === item.url}" v-for="item in list">
                <span class="name">{{item.url}}</span>
              </a>
            </div>`,
          computed: {
            currentFile: function () {
              return this.$store.state.app.current
            },
            list: function () {
              return this.$store.state.app.files
            }
          },
          methods: {
            click: function (item) {
              return this.$store.dispatch('go-file', item)
            }
          }
        },
        components: {
          template: `
            <div class="components-list">
            </div>`,
          computed: {
            list: function () {
              return this.$store.state.app.components
            }
          }
        }
      }
    }

    let preview = {
      template: `
        <div
          :class="{maximized: maximized === 'preview', minimized: maximized != '' && maximized !== 'preview'}"
          class="preview">
          <div class="column-header">
            <h2>Preview</h2>
            <div class="column-options">
              <nav class="tabs">
                <!-- <a href="#preview" title="Publish preview">Preview</a> -->
                <!-- <a href="#docs" title="Browse the help & documentation">Docs</a> -->
              </nav>
            </div>
          </div>
          <nav class="toolbar clearfix">
            <h2>{{path}}</h2>
            <a @click="refresh()"
               :class="{hidden: path == ''}">Reload</a>
            <a
              @click="maximized = 'preview'"
              :class="{hidden: maximized === 'preview'}"
              title="Maximize">◩</a>
            <a
              @click="maximized = ''"
              :class="{hidden: maximized !== 'preview'}"
              title="Minimize">▣</a>
          </nav>
          <iframe :src="src" class="publish-preview"></iframe>
        </div>
      `,
      data: function () {
        return {
          path: '',
          src: ''
        }
      },
      computed: {
        maximized: {
          get: function () {
            return this.$store.state.maximizedColumn
          },
          set: function (val) {
            this.$store.commit('maximize-column', val)
          }
        },
        url: {
          get: function () {
            return this.$store.state.previewUrl
          },
          set: function (val) {
            return this.$store.commit('preview-url', val)
          }
        }
      },
      watch: {
        url: function (url) {
          console.log('url watcher called: ' + url)
          this.refresh(url)
        }
      },
      methods: {
        refresh (url) {
          if (url === '') {
            this.path = ''
            this.src = ''
            return
          }
          // If the src attribute will not change the page
          // won't be refreshed so we need to call reload()
          if (url === this.path) {
            let frame = document.querySelector('.publish-preview')
            return frame.contentDocument.location.reload()
          }
          this.path = url
          this.src = this.getPreviewUrl(url)
        },
        getPreviewUrl (url) {
          if (url) {
            url = url.replace(/^\//, '')
          }
          return document.location.origin + data.app.url + (url || '')
        }
      }
    }

    let editor = {
      template: `
        <div
          :class="{maximized: maximized === 'editor', minimized: maximized != '' && maximized !== 'editor'}"
          class="editor">
          <div class="column-header">
            <h2>Editor</h2>
            <div class="column-options">
              <nav class="tabs">
                <a v-bind:class="{selected: currentView === 'file-editor', hidden: hidden}"
                  @click="currentView = 'file-editor'"
                  title="Show file editor">File</a>
                <a v-bind:class="{selected: currentView === 'source-editor', hidden: hidden}"
                  @click="currentView = 'source-editor'"
                  title="Show source editor">Source</a>
                <a v-bind:class="{selected: currentView === 'visual-editor', hidden: hidden}"
                  @click="currentView = 'visual-editor'"
                  title="Show visual editor">Visual</a>
              </nav>
            </div>
          </div>
          <nav class="toolbar clearfix">
            <h2>{{title}}</h2>
            <a @click="saveAndRun"
              v-bind:class="{hidden: currentView != 'source-editor'}" href="#" title="Save & Run">Save & Run</a>
            <a
              @click="maximized = 'editor'"
              :class="{hidden: maximized === 'editor'}"
              title="Maximize">◩</a>
            <a
              @click="maximized = ''"
              :class="{hidden: maximized !== 'editor'}"
              title="Minimize">▣</a>
          </nav>
          <component v-bind:is="currentView"></component>
        </div>
      `,
      computed: {
        maximized: {
          get: function () {
            return this.$store.state.maximizedColumn
          },
          set: function (val) {
            this.$store.commit('maximize-column', val)
          }
        },
        hidden: function () {
          return !this.$store.state.hasFile()
        },
        currentFile: function () {
          return this.$store.state.app.current
        },
        currentView: {
          get: function () {
            return this.$store.state.editorView
          },
          set: function (view) {
            this.$store.commit('editor-view', view)
          }
        }
      },
      watch: {
        currentFile: function (file) {
          this.title = file.url
        }
      },
      data: function () {
        return {
          title: ''
        }
      },
      methods: {
        saveAndRun: function (e) {
          e.preventDefault()
          let file = this.currentFile
          let value = file.content
          data.saveFile(file, value)
            .then((res) => {
              let doc = res.document
              if (doc.ok) {
                this.$store.dispatch('preview-refresh')
              }
            })
        }
      },
      components: {
        welcome: {
          template: `
            <div class="welcome scroll">
              <p>Select a page or file to start editing.</p>
            </div>
          `
        },
        'file-editor': {
          template: `<div class="file-editor">
              <div class="scroll panel">
                <h2 class="file-info"><span v-bind:class="{hidden: !file.dir}">🗀</span><span v-bind:class="{hidden: file.dir}">🗎</span>&nbsp;{{file.name}}</h2>
                <section>
                  <h3>Rename File</h3>
                  <p>Choose a new name for your file.</p>
                  <form class="rename">
                    <input type="text" name="fileName" :value="file.name" />
                    <input type="submit" name="Rename" value="Rename" />
                  </form>
                </section>
                <section>
                  <h3>Delete File</h3>
                  <p v-bind:class="{hidden: confirmDelete}">Danger zone: be careful!</p>
                  <div>
                    <button @click="confirmDelete = true"
                      v-bind:class="{hidden: confirmDelete}"
                      class="danger">Delete {{file.url}}</button>
                  </div>
                  <div v-bind:class="{hidden: !confirmDelete}">
                    <p>Are you sure you want to delete {{file.url}}?<br />
                    <small>
                      Deleting a file is irreversible, it cannot be undone.
                    </small>
                    </p>
                    <nav class="form-actions">
                      <button @click="confirmDelete = false">Cancel</button>
                      <button @click="doDelete" class="danger">Delete</button>
                    </nav>
                  </div>
                </section>
                <section>
                  <h3>File Info</h3>
                  <ul class="small compact-list">
                    <li>Name: {{file.name}}</li>
                    <li>URL : {{file.uri}}</li>
                    <li v-bind:class="{hidden: !file.dir}">Directory: yes</li>
                    <li v-bind:class="{hidden: file.dir}">Size: {{file.size}} bytes</li>
                    <li v-bind:class="{hidden: file.dir}">Mime: {{file.mime}}</li>
                  </ul>
                </section>
              </div>
            </div>`,
          data: function () {
            return {
              confirmDelete: false
            }
          },
          computed: {
            file: function () {
              return this.$store.state.current
            }
          },
          methods: {
            doDelete: function () {
              return this.$store.dispatch('delete-file', this.file)
            }
          }
        },
        'source-editor': {
          template: `<div class="source-editor">
              <div class="text-editor"></div>
            </div>`,
          data: function () {
            return {
              mirror: null
            }
          },
          computed: {
            currentFile: function () {
              return this.$store.state.app.current
            },
            value: {
              get: function () {
                return this.$store.state.current.content
              },
              set: function (val) {
                this.$store.state.current.content = val
              }
            }
          },
          watch: {
            currentFile: function (file) {
              if (file && file.mime) {
                this.setCodeMirror({value: file.content, mode: this.getModeForMime(file.mime)})
              }
            }
          },
          methods: {
            getModeForMime (mime) {
              // remove charset info
              mime = mime.replace(/;.*$/, '')
              switch (mime) {
                case 'text/html':
                  return 'htmlmixed'
                case 'text/x-markdown':
                  return 'yaml-frontmatter'
              }
              return mime
            },
            changes: function (cm, changes) {
              this.value = this.mirror.getValue()
            },
            setCodeMirror: function (options) {
              options = options || {}
              let p = document.querySelector('.text-editor')

              if (this.mirror) {
                // TODO: verify listener is removed
                this.mirror.off('changes', this.changes)
                if (p.firstChild) {
                  p.removeChild(p.firstChild)
                }
              }
              this.mirror = CodeMirror(p, {
                value: options.value || '',
                mode: options.mode || 'htmlmixed',
                theme: options.theme || 'solarized dark',
                lineNumbers: true,
                keyMap: 'vim'
              })
              this.mirror.on('changes', this.changes)
            }
          },
          mounted: function () {
            let item = this.currentFile
            // Handles setting file content when switching tabs
            if (item && item.content) {
              this.setCodeMirror({value: item.content, mode: this.getModeForMime(item.mime)})
            }
          }
        },
        'visual-editor': {
          template: `<div class="visual-editor"></div>`
        }
      }
    }

    let app = new Vue({
      template: `
        <main>
          <app-header></app-header>
          <app-main></app-main>
          <app-footer></app-footer>
        </main>
      `,
      store: this.store,
      components: {
        'app-header': {
          template: `
              <header class="clearfix">
                <nav class="main">
                  <a
                    @click="$store.dispatch('navigate', {href: 'apps'})"
                    :class="{selected: selectedView === 'apps'}"
                    title="View and edit applications">Apps</a>
                  <a
                    @click="$store.dispatch('navigate', {href: 'docs'})"
                    :class="{selected: selectedView === 'docs'}"
                    title="Documentation">Docs</a>
                  <a
                    @click="$store.dispatch('navigate', {href: 'edit'})"
                    :class="{selected: selectedView === 'edit', hidden: !this.$store.state.hasApplication()}"
                    title="Edit Current Application">{{name}}</a>
                  <a
                    @click="$store.dispatch('navigate', {href: 'settings'})"
                    :class="{selected: selectedView === 'settings'}"
                    title="Settings">Settings</a>
                </nav>
                <nav class="home">
                  <a
                    @click="$store.dispatch('navigate', {href: '/'})"
                    class="home"
                    title="Home page">Ꝏ </a>
                </nav>
              </header>
            `,
          computed: {
            name: function () {
              return this.$store.state.application
            },
            selectedView: function () {
              return this.$store.state.mainView
            }
          }
        },
        'app-main': {
          template: `
            <component v-bind:is="currentView"></component>
          `,
          computed: {
            currentView: function () {
              return this.$store.state.mainView
            }
          },
          components: {
            'not-found': {
              template: `
                <div class="content-main not-found">
                  <h2>Not Found</h2>
                  <p>Oops, we could not find a page that matched your request <code>{{flash}}</code></p>
                </div>
              `,
              computed: {
                flash: function () {
                  let flash = this.$store.state.flash
                  if (flash) {
                    flash = `#${flash}`
                  }
                  return flash
                }
              }
            },
            'apps': {
              template: `
                <div class="content-main">
                  <div class="content">
                    <div class="containers" v-for="container in list">
                      <span :class="{hidden: !container.protected}">🔒&nbsp;</span>
                      <span class="name container">{{container.name}}</span>
                      <p class="small">{{container.description}}</p>
                      <ul class="compact-list">
                        <div class="app" v-for="app in container.apps">
                            <span :class="{hidden: !app.protected}">🔒&nbsp;</span>
                            <span class="name">{{app.name}}</span>
                            <p class="small">URL: {{app.url}}<br />{{app.description}}
                              <p class="app-actions">
                                <a class="name" @click="$store.dispatch('navigate', {href: linkify(container, app)})" :title="title(app, 'Edit')">Edit</a>
                                <a class="name" :href="linkify(container, app, true)" :title="title(app, 'Open')">Open</a>
                              </p>
                            </p>
                        </div>
                      </ul>
                    </div>
                  </div>
                </div>
              `,
              computed: {
                list: function () {
                  return this.$store.state.containers
                }
              },
              methods: {
                linkify: function (c, a, open) {
                  if (open) {
                    return a.url
                  }
                  return `apps/${c.name}/${a.name}`
                },
                title: function (a, prefix) {
                  return `${prefix} ${a.name}`
                }
              }
            },
            'docs': {
              template: `
                <div class="content-main">
                  <iframe class="docs" src="/docs/"></iframe>
                </div>
              `
            },
            'edit': {
              template: `
                <div class="content-main">
                  <div class="content">
                    <sidebar></sidebar>
                    <editor></editor>
                    <preview></preview>
                  </div>
                </div>
              `,
              components: {
                sidebar,
                editor,
                preview
              }
            },
            'settings': {
              template: `
                <div class="content-main">
                  <h3>Settings</h3>
                </div>
              `
            }
          }
        },
        'app-footer': {
          template: `
            <footer>
              <p class="log" v-bind:class="{error: error}">{{prefix}}{{message}}</p>
            </footer>
          `,
          computed: {
            message: function () {
              return this.$store.state.log.toString()
            },
            error: function () {
              return (this.$store.state.log.last instanceof Error)
            },
            prefix: function () {
              if (this.message && this.error) {
                return '! '
              } else if (this.message && !this.error) {
                return '# '
              }
              return ''
            }
          }
        }
      }
    })

    app.$mount('main')
    return app
  }

  load (container, application) {
    let data = this.data

    this.data.setApplication(container, application)

    this.store.dispatch('log', `Loading app from ${data.url}`)
    return this.store.dispatch('app')
      .then(() => this.store.dispatch('list-files'))
      .then(() => this.store.dispatch('list-pages'))
      .then(() => {
        this.store.commit('sidebar-view', 'pages')
        this.store.dispatch('log', 'Done')
      })
      .catch((err) => this.store.dispatch('log', err))
  }

  init () {
    this.ui()
    this.router.start()
  }
}

let app = new EditorApplication()
app.init()
