/* globals Vue Vuex CodeMirror document fetch document history window */

class ColumnManager {
  constructor () {
    this.state = null
    this.styles = null
    this.maximized = ''

    this.doDrag = (e) => {
      e.stopImmediatePropagation()
      let target = this.state.target
      let parent = this.state.parent
      let index = this.state.index
      let maximum = this.state.maximum
      let tb = target.getBoundingClientRect()

      // Out of bounds cursor
      if (e.clientX < 0 || e.clientX > maximum || e.clientX < tb.left) {
        return
      }

      // Try to stop other events interfering
      document.querySelector('body').setAttribute('style', 'user-select: none; pointer-events: none; cursor: ew-resize;')

      // Resize target column by percentage
      let percent = Math.round(((e.clientX - tb.left) / maximum) * 100)
      target.setAttribute('style', 'max-width: none; width:' + percent + '%')

      // Columns after the target being resized, they are
      // compacted but maintain the aspect ratio
      let taken = this.state.widths + percent
      let remainder = 100 - taken
      let i, n
      for (i = index + 1; i < parent.childNodes.length; i++) {
        n = parent.childNodes[i]
        if (n.nodeType === 1) {
          n.setAttribute('style', 'max-width: none; width: ' + (remainder * this.state.ratios[i]) + '%')
        }
      }
    }

    this.stopDrag = (e) => {
      e.stopImmediatePropagation()
      document.querySelector('body').removeAttribute('style')
      window.removeEventListener('mousemove', this.doDrag)
      this.state = null
    }
  }

  startDrag (e) {
    e.stopImmediatePropagation()

    // Target to reize
    let target = e.currentTarget.parentNode

    // Parent gives overall available width
    let parent = target.parentNode

    // Width of all columns to calculate percentage
    let pb = parent.getBoundingClientRect()
    let maximum = pb.right - pb.left

    this.state = {
      target: target,
      parent: parent,
      maximum: maximum,
      index: undefined,
      widths: 0,
      ratios: []
    }

    // Used to track remaining available pixels
    let total = 0

    let i, n, b, w, ratio, percent
    for (i = 0; i < parent.childNodes.length; i++) {
      n = parent.childNodes[i]
      if (n.nodeType !== 1) {
        continue
      }

      b = n.getBoundingClientRect()
      w = b.right - b.left

      // Get ratios of subsequent columns
      if (this.state.index !== undefined) {
        // How much of the remaining space is used by this column
        ratio = w / (maximum - total)
        // Sparse array!
        this.state.ratios[i] = ratio
      }

      ratio = w / maximum
      percent = Math.round(ratio * 100)

      if (n === target) {
        this.state.index = i
        total += w
      }

      // Fix widths of previous columns
      if (this.state.index === undefined) {
        total += w
        this.state.widths += percent
        n.setAttribute('style', 'max-width: none; width:' + percent + '%')
      }
    }

    // Start the drag operation
    window.addEventListener('mousemove', this.doDrag)

    // Need to capture on the window for mouse up outside
    window.addEventListener('mouseup', this.stopDrag)
  }

  maximize (className) {
    const el = document.querySelector('.' + className)
    const parent = el.parentNode
    this.styles = {}
    parent.childNodes.forEach((n, index) => {
      if (n.nodeType !== 1) {
        return
      }
      this.styles[index] = n.getAttribute('style')
      if (n === el) {
        n.setAttribute('style', 'max-width: none; width: 100%;')
      } else {
        n.setAttribute('style', 'max-width: none; width: 0%;')
      }
    })
    this.maximized = className
  }

  minimize (className) {
    if (!className) {
      return
    }
    const el = document.querySelector('.' + className)
    const parent = el.parentNode
    parent.childNodes.forEach((n, index) => {
      if (n.nodeType !== 1) {
        return
      }
      if (this.styles[index]) {
        n.setAttribute('style', this.styles[index])
      } else {
        n.removeAttribute('style')
      }
    })
    this.styles = null
    this.maximized = ''
  }

  // Remove inline styles from columns will restore columns
  // to the defaults declared in the stylesheet.
  reset () {
    let parent = document.querySelector('.content-main > .content')
    let i, n
    this.styles = {}
    for (i = 0; i < parent.childNodes.length; i++) {
      n = parent.childNodes[i]
      if (n.nodeType !== 1) {
        continue
      }
      n.removeAttribute('style')
    }
  }
}

