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
}

export default Log
