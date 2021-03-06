// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'

import Vuex from 'vuex'

Vue.config.productionTip = false
Vue.use(Vuex)

import Router from './lib/router'
import State from './lib/state'
import Mutations from './lib/mutations'
import Actions from './lib/actions'
import Routes from './routes'

class Application {

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
    // Register a global custom directive called v-focus
    Vue.directive('focus', {
      inserted: function (el) {
        el.focus()
      }
    })
    /* eslint-disable no-new */
    new Vue({
      el: 'main',
      template: '<App/>',
      components: { App },
      store
    })
  }

  init () {
    window.onbeforeunload = (e) => {
      if (this.store.state.isDirty()) {
        return true
      }
    }
    this.ui(this.store)
    this.router.start()
  }
}

let app = new Application()
app.init()
