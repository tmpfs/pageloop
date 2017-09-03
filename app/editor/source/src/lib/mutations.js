const Mutations = {
  'flash': function (state, message) {
    state.flash = message
  },
  'log': function (state, message) {
    state.log.add(message)
  },
  'containers': function (state, list) {
    state.containers = list
  },
  'app': function (state, app) {
    // merge properties
    for (let k in app) {
      state.app[k] = app[k]
    }
    state.app.identifier = state.app.owner + ' / ' + state.app.name
  },
  'clear-app': function (state) {
    state.clearApplication()
  },
  'files': function (state, list) {
    state.app.files = list
  },
  'pages': function (state, list) {
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
  'reset-column': function (state) {
    state.columns.reset()
    if (state.columns.maximized) {
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
}

export default Mutations
