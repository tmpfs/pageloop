// Websocket endpoint
const WS = '/ws/'

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

class SocketConnection {
  constructor () {
    this.url = document.location.origin.replace(/^http/, 'ws') + WS
    this.protocols
    this.opts
    this._conn
    this._listeners = []

    this._keepalive = new SocketKeepalive(this, 30000)
  }

  get connected () {
    return this._conn && this._conn.readyState === WebSocket.OPEN
  }

  connect () {
    this._conn = new WebSocket(this.url, this.protocols, this.opts)

    this._conn.onopen = () => {
      // console.log('socket connection opened')
      this._keepalive.start()
    }

    this._conn.onmessage = (e) => {
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

    this._conn.onerror = (err) => {
      // TODO: log this error
      console.error(err)
    }

    this._conn.onclose = () => {
      console.log('socket connection closed')
      this.cleanup()
    }
  }

  cleanup () {
    this._keepalive.stop()
    this._conn.onopen = null
    this._conn.onmessage = null
    this._conn.onerror = null
    this._conn.onclose = null
    this._conn = null
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

          // TODO: reject on error???
          let doc = response.result || {}

          // Just vanilla rpc error, no status code available
          if (response.error) {
            doc.error = response.error
            res.status = 500
          } else if (response.result) {
            // Unwrap result object for status code
            doc = response.result.document
            res.status = response.result.status

            // TODO: handle wrapper error responses
          }

          resolve({response: res, document: doc})
        }
        this.send(payload)
      })
    }
  }
}

export default SocketConnection
