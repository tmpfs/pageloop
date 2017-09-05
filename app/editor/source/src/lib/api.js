class ApiClient {
  constructor (container, application) {
    this.host = ''
    this.api = `/api/`
    this.container = container
    this.application = application
    this.url = `${this.api}${container}/${application}/`
    this.raw = `/apps/raw/${container}/${application}`
  }

  upload (file) {
    return new Promise((resolve, reject) => {
      let u = this.url

      if (file.dir) {
        u += file.dir
      }

      // This is used to simulate a 405 response
      // u += file.name
      u += 'files/' + file.upload.name

      // Need to use XHR for upload progress :(
      let xhr = new XMLHttpRequest()
      xhr.open('PUT', u, true)

      xhr.upload.onprogress = (e) => {
        if (e.lengthComputable) {
          let ratio = (e.loaded / e.total)
          file.info.ratio = ratio
        }
      }

      xhr.onload = function (e) {
        if (this.status !== 201) {
          const doc = JSON.parse(this.responseText)
          return reject(new Error(`Upload failed for ${file.name}: ${doc.error || doc.message}`))
        }
        file.complete = true

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

  json (url, options) {
    return fetch(url, options)
      .then((res) => res.json())
      .catch((err) => err)
  }

  createNewApp (app) {
    let url = this.api + 'user/'
    let opts = {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(app)
    }

    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  getContainers () {
    return this.json(this.api)
  }

  getApplication () {
    return this.json(this.url)
  }

  getPages () {
    let url = this.url + 'pages/'
    return this.json(url)
  }

  getFiles () {
    let url = this.url + 'files/'
    return this.json(url)
  }

  getFileContents (pathname) {
    let url = this.raw
    return fetch(url + pathname)
      .catch((err) => err)
  }

  deleteFile (file) {
    let url = this.url + 'files' + file.url
    let opts = {
      method: 'DELETE'
    }
    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  renameFile (file, newName) {
    let url = this.url + 'files' + file.url
    let opts = {
      method: 'POST',
      headers: {
        Location: newName
      }
    }
    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  deleteApp (container, application) {
    let url = `${this.api}${container}/${application}/`
    let opts = {
      method: 'DELETE'
    }
    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  saveFile (file, value) {
    file.content = value

    let url = this.url + 'files' + file.url
    let opts = {
      method: 'POST',
      headers: {
        'Content-Type': file.mime
      },
      body: value
    }
    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }

  createNewFile (name, template) {
    let url = this.url + 'files' + name
    let opts = {
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

    return fetch(url, opts)
      .then((res) => {
        return res.json().then((doc) => {
          return {response: res, document: doc}
        })
      })
  }
}

export default ApiClient
