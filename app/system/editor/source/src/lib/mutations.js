const Mutations = {
  'meta': function (state, meta) {
    state.meta = meta
  },
  'flash': function (state, message) {
    state.flash.message = message
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
    state.main.view = view
  },
  'sidebar-view': function (state, view) {
    state.sidebar.view = view
  },
  'editor-view': function (state, view) {
    state.editor.view = view
    if (state.hasFile()) {
      state.current.editorView = view
    }
  },
  'current-file': function (state, file) {
    if (!file.editorView) {
      if (file.binary) {
        file.editorView = state.editor.defaultBinaryView
      } else {
        file.editorView = state.editor.defaultView
      }
    }
    state.current = file
  },
  'current-file-dirty': function (state, val) {
    state.current.dirty = val
  },
  'selected-file': function (state, file) {
    if (!state.sidebar.selection.length) {
      const view = state.sidebar.view
      if (view) {
        const list = state.app[view]
        // TODO: check this works for pages too
        if (~list.indexOf(file)) {
          state.sidebar.selection = [file]
        }
      }
    }
  },
  'preview-blank': function (state, file) {
    state.preview.blank = !state.preview.blank
  },
  'preview-change': function (state, file) {
    state.preview.file = file
  },
  'preview-refresh': function (state, toggle) {
    state.preview.refresh = toggle
  },
  'maximize-column': function (state, info) {
    // Maximizing
    if (info) {
      state.editor.columns.maximize(info)
    // Minimizing
    } else {
      state.editor.columns.minimize(state.editor.columns.maximized)
    }
  },
  'reset-column': function (state) {
    state.editor.columns.reset()
    if (state.editor.columns.maximized) {
      state.editor.columns.minimize(state.editor.columns.maximized)
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
  },
  'transfers': function (state, {files, dir}) {
    let i, f, url, exists
    const list = []
    for (i = 0; i < files.length; i++) {
      f = files[i]
      if (!dir) {
        url = '/' + f.name
      } else {
        url = dir.replace(/\/$/, '') + '/' + f.name
      }
      exists = state.app.getFileByUrl(url)
      list.push({
        name: f.name,
        size: f.size,
        upload: f,
        complete: false,
        dir: dir || '',
        exists: exists,
        info: {
          ratio: 0
        }
      })
    }
    state.transfer.transfers = list
  },
  'dismiss-hint': function (state, id) {
    state.settings.set(id, false)
  },
  'templates': function (state, templates) {
    state.templates = templates
  },
  'services': function (state, services) {
    state.services.list = services
  },
  'clear-local-storage': function (state) {
    state.settings.reset()
  },
  'new-app-view': function (state, view) {
    state.newApp.view = view
  },
  'app-list-selected': function (state, app) {
    state.appList.selected = app
  },
  'main-menu-visible': function (state, visibility) {
    state.settings.mainMenuVisible = visibility
  }
}

export default Mutations
