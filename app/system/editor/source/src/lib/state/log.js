class LogMessage {
  constructor () {
    this._level = ''
    this._message = ''
    this._error = false
    this._messages = []
  }

  get error () {
    return this._error
  }

  set error (err) {
    this._error = err
  }

  get level () {
    return this._level
  }

  get message () {
    return this._message
  }

  get messages () {
    return this._messages
  }

  parse (message) {
    if (typeof (message) === 'string') {
      this._level = 'Info'
      this._message = message
    } else if (message instanceof Error) {
      this._level = message.level || 'Warn'
      this._message = message.toString()
      this._error = message
    } else if (message && typeof (message) === 'object') {
      this._level = message.level || message.title
      this._message = message.message
    }
  }

  add (message) {
    let msg = message
    if (!(message instanceof LogMessage)) {
      msg = new LogMessage()
      msg.parse(message)
    }
    this.messages.unshift(msg)
  }
}

class Log {
  constructor () {
    this.maximum = 512
    this.messages = []
  }

  add (message) {
    let msg = message
    if (!(message instanceof LogMessage)) {
      msg = new LogMessage()
      msg.parse(message)
    }
    this.messages.unshift(msg)
    if (this.messages.length > this.maximum) {
      this.messages.pop()
    }
    return msg
  }
}

export default Log
