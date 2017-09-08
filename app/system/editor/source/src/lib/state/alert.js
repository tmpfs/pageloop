/**
 *  Encapsulates the state for a UI alert.
 */
class Alert {
  constructor () {
    this.visible = false
    this.title = 'Alert'
    this.message = ''
    this.note = ''
    this.ok = function noop () {}
  }
}

export default Alert
