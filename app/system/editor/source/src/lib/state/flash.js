/**
 *  Represents a flash message that can only be consumed once.
 *
 *  Useful when redirecting to a new view via the router and
 *  passing a single value. For example, used when redirecting
 *  to the page not found view.
 */
class Flash {
  constructor () {
    this._message = undefined
  }

  get message () {
    let f = this._message
    this._message = undefined
    return f
  }

  set message (msg) {
    this._message = msg
  }
}

export default Flash
