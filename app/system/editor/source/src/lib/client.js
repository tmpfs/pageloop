import SocketConnection from './socket'
import {fetchFromRpc} from './services'

// Singleton websocket client
let socket

// Message identifier counter
let id = 0

class RpcRequest {
  constructor (id, method, params, ...args) {
    this.id = id
    this.method = method
    if (params && !Array.isArray(params)) {
      params = [params]
    }
    this.params = params || []
    this.args = args
  }

  get body () {
    return this.args && this.args.length && this.args[0]
  }

  get parameters () {
    return this.params[0]
  }
}

class Request {
  // Get a JSON RPC request object.
  static rpc (method, params, ...args) {
    const req = new RpcRequest(++id, method, params, ...args)
    if (!req.params.length) {
      delete req.params
    }
    if (!req.args.length) {
      delete req.args
    }
    return req
  }
}

class ApiClient {
  constructor () {
    // log should be injected
    this.log = null
    this.socket = this.connect()
    // this.useWebsocket = true
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
  request (url, opts) {
    const log = this.preflight(url, opts)
    return fetch(url, opts)
      .then((res) => {
        res.transport = 'http://rest-api'
        this.postflight(log, res)
        if (opts.raw) {
          return res
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
      // Clean REST specific flags
      delete req.mime
      delete req.raw
      delete req.fetch
      return this.socket.request(req)
        .then((res) => {
          res.url = url
          res.status = res.response.status
          this.postflight(log, res)
          return res
        })
    }

    // Send via standard REST API if the socket is not available
    return this.request(url, options)
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

  // Get a single application
  getApplication (container, application) {
    return this.rpc(Request.rpc('Application.Read', {context: container, target: application}))
  }

  // Get the files for an application
  getFiles (container, application) {
    return this.rpc(Request.rpc('Application.ReadFiles', {context: container, target: application}))
  }

  // Get the pages for an application
  getPages (container, application) {
    return this.rpc(Request.rpc('Application.ReadPages', {context: container, target: application}))
  }

  // Delete a list of files from an application
  deleteFiles (container, application, files) {
    const urls = files.map((f) => {
      return f.url
    })
    return this.rpc(
      Request.rpc('Application.DeleteFiles',
      {context: container, target: application},
      urls)
    )
  }

  // Run an application build task
  runTask (container, application, task) {
    // TODO: get container from app reference
    return this.rpc(
      Request.rpc('Application.RunTask',
      {context: container, target: application, item: task})
    )
  }

  // Create a new application.
  createApp (app) {
    // TODO: fix with template reference
    return this.rpc(
      Request.rpc('Container.CreateApp',
      {context: 'user'},
      app)
    )
  }

  // Delete an application.
  deleteApp (container, application) {
    return this.rpc(
      Request.rpc('Application.Delete',
      {context: container, target: application})
    )
  }

  // Create a new file optionally using the specified
  // template reference.
  createFile (container, application, url, template) {
    if (template) {
      return this.rpc(
        Request.rpc('File.CreateTemplate',
        {context: container, target: application, item: url}, template)
      )
    }
    // Create an empty file
    return this.rpc(
      Request.rpc('File.Create',
      {context: container, target: application, item: url}, '')
    )
  }

  // Save a file sending value as the new file content
  saveFile (container, application, file, value) {
    const req = Request.rpc('File.Save',
      {context: container, target: application, item: file.url}, value)
    req.raw = true
    req.mime = file.mime
    return this.rpc(req)
  }

  // Move a file
  moveFile (container, application, file, newName) {
    const req = Request.rpc('File.Move',
      {context: container, target: application, item: file.url}, newName)
    return this.rpc(req)
  }

  // TODO: get binary data over websocket!?
  getFileSource (container, application, file, raw) {
    const req = Request.rpc(raw ? 'File.ReadSourceRaw' : 'File.ReadSource',
      {context: container, target: application, item: file.url})
    // passthrough the underlying fetch promise
    req.fetch = true
    return this.rpc(req, {http: true})
  }
}

export default ApiClient
