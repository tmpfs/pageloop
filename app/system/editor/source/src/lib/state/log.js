class Log {
  constructor () {
    this.maximum = 512
    this.messages = []
  }

  add (message) {
    if (typeof (message) === 'string') {
      message = {
        level: 'Info',
        message: message
      }
    } else if (message instanceof Error) {
      console.log('treat as log error')
      message = {
        level: message.level || 'Warn',
        message: message.toString(),
        error: message
      }
    }
    this.messages.unshift(message)
    if (this.messages.length > this.maximum) {
      this.messages.pop()
    }
  }

  get last () {
    let m = null
    if (this.messages.length) {
      m = this.messages[this.messages.length - 1]
    }
    return m
  }

  toString () {
    let m = this.last
    if (m instanceof Error) {
      return m.message || ('' + m)
    }
    return m
  }
}

export default Log
