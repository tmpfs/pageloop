import Application from './application'
import ApiClient from './api'

import MainState from './state/main'
import SidebarState from './state/sidebar'
import EditorState from './state/editor'

import Flash from './state/flash'
import Log from './state/log'

class State {
  constructor () {
    this.client = this.defaultClient = new ApiClient()
    this.containers = []
    this.setApplication('', '')

    this.main = new MainState()
    this.sidebar = new SidebarState()
    this.editor = new EditorState()

    this.previewUrl = ''
    this.previewRefresh = false

    this.log = new Log()
    this.flash = new Flash()

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
    this.app = new Application()

    this.client = this.defaultClient

    // File upload transfers
    this.transfers = []

    this.concurrentTransfers = 3
    this.currentTransfer = []
  }

  upload () {
    if (this.transfers.length) {
      let amount = Math.floor(this.transfers.length / this.concurrentTransfers)
      if (this.transfers.length % this.concurrentTransfers !== 0) {
        amount++
      }

      let chunks = []
      let i, ind, len
      for (i = 0; i < amount; i++) {
        ind = i * this.concurrentTransfers
        len = Math.min(this.transfers.length, ind + this.concurrentTransfers)
        chunks.push(this.transfers.slice(ind, len))
      }

      // Transfer a single chunk
      const transfer = (chunk, done) => {
        return new Promise((resolve, reject) => {
          let loaded = 0
          chunk.forEach((file) => {
            this.client.upload(file).then((file) => {
              loaded++
              if (loaded === chunk.length) {
                // Process next chunk
                if (chunks.length) {
                  this.currentTransfer = chunks.shift()
                  resolve(transfer(this.currentTransfer, done))
                // All done, upload completed
                } else {
                  done(this.transfers)
                }
              }
            })
            .catch(reject)
          })
        })
      }
      this.currentTransfer = chunks.shift()
      return new Promise((resolve, reject) => {
        transfer(this.currentTransfer, (files) => {
          this.transfers = []
          this.currentTransfer = []
          resolve(files)
        })
        .catch((err) => {
          this.transfers = []
          this.currentTransfer = []
          reject(err)
        })
      })
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
