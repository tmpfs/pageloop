// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'

import Vuex from '../static/vendor/vuex-2.3.1'

Vue.config.productionTip = false
Vue.use(Vuex)

import Router from './lib/router'
import State from './lib/state'
import Mutations from './lib/mutations'
import Actions from './lib/actions'
import Routes from './routes'

class EditorApplication {

  constructor () {
    let router = new Router('home', true)
    let store = new Vuex.Store({
      state: new State(),
      mutations: Mutations,
      actions: Actions(router)
    })

    Routes(router, store)

    this.store = store
    this.router = router
  }

  ui (store) {
    /* eslint-disable no-new */
    new Vue({
      el: 'main',
      template: '<App/>',
      components: { App },
      store
    })
  }

  init () {
    this.ui(this.store)
    this.router.start()
  }
}

let app = new EditorApplication()
app.init()
