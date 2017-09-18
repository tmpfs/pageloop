// Singleton websocket client
let socket

// Message identifier counter
let id = 0

function getDefaultOptions () {
  return {method: 'GET'}
}

// REST API endpoint
const API = '/api/'
// Websocket endpoint
const WS = '/ws/'

// Maps RPC function names to REST request URLs
const URLS = {
  'Core.Meta': function () {
    return API
  },
  'Core.Stats': function () {
    return API + 'stats/'
  },
  'Container.List': function () {
    return API + 'apps/'
  },
  'Template.ReadApplications': function () {
    return API + 'templates/'
  },
  'Jobs.ReadActiveJobs': function () {
    return API + 'jobs/'
  }
}

const OPTIONS = {
  'Core.Meta': getDefaultOptions,
  'Core.Stats': getDefaultOptions,
  'Container.List': getDefaultOptions,
  'Template.ReadApplications': getDefaultOptions,
  'Jobs.ReadActiveJobs': getDefaultOptions
}

class RpcRequest {
  constructor (id, method, params) {
    this.id = id
    this.method = method
    if (params && !Array.isArray(params)) {
      params = [params]
    }
    this.params = params || []
  }
}

class Request {

  // Translate a JSON RPC request to a standard HTTP request
  static translate (rpc) {
    const o = {}
    // console.log('translate request: ' + rpc.method)
    // console.log(rpc)
    o.url = URLS[rpc.method](rpc)
    o.options = OPTIONS[rpc.method](rpc)
    return o
  }

  // Get a JSON RPC request object.
  static rpc (method, params) {
    return new RpcRequest(++id, method, params)
  }
}

class SocketConnection {
  constructor () {
    this.url = document.location.origin.replace(/^http/, 'ws') + WS
    this.protocols
    this.opts
    this._conn
    this._listeners = []
  }

  get connected () {
    return this._conn && this._conn.readyState === WebSocket.OPEN
  }

  connect () {
    this._conn = new WebSocket(this.url, this.protocols, this.opts)

    console.log(this._conn)

    this._conn.onopen = () => {
      // this.send({message: 'foo'})
    }

    this._conn.onmessage = (e) => {
      console.log(e)
      if (e.data) {
        let doc
        try {
          doc = JSON.parse(e.data)
        } catch (e) {
          throw e
        }
        if (doc.id && this._listeners[id]) {
          this._listeners[id](doc)
          delete this._listeners[id]
        }
      }
    }

    this._conn.onerror = (err) => {
      // TODO: log this error
      console.error(err)
    }

    this._conn.onclose = () => {
      this.cleanup()
      console.log('socket conn closed')
    }
  }

  cleanup () {
    this._conn.onopen = null
    this._conn.onmessage = null
    this._conn.onerror = null
    this._conn.onclose = null
    this._conn = null
  }

  // Send a JSON payload and ignore any response
  send (payload) {
    if (this.connected) {
      this._conn.send(JSON.stringify(payload))
    }
  }

  request (payload) {
    if (this.connected) {
      return new Promise((resolve, reject) => {
        this._listeners[payload.id] = (response) => {
          // console.log('listener called')
          console.log(response)
          const res = {
            status: response.status,
            id: response.id,
            transport: 'ws://json-rpc'}
          const doc = response.error || response.result
          resolve({response: res, document: doc})
        }
        this._conn.send(JSON.stringify(payload))
      })
    }
  }
}

class ApiClient {
  constructor (container, application) {
    this.host = ''
    this.api = API
    this.container = container
    this.application = application
    this.url = `${this.api}apps/${container}/${application}/`
    this.raw = `/apps/raw/${container}/${application}`
    // should be injected
    this.log = null
    this.socket = this.connect()

    this.useWebsocket = true
  }

  // Singleton websocket
  connect () {
    if (!socket) {
      socket = new SocketConnection()
      socket.connect()
    }
    return socket
  }

  preflight (url, opts) {
    if (this.log) {
      return this.log.add({level: opts.method || 'GET', message: url})
    }
  }

  postflight (log, res) {
    if (log) {
      const err = !/^20(0|1|2)$/.test('' + res.status)
      const url = res.url.replace(document.location.origin, '')
      log.add({level: res.status, message: url, error: err})
    }
  }