class Router {
  constructor (href, strip) {
    this.defaultHref = href
    this.routes = []
    this.strip = strip
  }

  navigate (href, state) {
    let url = this.url(href)
    // history.pushState({href: href, state: state}, '', url)
    history.pushState({href: href, state: null}, '', url)
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
    return m
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
    this.previewRefresh = false

    this.log = new Log()

    this._flash = undefined

    // State for edit mode columns
    this.columns = new ColumnManager()
  }

  get flash () {
    let f = this._flash
    this._flash = undefined
    return f
  }

  set flash (msg) {
    this._flash = msg
  }

  isDirty () {
    if (this.app) {
      for (let i = 0; i < this.app.files.length; i++) {
        if (this.app.files[i].dirty) {
          return true
        }
      }
    }
    return false
  }

  getAppHref (...args) {
    let p = ['apps', this.container, this.application]

    // TODO: ensure we never get passed undefined / null etc
    // TODO: and remove this
    args = args.filter((val) => {
      return val
    })

    args = args.map((val) => {
      return val.replace(/^\//, '')
    })
    p.push(...args)
    return p.join('/')
  }

  getIndexFile () {
    let files = this.app.files || []
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
      let pages = this.app.pages || []
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

  isDirectory () {
    return this.hasFile() && this.current.dir
  }

  hasFile () {
    return this.app.current.url !== undefined
  }

  isPage (file) {
    return file && file.page !== undefined
  }

  json (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  createNewApp (app) {
    let url = this.api + 'user/'
    let opts = {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(app)
    }

    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
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

  renameFile (file, newName) {
    let url = this.url + 'files' + file.url
    let opts = {
      method: 'POST',
      headers: {
        Location: newName
      }
    }
    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  deleteApp (container, application) {
    let url = `${this.api}${container}/${application}/`
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
    let data = this.state = new AppDataSource()
    let store = this.store = new Vuex.Store({
      state: this.state,
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
          if (state.hasFile()) {
            state.current.editorView = view
          }
        },
        'current-file': function (state, file) {
          if (!file.editorView) {
            file.editorView = state.defaultEditorView
          }
          state.current = file
        },
        'current-file-dirty': function (state, val) {
          state.current.dirty = val
        },
        'preview-url': function (state, url) {
          state.previewUrl = url
        },
        'preview-refresh': function (state, toggle) {
          state.previewRefresh = toggle
        },
        'reset-current-file': function (state, url) {
          state.current = state.defaultFile
          state.previewUrl = ''
        },
        'maximize-column': function (state, info) {
          // Maximizing
          if (info) {
            state.columns.maximize(info)
          // Minimizing
          } else {
            state.columns.minimize(state.columns.maximized)
          }
        }
      },
      actions: {
        'resize-column': function (context, e) {
          context.state.columns.startDrag(e)
        },
        'log': function (context, message) {
          context.commit('log', message)
        },
        'navigate': function (context, request) {
          return r.navigate(request.href, request.state)
        },
        'containers': function (context) {
          return context.state.getContainers()
            .then((list) => {
              context.commit('containers', list)
            })
        },
        'new-app': function (context, app) {
          app.template = {
            container: 'template',
            application: 'pure',
            directory: ''
          }
          return context.state.createNewApp(app)
            .then((res) => {
              // Show error response
              if (res.response.status !== 201) {
                let doc = res.document
                let msg = doc.error || doc.message
                msg = `[${res.response.status}] ${msg}`
                let err = new Error(msg)
                context.dispatch('log', err)
                throw err
              }

              context.dispatch('log', `Created ${app.name}`)

              // Refresh containers list
              return context.dispatch('containers')
            })
        },
        'del-app': function (context, {container, application}) {
          return context.state.deleteApp(container, application)
            .then((res) => {
              // Show error response
              if (res.response.status !== 200) {
                let doc = res.document
                let msg = doc.error || doc.message
                msg = `[${res.response.status}] ${msg}`
                let err = new Error(msg)
                context.dispatch('log', err)
                throw err
              }

              context.dispatch('log', `Deleted ${application}`)

              // Refresh containers list
              return context.dispatch('containers')
            })
        },
        'app': function (context) {
          return context.state.getApplication()
            .then((doc) => {
              context.commit('app', doc)
            })
        },
        'list-files': function (context) {
          return context.state.getFiles()
            .then((list) => {
              context.commit('files', list)
            })
        },
        'list-pages': function (context) {
          return context.state.getPages()
            .then((list) => {
              context.commit('pages', list)
            })
        },
        'reload': function (context) {
          return context.dispatch('list-pages')
            .then(() => context.dispatch('list-files'))
        },
        'get-file-contents': function (context, item) {
          return context.state.getFileContents(item.url)
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
              if (file.editorView) {
                context.commit('editor-view', file.editorView)
              }
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
          context.commit('preview-refresh', true)
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
                let err = new Error(msg)
                context.dispatch('log', err)
                throw err
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
        'save-file': function (context, file) {
          if (!file) {
            file = context.state.current
          }
          let value = file.content
          return context.state.saveFile(file, value)
            .then((res) => {
              let doc = res.document
              if (res.response.status !== 200) {
                let msg = doc.error || doc.message
                msg = `[${res.response.status}] ${msg}`
                let err = new Error(msg)
                context.dispatch('log', err)
                throw err
              }
              if (doc.ok) {
                return context.dispatch('preview-refresh')
              }
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
                let err = new Error(msg)
                context.dispatch('log', err)
                throw err
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
        },
        'rename-file': function (context, {file, newName}) {
          return context.state.renameFile(file, newName)
            .then((res) => {
              let doc = res.document
              if (res.response.status !== 200) {
                let msg = doc.error || doc.message
                msg = `[${res.response.status}] ${msg}`
                let err = new Error(msg)
                context.dispatch('log', err)
                throw err
              }

              if (doc.ok) {
                context.dispatch('log', `Renamed ${file.url} to ${newName}`)
                // Update file data
                if (doc.file) {
                  file.name = doc.file.name
                  file.url = doc.file.url
                  file.uri = doc.file.uri
                }
                context.commit('preview-url', file.uri)
              }

              // return context.dispatch('reload')
            })
        }
      }
    })

    let r = this.router = new Router('home', true)
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
    r.add(/^apps\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+\/(files|pages|media|new|del)$/,
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
    r.add(/^(|home|apps|docs|edit|settings)$/, ['section'], (match) => {
      let section = match.map.section

      // Request with just the #
      if (section === '') {
        return r.replace('home', true)
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
              <a v-bind:class="{selected: currentView === 'media'}"
                @click="currentView = 'media'"
                title="Show media files">Media</a>
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
          <div class="column-drag" :class="{hidden: maximized}" @mousedown="resizeColumn"></div>
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
            return this.$store.state.columns.maximized
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
            var file = this.$store.state.current
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
            .catch((e) => console.error(e))
        },
        showNewFileView: function () {
          this.previousView = this.currentView
          this.currentView = 'new'
        },
        closeNewFileView: function () {
          this.currentView = this.previousView || 'pages'
        },
        resizeColumn: function (e) {
          this.$store.dispatch('resize-column', e)
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
              return this.$store.dispatch(
                'new-file', {name: this.fileName, template: this.template, action: action})
                .then(() => {
                  this.$parent.closeNewFileView()
                })
                .catch((e) => console.error(e))
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
        media: {
          template: `
            <div class="media-list">
            </div>`,
          computed: {
            list: function () {
              return this.$store.state.app.media
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
            <a @click="refresh(path)"
               title="Refresh preview"
               :class="{hidden: path == ''}">Refresh</a>
            <a
              @click="maximized = 'preview'"
              :class="{hidden: maximized === 'preview'}"
              title="Maximize">◩</a>
            <a
              @click="maximized = ''"
              :class="{hidden: maximized !== 'preview'}"
              title="Minimize">▣</a>
          </nav>
          <iframe :src="src" sandbox="allow-same-origin allow-scripts" class="publish-preview"></iframe>
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
            return this.$store.state.columns.maximized
          },
          set: function (val) {
            this.$store.commit('maximize-column', val)
          }
        },
        previewRefresh: function () {
          return this.$store.state.previewRefresh
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
          this.refresh(url)
        },
        previewRefresh: function (val) {
          if (val === true) {
            this.refresh(this.path)
          }
          this.$store.commit('preview-refresh', false)
        }
      },
      mounted: function () {
        // This catches the case when switching main views
        // and a refresh is needed
        if (this.url) {
          this.refresh(this.url)
        }
      },
      methods: {
        refresh (url) {
          let allowed = /\.(html?)$/
          if (!allowed.test(url)) {
            return
          }
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
          return document.location.origin + this.$store.state.app.url + (url || '')
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
                <a v-bind:class="{selected: currentView === 'file-editor', disabled: fileHidden}"
                  @click="currentView = 'file-editor'"
                  title="Show file editor">File</a>
                <a v-bind:class="{selected: currentView === 'data-editor', disabled: dataHidden}"
                  @click="currentView = 'data-editor'"
                  title="Show data editor">Data</a>
                <a v-bind:class="{selected: currentView === 'source-editor', disabled: hidden}"
                  @click="currentView = 'source-editor'"
                  title="Show source editor">Code</a>
                <a v-bind:class="{selected: currentView === 'visual-editor', disabled: hidden}"
                  @click="currentView = 'visual-editor'"
                  title="Show visual editor">Visual</a>
              </nav>
            </div>
          </div>
          <nav class="toolbar clearfix">
            <h2><span class="status-dirty" :class="{hidden: !isDirty}">✺</span>{{currentFile.name}}</h2>
            <a @click="save"
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
          <div class="column-drag" :class="{hidden: maximized}" @mousedown="resizeColumn">&nbsp;</div>
        </div>
      `,
      computed: {
        dirty: {
          get: function () {
            return this.currentFile.dirty
          },
          set: function (val) {
            if (this.changeGeneration > -1 && this.currentFile.document) {
              if (val === true && this.currentFile.document.isClean(this.changeGeneration)) {
                val = false
              }
            }
            this.$store.commit('current-file-dirty', val)
            this.isDirty = val
          }
        },
        maximized: {
          get: function () {
            return this.$store.state.columns.maximized
          },
          set: function (val) {
            this.$store.commit('maximize-column', val)
          }
        },
        dataHidden: function () {
          return this.hidden || !this.$store.state.current.page
        },
        fileHidden: function () {
          return !this.$store.state.hasFile()
        },
        hidden: function () {
          return !this.$store.state.hasFile() || this.$store.state.isDirectory()
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
          this.dirty = file.dirty
          if (file && file.dir) {
            this.currentView = 'file-editor'
          }
        }
      },
      data: function () {
        return {
          title: '',
          isDirty: false,
          changeGeneration: -1,
          codeMirror: null
        }
      },
      methods: {
        save: function (e) {
          if (e) {
            e.preventDefault()
          }
          this.$store.dispatch('save-file')
            .then(() => {
              this.dirty = false
              this.changeGeneration = this.currentFile.document.changeGeneration()
            })
            .catch((e) => console.error(e))
        },
        resizeColumn: function (e) {
          this.$store.dispatch('resize-column', e)
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
                  <form @submit="rename" class="rename">
                    <input type="text" name="fileName" v-model="newName" />
                    <div class="form-actions">
                      <input type="submit" name="Rename" value="Rename" />
                    </div>
                  </form>
                </section>
                <section>
                  <h3>Delete File</h3>
                  <p v-bind:class="{hidden: confirmDelete}">Danger zone: be careful!</p>
                  <div class="form-actions">
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
              confirmDelete: false,
              newName: ''
            }
          },
          computed: {
            file: function () {
              return this.$store.state.current
            }
          },
          created: function () {
            this.newName = this.file.url
          },
          methods: {
            rename: function (e) {
              e.preventDefault()
              return this.$store.dispatch('rename-file', {file: this.file, newName: this.newName})
                .catch((e) => console.error(e))
            },
            doDelete: function () {
              this.confirmDelete = false
              return this.$store.dispatch('delete-file', this.file)
            }
          },
          watch: {
            file: function (file) {
              this.newName = file.url
            }
          }
        },
        'data-editor': {
          computed: {
            pageDataJson: function () {
              return JSON.stringify(this.pageData, undefined, 2)
            },
            pageData: {
              get: function () {
                if (!this.$store.state.current || !this.$store.state.current.page) {
                  return {}
                }
                return this.$store.state.current.page.data
              },
              set: function (val) {
                //
              }
            }
          },
          render: function (h) {
            // We need recursion to render meta page data
            function list (target) {
              let isArr = Array.isArray(target)
              function it (o, fn) {
                if (Array.isArray(o)) {
                  o.forEach(fn)
                } else {
                  for (let k in o) {
                    fn(o[k], k)
                  }
                }
              }

              let el = h(isArr ? 'ol' : 'ul', null, [])
              it(target, (value, key) => {
                let li = h('li', {'data-type': typeof (value), 'data-key': key}, [])
                if (!isArr) {
                  let k = h('span', {class: 'data-key'}, ['' + key])
                  li.children.push(k)
                }
                if (Array.isArray(value) || value && typeof (value) === 'object') {
                  let nodes = list(value)
                  li.children.push(...nodes)
                } else {
                  let v = h('span', {class: 'data-value'}, ['' + value])
                  li.children.push(v)
                }
                el.children.push(li)
              })
              return [el]
            }

            let children = list(this.pageData)
            let el = h('div', {class: 'data-editor'}, children)
            return el
          }
        },
        'source-editor': {
          template: `<div class="source-editor">
              <div class="text-editor"></div>
            </div>
          `,
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
                this.setCodeMirror(file, {mode: this.getModeForMime(file.mime)})
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
              this.$parent.dirty = true
            },
            save: function () {
              this.$parent.save()
            },
            setCodeMirror: function (file, options) {
              options = options || {}
              let p = document.querySelector('.text-editor')

              if (this.mirror) {
                // TODO: verify listener is removed
                this.mirror.off('changes', this.changes)
                if (p.firstChild) {
                  p.removeChild(p.firstChild)
                }
              }
              if (file.document) {
                file.document.cm = null
              }
              this.mirror = CodeMirror(p, {
                value: file.document || file.content || '',
                mode: options.mode || 'htmlmixed',
                theme: options.theme || 'solarized dark',
                lineNumbers: true,
                keyMap: 'vim'
              })
              this.mirror.on('changes', this.changes)
              this.mirror.setOption('extraKeys', {
                'Ctrl-S': (cm) => {
                  this.save()
                },
                Tab: (cm) => {
                  var spaces = Array(cm.getOption('indentUnit') + 1).join(' ')
                  cm.replaceSelection(spaces)
                }
              })

              /*
              let wait
              let opts = {column: 80}
              let changing = false
              this.mirror.on('change', function (cm, change) {
                if (changing) return
                clearTimeout(wait)
                wait = setTimeout(function () {
                  changing = true
                  cm.wrapParagraphsInRange(change.from, CodeMirror.changeEnd(change), opts)
                  changing = false
                }, 200)
              })
              */

              // This gives us :w in vim mode
              CodeMirror.commands.save = (cm) => {
                this.save()
              }

              if (!file.document) {
                file.document = this.mirror.getDoc()
                this.$parent.changeGeneration = file.document.changeGeneration()
              }

              this.$parent.codeMirror = this.mirror
            }
          },
          mounted: function () {
            let file = this.currentFile
            // Handles setting file content when switching tabs
            if (file && file.mime) {
              this.setCodeMirror(file, {mode: this.getModeForMime(file.mime)})
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
          <app-alert></app-alert>
          <app-header></app-header>
          <app-main></app-main>
          <app-footer></app-footer>
        </main>
      `,
      store: this.store,
      components: {
        'app-alert': {
          template: `
            <div class="alert">
              <div class="background"></div>
              <div class="dialog">
                <h2>{{title}}</h2>
                <p v-if="message">{{message}}</p>
                <slot />
              </div>
            </div>
          `,
          data: function () {
            return {
              title: 'Alert',
              message: 'Are you sure?'
            }
          }
        },
        'app-header': {
          template: `
              <header class="clearfix">
                <nav class="home">
                  <a
                    @click="$store.dispatch('navigate', {href: 'edit'})"
                    :class="{selected: selectedView === 'edit', hidden: !this.$store.state.hasApplication()}"
                    title="Edit Current Application">{{name}}</a>
                </nav>
                <nav class="main">
                  <a
                    @click="$store.dispatch('navigate', {href: 'home'})"
                    class="home"
                    :class="{selected: selectedView === 'home'}"
                    title="Home page">Ꝏ </a>
                  <a
                    @click="$store.dispatch('navigate', {href: 'apps'})"
                    :class="{selected: selectedView === 'apps'}"
                    title="View and edit applications">Apps</a>
                  <a
                    @click="$store.dispatch('navigate', {href: 'docs'})"
                    :class="{selected: selectedView === 'docs'}"
                    title="Documentation">Docs</a>
                  <a
                    @click="$store.dispatch('navigate', {href: 'settings'})"
                    :class="{selected: selectedView === 'settings'}"
                    title="Settings">Settings</a>
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
            'home': {
              template: `
                <div class="content-main home">
                  <div class="scroll center">
                    <section class="document">
                      <h1><span>Ꝏ</span>&nbsp;Pageloop</h1>
                      <p class="leader">Collaborative realtime web based document manager</p>
                      <p>Designed for static websites pageloop enables you to create modern websites from lovingly crafted responsive templates in seconds.</p>
                      <p>To get started go to the <a href="#apps" title="Application manager">application manager</a> and choose from one of our beautiful, responsive templates for your next website.</p>
                      <p>You can use the application editor to customize your website using an intuitive drag and drop interface. Technical folks can edit the source code directly and you can export your files at any time if you need to.</p>
                      <p>For those working as a team you can assign semantic roles (such as designer, editor, translator etc) so that we can customize the interface to best fit your needs.</p>
                      <p>We suggest you read the <a href="#docs" title="Documentation & help">documentation</a> to learn more about how pageloop works.</p>
                    </section>
                  </div>
                </div>
              `
            },
            'apps': {
              template: `
                <div class="content-main">
                  <div class="content scroll">
                    <div class="new-app">
                      <h2>New Application</h2>
                      <form @submit="createApplication">
                        <p class="small">Choose an application name:</p>
                        <input type="text" name="name"
                          :value="applicationName" v-model="applicationName" />
                        <p class="small">Enter a publish URL:</p>
                        <input @change="urlChanged = true" type="text" name="url"
                          :value="applicationUrl" v-model="applicationUrl" />
                        <p class="small">Short description:</p>
                        <input type="text" name="description"
                          :value="applicationDescription" v-model="applicationDescription" />
                        <div class="form-actions">
                          <input type="submit" value="Create" />
                        </div>
                      </form>
                    </div>
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
                                <a class="name"
                                  @click="$store.dispatch('navigate', {href: linkify(container, app)})"
                                  :title="title(app, 'Edit')">Edit</a>
                                <a class="name"
                                  :href="linkify(container, app, true)"
                                  :title="title(app, 'Open')">Open</a>
                                <a v-if="!app.protected" class="name"
                                  @click="deleteApplication(container, app)"
                                  :title="title(app, 'Delete')">Delete</a>
                              </p>
                            </p>
                        </div>
                      </ul>
                    </div>
                  </div>
                </div>
              `,
              data: function () {
                return {
                  urlChanged: false,
                  applicationName: 'new-app',
                  applicationUrl: '/new-app',
                  applicationDescription: 'New application'
                }
              },
              watch: {
                applicationName: function (val) {
                  if (!this.urlChanged) {
                    this.applicationUrl = '/' + val
                  }
                }
              },
              computed: {
                list: function () {
                  return this.$store.state.containers
                }
              },
              methods: {
                createApplication: function (e) {
                  e.preventDefault()
                  let app = {}
                  if (this.applicationName) {
                    app.name = this.applicationName
                  }
                  if (this.applicationUrl) {
                    app.url = this.applicationUrl
                  }
                  if (this.applicationDescription) {
                    app.description = this.applicationDescription
                  }
                  this.$store.dispatch('new-app', app)
                    .then(() => {
                      return this.$store.dispatch('navigate', {href: `apps/user/${app.name}`})
                    })
                    .catch((e) => console.error(e))
                },
                deleteApplication: function (container, application) {
                  this.$store.dispatch('del-app', {container: container.name, application: application.name})
                    .catch((e) => console.error(e))
                  return false
                },
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
              let msg = this.$store.state.log.last
              if (msg && this.error) {
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
    this.state.setApplication(container, application)
    this.store.dispatch('log', `Loading app from ${this.state.url}`)
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
    window.onbeforeunload = (e) => {
      if (this.state.isDirty()) {
        return true
      }
    }
    this.ui()
    this.router.start()
  }
}

let app = new EditorApplication()
app.init()
