/* globals Vue Vuex CodeMirror document fetch document history window */

class Router {
  constructor (href) {
    this.defaultHref = href
    this.routes = []
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
    return document.location.hash.replace(/^#/, '')
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

    let r, ptn, fn
    for (let i = 0; i < this.routes.length; i++) {
      r = this.routes[i]
      ptn = r.ptn
      fn = r.fn
      if (typeof ptn === 'string' && href === ptn) {
        fn(result(href, r))
        break
      } else if (ptn instanceof RegExp && ptn.test(href)) {
        fn(result(href, r))
        break
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

class AppDataSource {
  constructor () {
    this.api = '/api/'
    this.containers = []
    this.setApplication('', '')
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
      current: {content: ''}
    }
  }

  hasApplication () {
    return this.container && this.application
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
    let bus = this.bus = new Vue()
    let data = this.data = new AppDataSource()

    let r = this.router = new Router('apps')
    r.add(/^apps\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+$/, ['section', 'container', 'application'], (match) => {
      this.load(match.map.container, match.map.application)
        .then(() => {
          bus.$emit('view:select', 'edit')
        })
    })
    r.add(/^(apps|docs|edit|settings)$/, ['section'], (match) => {
      let section = match.map.section
      if (section === 'apps') {
        return this.store.dispatch('containers')
          .then(() => {
            bus.$emit('view:select', section)
          })
      } else if (section === 'edit') {
        if (data.hasApplication()) {
          return r.replace('apps/' + data.container + '/' + data.application, true)
        } else {
          // no app being edited redirect to apps list
          return r.replace('apps', true)
        }
      }
      bus.$emit('view:select', section)
    })

    this.store = new Vuex.Store({
      state: this.data,
      mutations: {
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
        'current-file': function (state, file) {
          state.app.current = file
        }
      },
      actions: {
        'navigate': function (context, href) {
          r.navigate(href)
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
            .then(context.dispatch('list-files'))
        },
        'index-page-select': function (context) {
          let files = context.state.app.files
          for (let i = 0; i < files.length; i++) {
            // got a published index page whether the source is
            // HTML or markdown
            if (files[i].uri === '/index.html') {
              return context.dispatch('open-file', files[i])
            }
          }
        },
        'get-file-contents': function (context, item) {
          return data.getFileContents(item.url)
            .then((res) => {
              // TODO: get blob for binary types
              return res.text()
            })
        },
        'open-file': function (context, file) {
          // NOTE: we need to store in data source for
          // NOTE: switching editor tabs
          // context.state.app.currentFile = item

          /*
          if (this.currentView === 'welcome') {
            this.currentView = this.defaultOpenView
          }
          */

          console.log('opening file')
          console.log(file)

          return context.dispatch('get-file-contents', file)
            .then((content) => {
              file.content = content
              context.commit('current-file', file)
            })

          /*
          data.getFileContents(item.url)
            .then((res) => {
              // TODO: get blob for binary types
              return res.text()
            }).then((content) => {
              item.content = content
              this.title = item.url
              if (this.$children[0] && this.$children[0].showSourceText) {
                this.$children[0].showSourceText(item)
              }
              bus.$emit('open:complete', item)
            })
          */
        }
      }
    })
  }

  get (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  ui () {
    let data = this.data
    let bus = this.bus

    let sidebar = {
      template: `
        <div class="sidebar">
          <div class="column-header">
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
          </div>
          <nav class="toolbar">
            <a
              v-bind:class="{disabled: currentView === 'new-file'}"
              href="#" title="Delete File">➖</a>
            <a
              @click="showNewFileView"
              v-bind:class="{disabled: currentView === 'new-file'}"
              href="#" title="New File">➕</a>
          </nav>
          <div class="scroll">
            <component v-bind:is="currentView"></component>
          </div>
        </div>
      `,
      data: function () {
        return {
          currentView: ''
        }
      },
      created: function () {
        bus.$on('sidebar:reload', this.reload)
        bus.$on('sidebar:select', (view) => {
          this.currentView = view
        })
      },
      methods: {
        showNewFileView: function () {
          this.previousView = this.currentView
          this.currentView = 'new-file'
        },
        closeNewFileView: function () {
          this.currentView = this.previousView
        },
        reload: function (next) {
          this.$store.dispatch('reload')
            .then(next)
        }
      },
      components: {
        'new-file': {
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

              let name = this.fileName
              if (!/^\//.test(name)) {
                name = '/' + name
              }

              return data.createNewFile(name, this.template)
                .then((res) => {
                  // Show error response
                  if (res.response.status !== 201) {
                    let doc = res.document
                    let msg = doc.error || doc.message
                    msg = `[${res.response.status}] ${msg}`
                    return bus.$emit('log', new Error(msg))
                  }

                  bus.$emit('log', `Created ${this.fileName}`)
                  bus.$emit('sidebar:reload', () => {
                    // Open the newly created file
                    for (let i = 0; i < data.app.files.length; i++) {
                      if (data.app.files[i].url === this.fileName) {
                        bus.$emit('open:file', data.app.files[i])
                        break
                      }
                    }
                  })

                  this.$parent.closeNewFileView()
                })
            }
          }
        },
        pages: {
          template: `
            <div class="pages-list">
              <a @click="click(item)" class="page" v-for="item in list">
                <span class="name">{{item.url}}</span>
              </a>
            </div>`,
          computed: {
            list: function () {
              return this.$store.state.app.pages
            }
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
          computed: {
            list: function () {
              return this.$store.state.app.files
            }
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
        <div class="preview">
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
            <a @click="refresh(true)"
               :class="{hidden: path == ''}">Reload</a>
          </nav>
          <iframe :src="url" class="publish-preview"></iframe>
        </div>
      `,
      data: function () {
        return {
          path: '',
          url: ''
        }
      },
      created: function () {
        bus.$on('open:complete', (item) => {
          let url = item.uri
          let all = /\.html?$/
          // Refresh preview when switching on page types
          if (all.test(url)) {
            this.refresh(url)
          }
        })

        bus.$on('preview:refresh', (url) => {
          this.refresh(url)
        })
      },
      methods: {
        refresh (url) {
          // If the src attribute will not change the page
          // won't be refreshed so we need to call reload()
          if (url === true || url === this.path || url === this.url) {
            let frame = document.querySelector('.publish-preview')
            return frame.contentDocument.location.reload()
          }
          this.path = url || '/'
          this.url = this.getPreviewUrl(url)
        },
        getPreviewUrl: function (url) {
          if (url) {
            url = url.replace(/^\//, '')
          }
          return document.location.origin + data.app.url + (url || '')
        }
      }
    }

    let editor = {
      template: `
        <div class="editor">
          <div class="column-header">
            <h2>Editor</h2>
            <div class="column-options">
              <nav class="tabs">
                <a v-bind:class="{selected: currentView === 'file-editor', hidden: currentFile === null}"
                  @click="currentView = 'file-editor'"
                  href="#" title="Show file editor">File</a>
                <a v-bind:class="{selected: currentView === 'source-editor', hidden: currentFile === null}"
                  @click="currentView = 'source-editor'"
                  href="#" title="Show source editor">Source</a>
                <a v-bind:class="{selected: currentView === 'visual-editor', hidden: currentFile === null}"
                  @click="currentView = 'visual-editor'"
                  href="#"  title="Show visual editor">Visual</a>
              </nav>
            </div>
          </div>
          <nav class="toolbar clearfix">
            <h2>{{title}}</h2>
            <a @click="saveAndRun"
              v-bind:class="{hidden: currentView != 'source-editor'}" href="#" title="Save & Run">Save & Run</a>
          </nav>
          <component v-bind:is="currentView"></component>
        </div>
      `,
      data: function () {
        let defaultTitle = ''
        return {
          title: defaultTitle,
          defaultTitle: defaultTitle,
          currentView: 'welcome',
          defaultOpenView: 'source-editor',
          currentFile: null
        }
      },
      created: function () {
        bus.$on('open:file', (item) => {
          this.open(item)
        })
        bus.$on('close:file', () => {
          this.close()
        })
      },
      methods: {
        saveAndRun: function (e) {
          e.preventDefault()
        },
        close: function () {
          if (this.currentFile) {
            this.currentView = 'welcome'
            this.title = this.defaultTitle
            this.currentFile = null
            // Must clear for tabs too
            data.currentFile = null
          }
        },
        open: function (item) {
          // NOTE: we need to store in data source for
          // NOTE: switching editor tabs
          data.currentFile = this.currentFile = item

          if (this.currentView === 'welcome') {
            this.currentView = this.defaultOpenView
          }

          data.getFileContents(item.url)
            .then((res) => {
              // TODO: get blob for binary types
              return res.text()
            }).then((content) => {
              item.content = content
              this.title = item.url
              if (this.$children[0] && this.$children[0].showSourceText) {
                this.$children[0].showSourceText(item)
              }
              bus.$emit('open:complete', item)
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
              file: {name: '', url: ''},
              confirmDelete: false
            }
          },
          created: function () {
            bus.$on('open:complete', (item) => {
              this.file = item
            })
          },
          mounted: function () {
            if (data.currentFile) {
              this.file = data.currentFile
            }
          },
          methods: {
            doDelete: function () {
              return data.deleteFile(this.file)
                .then((res) => {
                  let doc = res.document
                  if (res.response.status !== 200) {
                    let msg = doc.error || doc.message
                    msg = `[${res.response.status}] ${msg}`
                    return bus.$emit('log', new Error(msg))
                  }
                  bus.$emit('log', `Deleted ${this.file.url}`)
                  bus.$emit('sidebar:reload')
                })
            }
          }
        },
        'source-editor': {
          template: `<div class="source-editor">
              <div class="text-editor"></div>
            </div>`,
          data: function () {
            return {
              mirror: null,
              canSave: false
            }
          },
          methods: {
            closeFile: function (e) {
              e.preventDefault()
              this.canSave = false
              bus.$emit('close:file')
            },
            saveAndRun: function (e) {
              e.preventDefault()
              let file = data.currentFile
              let value = this.mirror.getValue()
              data.saveFile(file, value)
                .then((res) => {
                  let doc = res.document
                  if (doc.ok) {
                    bus.$emit('preview:refresh', file.uri)
                  }
                })
            },
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
              // console.log(changes)
            },
            setCodeMirror: function (options) {
              this.canSave = true
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
            },
            showSourceText: function (item) {
              this.setCodeMirror({value: item.content, mode: this.getModeForMime(item.mime)})
            }
          },
          mounted: function () {
            let item = this.$store.state.app.current
            console.log('source mounted')
            console.log(item)
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
                <nav>
                  <a
                    @click="$store.dispatch('navigate', 'apps')"
                    :class="{selected: selectedView === 'apps'}"
                    href="#apps" title="All applications">Apps</a>
                  <a
                    @click="$store.dispatch('navigate', 'docs')"
                    :class="{selected: selectedView === 'docs'}"
                    href="#docs" title="Documentation">Docs</a>
                  <a
                    @click="$store.dispatch('navigate', 'edit')"
                    :class="{selected: selectedView === 'edit', hidden: $store.state.container === ''}"
                    href="#edit" title="Edit Application">Edit</a>
                  <a
                    @click="$store.dispatch('navigate', 'settings')"
                    :class="{selected: selectedView === 'settings'}"
                    href="#settings" title="Settings">Settings</a>
                </nav>
                <div class="app-id">
                  <a href="/" title="Home page">Ꝏ</a>
                  <span class="name">{{name}}</span>
                </div>
              </header>
            `,
          data: function () {
            return {
              selectedView: ''
            }
          },
          computed: {
            name: function () {
              return this.$store.state.app.identifier
            }
          },
          created: function () {
            bus.$on('view:select', (view) => {
              this.selectedView = view
            })
          }
        },
        'app-main': {
          template: `
            <component v-bind:is="currentView"></component>
          `,
          data: function () {
            return {
              currentView: ''
            }
          },
          created: function () {
            bus.$on('view:select', (view) => {
              this.currentView = view
            })
          },
          components: {
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
                                <a class="name" @click="$store.dispatch('navigate', linkify(container, app))" :title="title(app, 'Edit')">Edit</a>
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
              <p class="log" v-bind:class="{error: error}">{{message}}</p>
            </footer>
          `,
          data: function () {
            return {
              message: '',
              error: false
            }
          },
          created: function () {
            bus.$on('log', this.log)
          },
          methods: {
            log: function (msg) {
              let err = (msg instanceof Error)
              if (err) {
                msg = '! ' + msg
                this.error = true
              } else {
                msg = '# ' + msg
                this.error = false
              }
              this.message = msg
            }
          }
        }
      }
    })

    app.$mount('main')
    return app
  }

  load (container, application) {
    let bus = this.bus
    let data = this.data

    this.data.setApplication(container, application)

    bus.$emit('log', `Loading app from ${data.url}`)
    return this.store.dispatch('app')
    .then(() => {
      this.store.dispatch('list-files')
        .then(this.store.dispatch('list-pages'))
        .then(() => {
          bus.$emit('sidebar:select', 'pages')
          bus.$emit('log', 'Done')
          return this.store.dispatch('index-page-select')
        })
        .catch((err) => bus.$emit('log', err))
    })
  }

  init () {
    this.ui()
    this.router.start()
  }
}

let app = new EditorApplication()
app.init()
