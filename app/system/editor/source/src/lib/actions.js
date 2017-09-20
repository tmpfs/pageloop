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
    'get-meta': function (context) {
      return context.state.client.getMeta()
        .then((res) => {
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }
          context.commit('meta', res.document)
        })
    },
    'reset-current-file': function (context, url) {
      context.state.current = context.state.app.defaultFile
      context.commit('preview-blank')
      context.commit('editor-view', 'welcome')
    },
    'error': function (context, err) {
      // console.log('showing error')
      // Log the error
      context.dispatch('log', err)
      // Notify the user
      context.state.notify({error: err})
      // Return the error to pass back to calling code
      return err
    },
    'load': function (context, {container, application}) {
      context.state.setApplication(container, application)
      context.dispatch('log', `Loading app ${application}`)
      return context.dispatch('app')
        .then(() => context.dispatch('list-files'))
        .then(() => context.dispatch('list-pages'))
        .then(() => {
          context.commit('sidebar-view', 'pages')
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
        .then((res) => {
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }
          context.commit('containers', res.document.containers)
        })
    },
    'run-task': function (context, {app, task}) {
      const container = context.state.container
      return context.state.client.runTask(container, app.name, task)
        .then((res) => {
          // Show error response
          if (res.response.status !== 202) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }
          context.state.notify({title: 'Task Info', message: `Task ${task} started in ${app.name}`})
          return res
        })
    },
    'new-app': function (context, app) {
      return context.state.client.createApp(app)
        .then((res) => {
          // Show error response
          if (res.response.status !== 201) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }

          context.state.notify({title: 'App Info', message: `Created application ${app.name}`})

          // Refresh containers list
          return context.dispatch('containers')
        })
    },
    'edit-app': function (context, {container, application}) {
      let href = `apps/${container.name}/${application.name}`
      context.state.log.add({level: 'Info', message: `Edit application ${application.name}`})
      return context.dispatch('navigate', {href: href})
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
      const container = context.state.container
      const application = context.state.application
      return context.state.client.getApplication(container, application)
        .then((res) => {
          // Show error response
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }
          context.commit('app', res.document)
        })
    },
    'list-files': function (context) {
      const container = context.state.container
      const application = context.state.application
      return context.state.client.getFiles(container, application)
        .then((res) => {
          // Show error response
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }
          context.commit('files', res.document)
        })
    },
    'list-pages': function (context) {
      const container = context.state.container
      const application = context.state.application
      return context.state.client.getPages(container, application)
        .then((res) => {
          // Show error response
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }
          context.commit('pages', res.document)
        })
    },
    'reload': function (context) {
      return context.dispatch('list-pages')
        .then(() => context.dispatch('list-files'))
    },
    'get-file-contents': function (context, file) {
      const container = context.state.container
      const application = context.state.application
      return context.state.client.getFileSource(container, application, file, true)
        .then((res) => {
          return file.binary ? res.blob() : res.text()
        })
    },
    'open-file': function (context, file) {
      context.state.log.add(
        {level: 'Info', message: `Open file ${file.url}`})
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
          context.commit('preview-change', file)
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
      const container = context.state.container
      const application = context.state.application
      return context.state.client.createFile(container, application, name, template)
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
      const container = context.state.container
      const application = context.state.application
      return context.state.client.saveFile(container, application, file, value)
        .then((res) => {
          let doc = res.document
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }

          // file.content = value

          context.state.log.add({level: 'Info', message: `Saved file ${file.url}`})

          // Currently YAML is in the source code and
          // can be edited directly we need to sync
          // the data in case it changed
          if (doc && doc.data) {
            file.data = doc.data
          }
          return context.dispatch('preview-refresh')
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
      const container = context.state.container
      const application = context.state.application
      return context.state.client.deleteFiles(container, application, files)
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
          // Current selected file is in deleted list
          if (~files.indexOf(context.state.current)) {
            // TODO: remove file from window location if already selected!
            context.dispatch('reset-current-file')
          }
          return context.dispatch('reload')
            .then(() => {
              // Deleted the currently selected file
              if (context.state.current && ~files.indexOf(context.state.current)) {
                context.dispatch('reset-current-file')
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
      const container = context.state.container
      const application = context.state.application
      return context.state.client.moveFile(container, application, file, newName)
        .then((res) => {
          let doc = res.document
          if (res.response.status !== 200) {
            return context.dispatch('error', error(res))
              .then((err) => {
                throw err
              })
          }

          context.state.notify({title: 'File Info', message: `Renamed ${file.url} to ${newName}`})

          context.state.app.updateFile(file, doc)

          const view = context.state.sidebar.view
          router.replace(context.state.getAppHref(view, file.url), false)

          context.commit('preview-change', file)
        })
    },
    'upload': function (context, info) {
      const container = context.state.container
      const application = context.state.application
      return context.state.transfer.upload(container, application, info.files)
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
