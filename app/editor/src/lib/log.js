class Log {
  constructor () {
    this.maximum = 1024
    this.messages = []
  }

  add (message) {
    this.messages.push(message)
    if (this.messages.length > this.maximum) {
      this.messages.shift()
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
