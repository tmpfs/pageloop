import Log from './log'
import ColumnManager from './columns'
import ApiClient from './api'

class State {
  constructor () {
    this.client = this.defaultClient = new ApiClient()
    this.containers = []
    this.setApplication('', '')

    this.mainView = ''
    this.sidebarView = ''
    this.editorView = ''
    this.defaultEditorView = 'code-editor'

    this.defaultFile = {content: ''}

    this.previewUrl = ''
    this.previewRefresh = false

    this.log = new Log()

    this._flash = undefined

    // State for edit mode columns
    this.columns = new ColumnManager()

    this.alert = {
      visible: false,
      title: 'Alert',
      message: '',
      note: '',
      ok: function noop () {}
    }

    this.notifications = []
  }

  getAppHref (...args) {
    let p = ['apps', this.container, this.application]

    // TODO: ensure we never get passed undefined / null etc
    // TODO: and remove this
    args = args.filter((val) => {
      return val
    })

    args = args.map((val) => {
      return val.replace(/^\//, '')
    })
    p.push(...args)
    return p.join('/')
  }

  notify (info, del) {
    if (del) {
      for (let i = 0; i < this.notifications.length; i++) {
        if (info === this.notifications[i]) {
          this.notifications.splice(i, 1)
          break
        }
      }
      return
    }

    info.reveal = true

    this.notifications.unshift(info)
  }

  get flash () {
    let f = this._flash
    this._flash = undefined
    return f
  }

  set flash (msg) {
    this._flash = msg
  }

  isDirty () {
    if (this.app) {
      for (let i = 0; i < this.app.files.length; i++) {
        if (this.app.files[i].dirty) {
          return true
        }
      }
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
    this.app = {
      url: '',
      identifier: '',
      owner: '',
      pages: [],
      files: [],
      // current selected file
      current: this.defaultFile
    }

    this.client = this.defaultClient

    // File upload transfers
    this.transfers = []
  }

  upload () {
    if (this.transfers.length) {
      return this.client.upload(this.transfers[0])
    }
  }

  setApplication (container, application) {
    this.clearApplication()
    this.app.owner = container

    this.container = container
    this.application = application

    // Set up new API client
    this.client = new ApiClient(container, application)
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
