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

    this.concurrentTransfers = 3
    this.currentTransfer = null
  }

  upload () {
    if (this.transfers.length) {
      this.currentTransfer = [this.transfers[0]]

      let amount = this.transfers.length / this.concurrentTransfers
      if (this.transfers.length % this.concurrentTransfers !== 0) {
        amount++
      }

      let chunks = []
      let i, ind
      for (i = 0; i < amount; i++) {
        ind = i * this.concurrentTransfers
        chunks.push(this.transfers.slice(ind, ind + this.concurrentTransfers))
      }

      const transfer = (chunk) => {
        return new Promise((resolve, reject) => {
          let loaded = 0
          chunk.forEach((file) => {
            this.client.upload(file)
              .then(() => {
                loaded++
                if (loaded === chunk.length) {
                  // More chunks to process
                  if (chunks.length) {
                    this.currentTransfer = chunks.shift()
                    transfer(this.currentTransfer)
                  // All done, upload completed
                  } else {
                    resolve()
                  }
                }
              })
              .catch(reject)
          })
        })
      }
      this.currentTransfer = chunks.shift()
      return transfer(this.currentTransfer)
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
