import ApiClient from './client'
import Application from './application'

import MainState from './state/main'
import SidebarState from './state/sidebar'
import EditorState from './state/editor'
import PreviewState from './state/preview'

import Meta from './state/meta'
import Notifier from './state/notifier'
import Alert from './state/alert'
import Flash from './state/flash'
import Hints from './state/hints'

import Transfer from './transfer'
import Log from './state/log'
import Activity from './state/activity'

import {KeyManager} from './keymap'

class State {
  constructor () {
    this.keymap = new KeyManager()

    this.meta = new Meta()

    this.main = new MainState()
    this.sidebar = new SidebarState()
    this.editor = new EditorState()
    this.preview = new PreviewState()

    this.notifier = new Notifier()
    this.alert = new Alert()
    this.flash = new Flash()
    this.hints = new Hints()

    this.transfer = new Transfer()
    this.log = new Log()
    this.network = new Log()

    // We use the defaultClient when no application
    // is selected
    this.client = this.defaultClient = new ApiClient()
    this.defaultClient.log = this.network
    this.client.log = this.network

    this.containers = []
    this.templates = []
    this.setApplication('', '')

    this.activity = new Activity(this.log)
  }

  getContainerByName (name) {
    for (let i = 0; i < this.containers.length; i++) {
      if (name === this.containers[i].name) {
        return this.containers[i]
      }
    }
  }

  get client () {
    return this._client
  }

  set client (val) {
    this._client = val
    this.transfer.client = val
  }

  notify (info, del) {
    if (!del) {
      this.activity.add(info)
    }
    return this.notifier.notify(info, del)
  }

  getAppHref (...args) {
    let p = ['apps', this.container, this.application]

    // TODO: ensure we never get passed undefined / null etc
    // TODO: and remove this call to filter()
    args = args.filter((val) => {
      return val
    })

    args = args.map((val) => {
      return val.replace(/^\//, '')
    })
    p.push(...args)
    return p.join('/')
  }

  isDirty () {
    if (this.app) {
      return this.app.isDirty()
    }
    return false
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

  clearApplication () {
    this.container = ''
    this.application = ''

    // current application
    this.app = new Application()

    // Reset to default client
    this.client = this.defaultClient
  }

  setApplication (container, application) {
    this.clearApplication()
    this.app.owner = container

    this.container = container
    this.application = application

    // Set up new API client
    this.client = new ApiClient(container, application)
    this.client.log = this.network
  }

  get current () {
    return this.app.current
  }

  set current (file) {
    if (file) {
      // TODO: get the server to send the file?
      let pages = this.app.pages || []
      for (let i = 0; i < pages.length; i++) {
        if (pages[i].url === file.url) {
          file.page = pages[i]
          break
        }
      }
    }
    if (!file) {
      file = this.app.defaultFile
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
}

export default State
