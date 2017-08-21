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
    this.loc = loc
    this.url = `/api/${loc.container}/${loc.application}/`
    this.raw = `/apps/raw/${loc.container}/${loc.application}`

    // current application
    this.app = {
      url: '',
      identifier: '',
      owner: loc.container,
      pages: [],
      files: []
    }
  }

  json (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  getFileContents (pathname) {
    let url = this.raw
    return fetch(url + pathname)
      .catch((err) => err)
  }

  getApplication () {
    return this.json(this.url)
      .then((app) => {
        // merge properties
        for (let k in app) {
          this.app[k] = app[k]
        }
        this.app.identifier = this.app.owner + ' / ' + this.app.name
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
  constructor (loc) {
    this.loc = loc
    this.bus = new Vue()
    this.data = new AppDataSource(this.loc)
  }

  get (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  ui () {
    let data = this.data
    let bus = this.bus

    this.header = new Vue({
      template: `
          <header>
            <nav>
              <a @click="navigate" href="#apps" title="All applications">Apps</a>
              <a @click="navigate" href="#docs" title="Documentation">Docs</a>
              <a @click="navigate" href="#editor" title="Editor">Editor</a>
              <a @click="navigate" href="#settings" title="Settings">Settings</a>
              <a @click="navigate" href="/" title="Home page">Ꝏ</a>
            </nav>
            <div class="app-id"></div>
          </header>
        `,
      data: {
        currentView: 'editor'
      },
      methods: {
        navigate: function (e) {
          e.preventDefault()
          console.log(this)
        }
      }
    })

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
          <span class="name">{{name}}</span>
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
          pages: data.app.pages,
          files: data.app.files,
          components: [],
          currentView: ''
        }
      },
      created: function () {
        bus.$on('sidebar:reload', this.reload)
      },
      methods: {
        showNewFileView: function () {
          this.previousView = this.currentView
          this.currentView = 'new-file'
        },
        closeNewFileView: function () {
          this.currentView = this.previousView
        },
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
              bus.$emit('files:load', list)
            })
        },
        reload: function (fn) {
          return this.loadPages()
            .then(() => {
              return this.loadFiles()
                .then(() => {
                  if (typeof fn === 'function') {
                    return fn()
                  }
                })
            })
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
              console.log('extension changed to: ' + val)
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
          <div class="column-header">
            <h2>{{path}}</h2>
            <div class="column-options">
              <nav class="tabs">
                <a href="#preview" title="Publish preview">Preview</a>
                <a href="#docs" title="Browse the help & documentation">Docs</a>
              </nav>
            </div>
          </div>
          <nav class="toolbar clearfix">
            <a href="#reload">Reload</a>
          </nav>
          <iframe :src="url" class="live"></iframe>
        </div>
      `,
      data: function () {
        return {
          path: '/',
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
          if (url === this.path) {
            let frame = document.querySelector('.live')
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
    })

    this.editor = new Vue({
      template: `
        <div class="editor">
          <div class="column-header">
            <h2>{{title}}</h2>
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
          <nav class="toolbar">
            <a @click="saveAndRun"
              v-bind:class="{hidden: currentView != 'source-editor'}" href="#" title="Save & Run">Save & Run</a>
          </nav>
          <component v-bind:is="currentView"></component>
        </div>
      `,
      data: function () {
        let defaultTitle = 'Editor'
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
              if (this.currentView === 'source-editor') {
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
            let item = data.currentFile
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
    })

    this.footer = new Vue({
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

    let main = new Vue({el: 'main'})

    // mount views
    this.header.$mount('header')
    this.identifier.$mount('.app-id')
    this.switcher.$mount('.switcher')
    this.sidebar.$mount('.sidebar')
    this.editor.$mount('.editor')
    this.preview.$mount('.preview')
    this.footer.$mount('footer')

    return main
  }

  log (msg) {
    this.bus.$emit('log', msg)
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
        this.preview.refresh()
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
        // console.log(this.data)
      })
  }
}

let app = new EditorApplication(new LocationParser().parse())
app.init()
