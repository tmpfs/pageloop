import SocketConnection from './socket'
import {fetchFromRpc} from './services'

// Singleton websocket client
let socket

// Message identifier counter
let id = 0

const TypeJson = 0
const TypeByte = 1

class Request {
  constructor (id, method, params) {
    this.id = id
    this.method = method
    if (params && !Array.isArray(params)) {
      params = [params]
    }
    this.params = params || []

    Object.defineProperty(this, 'request', {
      enumerable: false,
      configurable: false,
      writable: true
    })
  }

  // Set an object to be serialized to JSON and
  // sent as the request body (REST only).
  json (value) {
    this.request = {
      type: TypeJson,
      value: value
    }
  }

  // Set a raw body for the request.
  body (value, mime) {
    this.request = {
      type: TypeByte,
      value: value,
      mime: mime
    }
  }

  get parameters () {
    return this.params[0]
  }

  // Get a JSON RPC request object.
  static rpc (method, params) {
    return new Request(++id, method, params)
  }
}

class ApiClient {
  constructor () {
    // log should be injected
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

  // Perform an API request and assume a JSON response unless
  // the raw option is set.
  request (url, opts, id, method) {
    const log = this.preflight(url, opts)

    console.log(`[rest] (${id}) ${method} ${opts.method} ${url}`)
    console.log(JSON.stringify(opts))

    // console.log(opts)
    return fetch(url, opts)
      .then((res) => {
        res.transport = 'http://rest-api'
        this.postflight(log, res)
        const resType = parseInt(res.headers.get('x-response-type'))
        if (!isNaN(resType)) {
          // Server says this is a binary response, do not
          // interpret it as JSON
          if (resType === TypeByte) {
            return res
          }
        }
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  rpc (req, opts = {}) {
    const {url, options} = fetchFromRpc(req)
    // Try to send over socket connection first
    if (this.useWebsocket && this.socket.connected && !opts.http) {
      const log = this.preflight(url, {method: 'RPC'})
      return this.socket.request(req)
        .then((res) => {
          // console.log(res)
          res.url = url
          res.status = res.response.status
          this.postflight(log, res)
          return res
        })
    }

    // Send via standard REST API if the socket is not available
    return this.request(url, options, req.id, req.method)
  }

  upload (container, application, file) {
    return new Promise((resolve, reject) => {
      let u = `/api/apps/${container}/${application}/files/`
      let dir = file.dir
      if (dir) {
        u += dir
      }

      if (!/\/$/.test(u)) {
        u += '/'
      }

      u += file.name

      // TODO: log file uploads
      //
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
      .then((res) => {
        if (res.response.status !== 200) {
          return res
        }
        const meta = {info: res.document}
        return this.getStats()
          .then((res) => {
            meta.stats = res.document
            return {response: res.response, document: meta}
          })
      })
  }

  // List application templates
  listTemplates () {
    return this.rpc(Request.rpc('Template.List'))
  }

  // List all containers
  getContainers () {
    return this.rpc(Request.rpc('Host.List'))
  }

  getApplicationReference (container, application) {
    return {
      name: application,
      container: container
    }
  }

  getFileReference (container, application, url) {
    return {
      // For websocket file refs
      ref: `file://${container}/${application}#${url}`,

      // For rest parameter path building
      container: container,
      application: application,
      url: url
    }
  }

  // Get a single application
  getApplication (container, application) {
    const ref = this.getApplicationReference(container, application)
    return this.rpc(Request.rpc('Application.Read', ref))
  }

  // Get the files for an application
  getFiles (container, application) {
    const ref = this.getApplicationReference(container, application)
    return this.rpc(Request.rpc('Application.ReadFiles', ref))
  }

  // Get the pages for an application
  getPages (container, application) {
    const ref = this.getApplicationReference(container, application)
    return this.rpc(Request.rpc('Application.ReadPages', ref))
  }

  // Delete a list of files from an application
  deleteFiles (container, application, files) {
    const urls = files.map((f) => {
      return f.url
    })
    const ref = this.getApplicationReference(container, application)
    ref.batch = urls
    const req = Request.rpc('Application.DeleteFiles', ref)
    req.json(urls)
    return this.rpc(req)
  }

  // Run an application build task
  runTask (container, name, task) {
    return this.rpc(Request.rpc('Application.RunTask', {container, name, task}))
  }

  // Create a new application.
  createApp (app) {
    // TODO: fix with template reference
    app.container = 'user'
    const req = Request.rpc('Container.CreateApp', app)
    req.json(app)
    return this.rpc(req)
  }

  // Delete an application.
  deleteApp (container, application) {
    const ref = this.getApplicationReference(container, application)
    return this.rpc(Request.rpc('Application.Delete', ref))
  }

  // Create a new file optionally using the specified
  // template reference.
  createFile (container, application, url, template) {
    const ref = this.getFileReference(container, application, url)
    if (template) {
      ref.template = template
      return this.rpc(Request.rpc('File.CreateTemplate', ref))
    }
    // Create an empty file
    return this.rpc(Request.rpc('File.Create', ref))
  }

  // Save a file sending value as the new file content
  saveFile (container, application, file, value) {
    const ref = this.getFileReference(container, application, file.url)

    // Text files can be sent as strings (websocket only)
    if (!file.binary) {
      ref.value = value
    } else {
      // TODO: send as binary when file is binary
    }

    const req = Request.rpc('File.Save', ref)

    // Send as raw request (REST requests only)
    req.body(value, file.mime)

    return this.rpc(req)
  }

  // Move a file
  moveFile (container, application, file, newName) {
    const ref = this.getFileReference(container, application, file.url)
    ref.destination = newName
    const req = Request.rpc('File.Move', ref)
    return this.rpc(req)
  }

  // Get file contet.
  //
  // When the raw option is given the response document will include
  // frontmatter data when available.
  getFileSource (container, application, file, raw) {
    const ref = this.getFileReference(container, application, file.url)
    const req = Request.rpc(raw ? 'File.ReadSourceRaw' : 'File.ReadSource', ref)

    // TODO: allow this over websocket so we don't need to force a transport
    return this.rpc(req, {http: true})
  }
}

export default ApiClient
