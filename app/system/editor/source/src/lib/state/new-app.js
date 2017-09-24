class NewApp {
  constructor () {
    this.templateUrl = ''
    this.template = undefined
    this.description = ''

    this._id = ''
    this._name = ''
  }

  get id () {
    return this._id
  }

  get name () {
    return this._name
  }

  set name (val) {
    this._name = val
  }

  get valid () {
    return this.name !== '' && this.description !== ''
  }
}

export default NewApp
