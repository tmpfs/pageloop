function Actions (router) {
  return {
    'load': function (context, {container, application}) {
      context.state.setApplication(container, application)
      context.dispatch('log', `Loading app from ${context.state.url}`)
      return context.dispatch('app')
        .then(() => context.dispatch('list-files'))
        .then(() => context.dispatch('list-pages'))
        .then(() => {
          context.commit('sidebar-view', 'pages')
          context.dispatch('log', 'Done')
        })
        .catch((err) => context.dispatch('log', err))
    },
    'resize-column': function (context, e) {
      context.state.columns.startDrag(e)
    },
    'log': function (context, message) {
      context.commit('log', message)
    },
    'navigate': function (context, request) {
      return router.navigate(request.href, request.state)
    },
    'containers': function (context) {
      return context.state.client.getContainers()
        .then((list) => {
          context.commit('containers', list)
        })
    },
    'new-app': function (context, app) {
      app.template = {
        container: 'template',
        application: 'pure'
      }
      return context.state.client.createNewApp(app)
        .then((res) => {
          // Show error response
          if (res.response.status !== 201) {
            let doc = res.document
            let msg = doc.error || doc.message
            msg = `[${res.response.status}] ${msg}`
            let err = new Error(msg)
            context.dispatch('log', err)
            context.state.notify({error: err})
            throw err
          }

          context.state.notify({title: 'App Info', message: `Created ${app.name}`})

                  // Refresh containers list
          return context.dispatch('containers')
        })
    },
    'del-app': function (context, {container, application}) {
      return context.state.client.deleteApp(container, application)
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
      return context.state.client.getApplication()
        .then((doc) => {
          context.commit('app', doc)
        })
    },
    'list-files': function (context) {
      return context.state.client.getFiles()
        .then((list) => {
          context.commit('files', list)
        })
    },
    'list-pages': function (context) {
      return context.state.client.getPages()
        .then((list) => {
          context.commit('pages', list)
        })
    },
    'reload': function (context) {
      return context.dispatch('list-pages')
        .then(() => context.dispatch('list-files'))
    },
    'get-file-contents': function (context, item) {
      return context.state.client.getFileContents(item.url)
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
      return context.state.client.createNewFile(name, template)
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
      return context.state.client.saveFile(file, value)
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
      return context.state.client.deleteFile(file)
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
                context.commit('reset-current-file')
                context.commit('editor-view', 'welcome')
              } else if (index > -1) {
                // select next nearest file
                let neighbour = list[index - 1] || list[index + 1]
                let href = context.state.getAppHref(ctx, neighbour.url)
                return context.dispatch('navigate', {href: href})
              }
            })
        })
    },
    'rename-file': function (context, {file, newName}) {
      return context.state.client.renameFile(file, newName)
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
}

export default Actions
