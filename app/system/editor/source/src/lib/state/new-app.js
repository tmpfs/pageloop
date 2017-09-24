const defaultView = 'new-app-info'

class NewApp {
  constructor () {
    this.reset()
  }

  get id () {
    return this._id
  }

  get name () {
    return this._name
  }

  set name (val) {
    this._name = val
    // Remove invalid characters
    let id = val.replace(/[^-a-zA-Z0-9 ]/g, '')
    // Normalize whitespace to hyphens
    id = id.replace(/\s+/g, '-')
    // May not begin with a hyphen
    id = id.replace(/^-/, '')
    // Should be lowercase
    id = id.toLowerCase()
    this._id = id
  }

  get valid () {
    return this.name !== '' && this.description !== ''
  }

  reset () {
    this.view = defaultView
    this.templateUrl = ''
    this.template = undefined
    this.description = ''

    this._id = ''
    this._name = ''
  }
}

export default NewApp
