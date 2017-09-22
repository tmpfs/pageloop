// Websocket endpoint
const WS = '/ws/'

/**
 *  Sends ping messages over the socket at regular intervals
 *  to keep the connection alive.
 *
 *  A ping message is the empty object: {}.
 */
class SocketKeepalive {
  constructor (conn, duration) {
    this.conn = conn
    this.duration = duration
  }

  stop () {
    clearInterval(this._interval)
  }

  ping () {
    this.conn.send({})
  }

  start () {
    this.stop()
    this._interval = setInterval(() => {
      this.ping()
    }, this.duration)
  }
}

/**
 *  Attempts to open a closed socket connection.
 */
class SocketOpener {
  constructor (conn, duration, limit) {
    this.conn = conn
    this.duration = duration
    this.limit = limit
    this._connecting = false
  }

  retry () {
    if (this._connecting) {
      return false
    }
    this._connecting = true
    this.connect(0)
    return true
  }

  /**
   *  @private
   */
  connect (retries) {
    const attempts = retries + 1
    const duration = this.duration
    const exp = duration * attempts

    if (this.limit !== undefined && retries === this.limit) {
      console.log(`[sock] retry limit ${this.limit} reached`)
      return false
    }

    console.log(`[sock] retry attempt ${attempts}, timeout in ${exp / 1000}s`)

    const timeout = setTimeout(() => {
      this.connect(retries + 1)
    }, exp)

    this.conn.connect(() => {
      // Connection re-established
      this._connecting = false
      clearTimeout(timeout)
      console.log(`[sock] connection re-established`)
    })
  }
}

/**
 *  Represents a socket connection.
 *
 *  Note that retryTimeout is multiplied on each attempt so it is ok
 *  to use a low value like 2000.
 *
 *  @param {Object} options
 *
 *  @option {Object} protocols protocols for Websocket.
 *  @option {Object} websocket options for Websocket.
 *  @option {Boolean=true} keepalive send keepalive ping requests.
 *  @option {Number=30000} keepaliveDuration interval for keepalive ping messages.
 *  @option {Boolean=true} retry attempt to re-connect on close event.
 *  @option {Number=2000} retryTimeout interval for re-connection attempts.
 *  @option {Number=64} retryLimit number of re-connection attempts before giving up.
 */
class SocketConnection {
  constructor (options = {}) {
    this.url = document.location.origin.replace(/^http/, 'ws') + WS
    this.protocols = options.protocols
    this.opts = options.websocket

    this._conn
    this._listeners = []
    if (options.keepalive === undefined || options.keepalive) {
      this._keepalive = new SocketKeepalive(this, options.keepaliveDuration || 30000)
    }
    if (options.retry === undefined || options.retry) {
      this._opener = new SocketOpener(this, options.retryTimeout || 2000, options.retryLimit || 64)
    }
  }

  get socket () {
    return this._conn
  }

  get connected () {
    return this._conn && this._conn.readyState === WebSocket.OPEN
  }

  connect (cb) {
    this._conn = new WebSocket(this.url, this.protocols, this.opts)

    this._conn.onopen = (e) => {
      if (typeof cb === 'function') {
        cb(e)
      }
      this.onOpen(e)
    }

    this._conn.onmessage = (e) => {
      this.onMessage(e)
    }

    this._conn.onerror = (err) => {
      this.onError(err)
    }

    this._conn.onclose = (e) => {
      this.onClose(e)
    }
  }

  onOpen (e) {
    // console.log('socket connection opened')
    this._keepalive.start()
  }

  onMessage (e) {
    // console.log(e)
    if (e.data) {
      let doc
      try {
        doc = JSON.parse(e.data)
      } catch (e) {
        throw e
      }
      // console.log(doc)
      if (doc.id && this._listeners[doc.id]) {
        this._listeners[doc.id](doc)
        delete this._listeners[doc.id]
      }
    }
  }

  onError (err) {
    // TODO: log this error
    console.error(err)
  }

  onClose (e) {
    console.log('socket connection closed')
    this.cleanup()
    this._opener.retry()
  }

  cleanup () {
    this._keepalive.stop()
    this._conn.onopen = null
    this._conn.onmessage = null
    this._conn.onerror = null
    this._conn.onclose = null
    // this._conn = null
  }

  // Send a JSON payload and ignore any response
  send (payload) {
    if (this.connected) {
      const encoded = JSON.stringify(payload)
      console.log(`[sock] (${payload.id}) ${payload.method}`)
      console.log(encoded)
      this._conn.send(encoded)
    }
  }

  request (payload) {
    if (this.connected) {
      return new Promise((resolve, reject) => {
        // TODO: set timeout to remove listener
        this._listeners[payload.id] = (response) => {
          const res = {
            status: response.status,
            id: response.id,
            transport: 'ws://json-rpc'}

          let doc = response.result || {}

          // Just vanilla rpc error, no status code available
          if (response.error) {
            doc.error = response.error
            res.status = 500
          } else if (response.result) {
            // Server replied with result containing custom error data
            if (doc.error) {
              // Unwrap error from rpc result object
              // allows passing custom status codes
              res.status = doc.error.status
              doc.error = doc.error.message
            // Server replied with valid result
            } else {
              // Unwrap result object for status code
              doc = response.result.document
              res.status = response.result.status
            }
          }

          resolve({response: res, document: doc})
        }
        this.send(payload)
      })
    }
  }
}

export default SocketConnection
