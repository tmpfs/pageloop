// Websocket endpoint
const WS = '/ws/'

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

    // TODO: ping control functions and configurable interval
    this._pinger = setInterval(() => {
      // console.log('sending ping message')
      this.send({})
    }, 30000)

    this._conn.onopen = () => {
      console.log('socket connection opened')
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
        console.log(doc)
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
    this._conn.onopen = null
    this._conn.onmessage = null
    this._conn.onerror = null
    this._conn.onclose = null
    this._conn = null
  }

  // Send a JSON payload and ignore any response
  send (payload) {
    if (this.connected) {
      console.log('sending websocket request')
      console.log(payload)
      this._conn.send(JSON.stringify(payload))
    }
  }

  request (payload) {
    if (this.connected) {
      /*
      console.log('requesting with websocket connection')
      console.log(payload)
      */
      return new Promise((resolve, reject) => {
        // TODO: set timeout to remove listener
        this._listeners[payload.id] = (response) => {
          // console.log(response)
          const res = {
            status: response.status,
            id: response.id,
            transport: 'ws://json-rpc'}

          // TODO: reject on error???
          let doc = response.error || response.result

          // Unwrap result object for status code
          if (response.result) {
            doc = response.result.document
            res.status = response.result.status
          }

          // Work out how to get the error status
          if (response.error) {
            console.log('got rpc error response!!!:' + response.id)
            console.log('got rpc error response!!!:' + response.error)
            // res.status = doc.error.status || 500
          }

          console.log(doc)

          resolve({response: res, document: doc})
        }
        this.send(payload)
      })
    }
  }
}

export default SocketConnection
