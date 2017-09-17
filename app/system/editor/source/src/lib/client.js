class SocketConnection {
  constructor () {
    this.url = document.location.origin.replace(/^http/, 'ws') + '/ws/'
    this.protocols = undefined
    this.opts = undefined
    this._conn
  }

  connect () {
    this._conn = new WebSocket(this.url, this.protocols, this.opts)

    console.log(this._conn)

    this._conn.onopen = () => {
      console.log('socket conn opened')
      this._conn.send('Foo')
    }

    this._conn.onmessage = (e) => {
      console.log(e)
    }

    this._conn.onerror = (err) => {
      console.error(err)
    }

    this._conn.onclose = () => {
      console.log('socket conn closed')
    }

    /*
    function onOpen() {
      console.log('open called')
    }

    function onMessage(event) {
      console.log(event)
    }

    function onError(err) {
      console.error(err)
    }

    function onClose() {
      this.ws.onclose = null;
    }
    */
  }
}

class ApiClient {
  constructor (container, application) {
    this.host = ''
    this.api = `/api/`
    this.container = container
    this.application = application
    this.url = `${this.api}apps/${container}/${application}/`
    this.raw = `/apps/raw/${container}/${application}`
    // should be injected
    this.log = null

    this.websocket = new SocketConnection()
    this.websocket.connect()
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

  getVersion () {
    const url = this.api
    return this.request(url, {})
  }

  getStats () {
    const url = this.api + 'stats'
    return this.request(url, {})
  }

  getMeta () {
    return this.getVersion()
      .then(({response, document}) => {
        const meta = {info: document}
        return this.getStats()
          .then(({response, document}) => {
            meta.stats = document
            return {response: response, document: meta}
          })
      })
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

  getContainers () {
    const url = this.api + 'apps/'
    return this.request(url, {})
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

  listTemplates () {
    const url = this.api + 'templates'
    return this.request(url, {})
  }
}

export default ApiClient
