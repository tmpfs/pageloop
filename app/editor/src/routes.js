function Routes (router, store) {
  let state = store.state

  router.add(/^apps\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+\/(pages|files)\/(.*)$/,
    ['section', 'container', 'application', 'action'],
    (match) => {
      let href = '/' + match.parts.slice(4).join('/')
      let container = match.map.container
      let application = match.map.application
      let action = match.map.action
      let file

      function findAndOpen (href) {
        let arr = state.app.files
        if (action === 'pages') {
          arr = state.app.pages
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
      if (container !== state.container || (container === state.container && application !== state.application)) {
        store.dispatch('load', {container: match.map.container, application: match.map.application})
          .then(() => {
            return trigger()
          })
      } else {
        return trigger()
      }
    })

  router.add(/^apps\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+\/(files|pages|media|new|del)$/,
    ['section', 'container', 'application', 'action'],
    (match) => {
      let container = match.map.container
      let application = match.map.application
      let action = match.map.action

      // Need to load application data
      if (container !== state.container || (container === state.container && application !== state.application)) {
        store.dispatch('load', {container: match.map.container, application: match.map.application})
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

  router.add(/^apps\/[a-zA-Z0-9-]+\/[a-zA-Z0-9-]+$/,
    ['section', 'container', 'application'],
    (match) => {
      store.dispatch('load', {container: match.map.container, application: match.map.application})
        .then(() => {
          let index = store.state.getIndexFile()
          if (index) {
            let href = match.href + '/files' + index.url
            // Redirect to index page if there is one
            return router.replace(href, true)
          }
          store.commit('reset-current-file')
          store.commit('main-view', 'edit')
          store.commit('editor-view', 'welcome')
        })
    })
  router.add(/^(|home|apps|docs|edit|settings)$/, ['section'], (match) => {
    let section = match.map.section

    // Request with just the #
    if (section === '') {
      return router.replace('home', true)
    } else if (section === 'apps') {
      return store.dispatch('containers')
        .then(() => {
          console.log('loaded containers!!!')
          store.commit('main-view', section)
        })
    } else if (section === 'edit') {
      if (state.hasApplication()) {
        return router.replace('apps/' + state.container + '/' + state.application, true)
      } else {
        // no app being edited redirect to apps list
        return router.replace('apps', true)
      }
    }
    store.commit('main-view', section)
  })

  router.add(/^404$/, (match) => {
    store.commit('main-view', 'not-found')
  })

  router.add(/.*/, (match) => {
    store.commit('flash', router.hash)
    router.replace('404', true)
  })
}

export default Routes
