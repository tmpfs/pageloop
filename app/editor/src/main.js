// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'

import Vuex from '../static/vendor/vuex-2.3.1'

Vue.use(Vuex)

import Router from './lib/router'
import State from './lib/state'

let r = new Router()

let store = new Vuex.Store({
  state: new State(),
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
    },
    'alert-show': function (state, details) {
      for (let k in details) {
        state.alert[k] = details[k]
      }
      state.alert.visible = true
    },
    'alert-hide': function (state, details) {
      state.alert.visible = false
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

          context.state.notify({title: 'App Info', message: `Created ${app.name}`})

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

          context.state.notify({title: 'App Info', message: `Deleted ${application}`})

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

          context.state.notify({title: 'File Info', message: `Created ${name}`})

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

          context.state.notify({title: 'File Info', message: `Deleted ${file.name}`})

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
        })
    }
  }
})

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: 'main',
  template: '<App/>',
  components: { App },
  store
})