  // Perform an API request and assume a JSON response.
  request (url, opts) {
    const log = this.preflight(url, opts)
    return fetch(url, opts)
      .then((res) => {
        this.postflight(log, res)
        if (opts.raw) {
          return res
        }
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  rpc (req) {
    // TODO: implement API requests over websockets
    if (this.useWebsocket && this.socket.connected) {
      return this.socket.request(req)
    }

    const {url, options} = Request.translate(req)
    console.log('translated: ' + url)
    console.log(options)
    return this.request(url, options)
  }

  upload (file) {
    return new Promise((resolve, reject) => {
      let u = this.url + 'files'
      let dir = file.dir
      if (dir) {
        u += dir
      }

      if (!/\/$/.test(u)) {
        u += '/'
      }

      u += file.name

      // TODO: log file uploads

      const method = file.exists ? 'POST' : 'PUT'

      // Need to use XHR for upload progress :(
      let xhr = new XMLHttpRequest()
      xhr.open(method, u, true)

      if (file.exists) {
        xhr.setRequestHeader('Content-Type', file.exists.mime)
      }

      xhr.upload.onprogress = (e) => {
        if (e.lengthComputable) {
          let ratio = (e.loaded / e.total)
          file.info.ratio = ratio
        }
      }

      xhr.onload = function (e) {
        const doc = JSON.parse(this.responseText)
        if (this.status !== 201 && this.status !== 200) {
          return reject(
            new Error(`Upload failed for ${file.name}: ${doc.error || doc.message}`))
        }
        file.complete = true

        // File object returned by the server
        file.handle = doc

        // Set a timeout before completion so
        // progress preloaders are visible on fast
        // uploads
        setTimeout(() => {
          resolve(file)
        }, 3000)
      }

      xhr.onerror = function (err) {
        reject(err)
      }

      xhr.send(file.upload)
    })
  }

  // Get meta version information
  getVersion () {
    return this.rpc(Request.rpc('Core.Meta'))
  }

  // Get server statistics
  getStats () {
    return this.rpc(Request.rpc('Core.Stats'))
  }

  // Combine version meta information with server statistics
  getMeta () {
    return this.getVersion()
      .then(({response, document}) => {
        if (response.status !== 200) {
          return {response: response, document: document}
        }
        const meta = {info: document}
        return this.getStats()
          .then(({response, document}) => {
            meta.stats = document
            return {response: response, document: meta}
          })
      })
  }

  // List application templates
  listTemplates () {
    return this.rpc(Request.rpc('Template.ReadApplications'))
  }

  getContainers () {
    return this.rpc(Request.rpc('Container.List'))
  }

  runTask (app, task) {
    const url = this.url + `tasks/${task}`
    const opts = {
      method: 'PUT'
    }
    return this.request(url, opts)
  }

  createNewApp (app) {
    const url = this.api + 'apps/user/'
    const opts = {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(app)
    }
    return this.request(url, opts)
  }

  getApplication () {
    const url = this.url
    return this.request(url, {})
  }

  getPages () {
    const url = this.url + 'pages/'
    return this.request(url, {})
  }

  getFiles () {
    const url = this.url + 'files/'
    return this.request(url, {})
  }

  getFileContents (pathname) {
    const url = this.raw + pathname
    return this.request(url, {raw: true})
  }

  deleteFiles (files) {
    const urls = files.map((f) => {
      return f.url
    })
    const url = this.url + 'files'
    const opts = {
      method: 'DELETE',
      body: JSON.stringify(urls),
      headers: {
        'Content-Type': 'application/json; charset=utf-8'
      }
    }
    opts.headers['Content-Length'] = opts.body.length
    return this.request(url, opts)
  }

  renameFile (file, newName) {
    const url = this.url + 'files' + file.url
    const opts = {
      method: 'POST',
      headers: {
        Location: newName
      }
    }
    return this.request(url, opts)
  }

  deleteApp (container, application) {
    const url = this.api + `apps/${container}/${application}`
    const opts = {
      method: 'DELETE'
    }
    return this.request(url, opts)
  }

  saveFile (file, value) {
    file.content = value
    const url = this.url + 'files' + file.url
    const opts = {
      method: 'POST',
      headers: {
        'Content-Type': file.mime
      },
      body: value
    }
    return this.request(url, opts)
  }

  createNewFile (name, template) {
    const url = this.url + 'files' + name
    const opts = {
      method: 'PUT',
      headers: {},
      body: ''
    }

    // Create file from template
    if (template) {
      opts.body = JSON.stringify(template)
      opts.headers['Content-Type'] = 'application/json; charset=utf-8'
    }

    opts.headers['Content-Length'] = opts.body.length
    return this.request(url, opts)
  }
}

export default ApiClient
