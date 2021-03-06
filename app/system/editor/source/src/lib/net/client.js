import SocketConnection from './socket'
import {Request, TypeByte} from './request'
import {fetchFromRpc} from './services'

// Singleton websocket client
let socket

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

    const startTime = Date.now()
    // console.log(opts)
    return fetch(url, opts)
      .then((res) => {
        res.method = opts.method
        res.transport = 'application/rest+api'
        const time = {
          start: startTime,
          end: Date.now()
        }
        time.duration = time.end - time.start
        res.time = time
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
    // Build url and options here so we can log the URL
    const {url, options} = fetchFromRpc(req)

    if (/\/\//.test(url) || /\.\.\/?/.test(url)) {
      throw new Error(`Bad request URL: ${url}`)
    }

    // Try to send over socket connection first
    if (this.useWebsocket && this.socket.connected && !opts.http) {
      const log = this.preflight(url, {method: 'RPC'})
      const startTime = Date.now()
      return this.socket.request(req)
        .then((res) => {
          // console.log(res)
          res.url = url
          res.response.duration = Date.now() - startTime

          const time = {
            start: startTime,
            end: Date.now()
          }
          time.duration = time.end - time.start
          res.response.time = time
          res.response.url = socket.url
          res.response.method = 'RPC'
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

  // List services
  listServices () {
    return this.rpc(Request.rpc('Service.List'))
  }

  // List all containers
  getContainers () {
    return this.rpc(Request.rpc('Host.List'))
  }

  getApplicationReference (container, application) {
    /*
    return {
      name: application,
      container: container
    }
    */
    return {
      // For websocket app refs
      ref: `file://${container}/${application}`,

      // For rest parameter path building
      container: container,
      name: application
    }
  }

  getAppRef (container, application) {
    return `file://pageloop.com/${container}/${application}`
  }

  getFileRef (container, application, url) {
    return `file://pageloop.com/${container}/${application}#${url}`
  }

  // Get a single application
  getApplication (container, application) {
    const params = {
      ref: this.getAppRef(container, application)
    }
    const req = Request.rpc('Application.Read', params)
    return this.rpc(req)
  }

  // Delete an application.
  deleteApp (container, application) {
    /*
    const ref = this.getApplicationReference(container, application)
    console.log(ref)
    return this.rpc(Request.rpc('Application.Delete', ref))
    */

    const params = {
      ref: this.getAppRef(container, application)
    }
    const req = Request.rpc('Application.Delete', params)
    return this.rpc(req)
  }

  // Get the files for an application
  getFiles (container, application) {
    const params = {
      ref: this.getAppRef(container, application)
    }
    const req = Request.rpc('Application.ReadFiles', params)
    return this.rpc(req)
  }

  // Get the pages for an application
  getPages (container, application) {
    const params = {
      ref: this.getAppRef(container, application)
    }
    const req = Request.rpc('Application.ReadPages', params)
    return this.rpc(req)
  }

  // Delete a list of files from an application
  deleteFiles (container, application, files) {
    const urls = files.map((f) => {
      return f.url
    })
    const params = {
      ref: this.getAppRef(container, application),
      batch: urls
    }
    const req = Request.rpc('Application.DeleteFiles', params)
    req.json(urls)
    return this.rpc(req)
  }

  // Run an application build task
  runTask (container, application, task) {
    const params = {
      ref: this.getAppRef(container, application),
      task: task
    }
    const req = Request.rpc('Application.RunTask', params)
    return this.rpc(req)
  }

  // Create a new application.
  createApp (app) {
    // TODO: fix with template reference
    app.container = 'user'
    const req = Request.rpc('Container.CreateApp', app)
    req.json(app)
    return this.rpc(req)
  }

  // Create a new file optionally using the specified
  // template reference.
  createFile (container, application, url, template) {
    const params = {
      ref: this.getFileRef(container, application, url)
    }
    if (template) {
      params.template = template
      return this.rpc(Request.rpc('File.CreateTemplate', params))
    }
    // Create an empty file
    return this.rpc(Request.rpc('File.Create', params))
  }

  // Save a file sending value as the new file content
  saveFile (container, application, file, value) {
    const params = {
      ref: this.getFileRef(container, application, file.url)
    }

    // Text files can be sent as strings (websocket only)
    if (!file.binary) {
      params.value = value
    } else {
      // TODO: send as binary when file is binary
    }

    const req = Request.rpc('File.Save', params)

    // Send as raw request (REST requests only)
    req.body(value, file.mime)

    return this.rpc(req)
  }

  // Move a file
  moveFile (container, application, file, newName) {
    const params = {
      ref: this.getFileRef(container, application, file.url)
    }
    params.destination = newName
    const req = Request.rpc('File.Move', params)
    return this.rpc(req)
  }

  // Get file contet.
  //
  // When the raw option is given the response document will include
  // frontmatter data when available.
  getFileSource (container, application, file, raw) {
    const params = {
      ref: this.getFileRef(container, application, file.url)
    }
    const req = Request.rpc(raw ? 'File.ReadSourceRaw' : 'File.ReadSource', params)

    // TODO: allow this over websocket so we don't need to force a transport
    return this.rpc(req, {http: true})
  }

  // Export public zip
  exportArchive (container, application, filter) {
    const params = {
      ref: this.getAppRef(container, application)
    }
    if (filter) {
      params.filter = filter
    }
    /*
    const ref = this.getApplicationReference(container, application)
    if (filter) {
      ref.filter = filter
    }
    */
    const req = Request.rpc('Archive.Export', params)
    const fetch = fetchFromRpc(req)
    console.log('download from: ' + fetch.url)
    // TODO: try not to use replace()
    document.location.replace(fetch.url)
    // return this.rpc(req, {http: true})
  }
}

export {ApiClient, Request}
