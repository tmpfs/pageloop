function error (res) {
  let doc = res.document
  let msg = doc.error || doc.message
  msg = `${msg}`
  let err = new Error(msg)
  err.status = res.response.status
  err.response = res
  err.document = doc
  return err
}

function Actions (router) {
  return {
    'error': function (context, err) {
      // Log the error
      context.dispatch('log', err)
      // Notify the user
      context.state.notify({error: err})
      // Return the error to pass back to calling code
      return err
    },
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
        .catch((err) => {
          return context.dispatch('error', err)
            .then((err) => {
              throw err
            })
        })
    },
    'resize-column': function (context, e) {
      context.state.editor.columns.startDrag(e)
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
      return context.state.client.createNewApp(app)
        .then((res) => {
          // Show error response
          if (res.response.status !== 201) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
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
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }

          let state = context.state

          if (state.hasApplication()) {
            if (state.container === container && state.application === application) {
              context.commit('clear-app')
            }
          }

          state.notify({title: 'App Info', message: `Deleted ${application}`})

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
          return item.binary ? res.blob() : res.text()
        })
    },
    'open-file': function (context, file) {
      return context.dispatch('get-file-contents', file)
        .then((content) => {
          if (!file.binary) {
            file.content = content
          } else {
            file.blob = content
          }
          context.commit('current-file', file)
          if (file.editorView) {
            context.commit('editor-view', file.editorView)
          }
          context.commit('preview-url', file)
          if (context.state.editor.view === 'welcome') {
            if (!file.binary) {
              context.commit('editor-view', context.state.editor.defaultView)
            } else {
              context.commit('editor-view', 'visual-editor')
            }
          }
          context.commit('selected-file', file)
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
    'go-media': function (context, {filter, file}) {
      let href = context.state.getAppHref(filter, file.url)
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
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
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
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }
          if (doc.ok) {
            // Currently YAML is in the source code and
            // can be edited directly we need to sync
            // the data in case it changed
            if (doc.file && doc.file.data) {
              file.data = doc.file.data
            }
            return context.dispatch('preview-refresh')
          }
        })
    },
    'delete-files': function (context, files) {
      let list = context.state.app.files
      let ctx = context.state.sidebar.view
      if (ctx === 'pages') {
        list = context.state.app.pages
      }

      // TODO: restore nearest neighbour selection

      // let index = list.indexOf(file)
      let len = list.length
      return context.state.client.deleteFiles(files)
        .then((res) => {
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }

          const urls = files.map((f) => {
            return f.url
          })

          context.state.notify({title: 'File Info', message: `Deleted ${urls.join(', ')}`})
          return context.dispatch('reload')
            .then(() => {
              // Deleted the currently selected file
              if (context.state.current && ~files.indexOf(context.state.current)) {
                context.commit('reset-current-file')
                context.commit('editor-view', 'welcome')
              }

              if (len <= 1) {
              /*
              } else if (index > -1) {
                // select next nearest file
                let neighbour = list[index - 1] || list[index + 1]
                let href = context.state.getAppHref(ctx, neighbour.url)
                return context.dispatch('navigate', {href: href})
              */
              }
            })
        })
    },
    'rename-file': function (context, {file, newName}) {
      return context.state.client.renameFile(file, newName)
        .then((res) => {
          let doc = res.document
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }

          if (doc.ok) {
            context.state.notify({title: 'File Info', message: `Renamed ${file.url} to ${newName}`})
            // Update file data
            if (doc.file) {
              file.name = doc.file.name
              file.url = doc.file.url
              file.uri = doc.file.uri
            }
            context.commit('preview-url', file)
          }
        })
    },
    'upload': function (context, info) {
      return context.state.transfer.upload(info.files)
        .then((transfers) => {
          // Reload file list for the moment
          return context.dispatch('reload')
            .then(() => transfers)
        })
    },
    'list-templates': function (context) {
      return context.state.client.listTemplates()
        .then(({response, document}) => {
          context.commit('templates', document)
        })
    }
  }
}

export default Actions
