// TODO: expand and improve these extension matches
const textPattern = /\.(pdf|doc|docx|txt|md|textile|log)$/
const imagePattern = /\.(jpe?g|png|gif|tiff)$/
const scriptPattern = /\.(jsx?|ts|coffee|es6|sh)$/
const stylePattern = /\.(css|sss|scss|less)$/
const audioPattern = /\.(aac|mp3|wav|aiff?)$/
const videoPattern = /\.(mp4)$/

class Application {
  constructor () {
    this.defaultFile = {content: ''}
    this.url = ''
    this.identifier = ''
    this.owner = ''
    this._pages = []
    this._files = []
    // current selected file
    this.current = this.defaultFile

    this.build = {
      tasks: {}
    }
  }

  get files () {
    return this._files || []
  }

  set files (val) {
    this._files = val
  }

  get pages () {
    return this._pages || []
  }

  set pages (val) {
    this._pages = val
  }

  get media () {
    let list = this.files.filter((f) => {
      return textPattern.test(f.name) ||
        imagePattern.test(f.name) ||
        scriptPattern.test(f.name) ||
        stylePattern.test(f.name) ||
        audioPattern.test(f.name) ||
        videoPattern.test(f.name)
    })
    return list
  }

  get text () {
    let list = this.files.filter((f) => {
      return textPattern.test(f.name)
    })
    return list
  }

  get images () {
    let list = this.files.filter((f) => {
      return imagePattern.test(f.name)
    })
    return list
  }

  get scripts () {
    let list = this.files.filter((f) => {
      return scriptPattern.test(f.name)
    })
    return list
  }

  get styles () {
    let list = this.files.filter((f) => {
      return stylePattern.test(f.name)
    })
    return list
  }

  get audio () {
    let list = this.files.filter((f) => {
      return audioPattern.test(f.name)
    })
    return list
  }

  get video () {
    let list = this.files.filter((f) => {
      return videoPattern.test(f.name)
    })
    return list
  }

  isDirty () {
    for (let i = 0; i < this.files.length; i++) {
      if (this.files[i].dirty) {
        return true
      }
    }
    return false
  }

  getFileByUrl (url) {
    let i
    for (i = 0; i < this.files.length; i++) {
      if (this.files[i].url === url) {
        return this.files[i]
      }
    }
  }

  getPageByUrl (url) {
    let i
    for (i = 0; i < this.pages.length; i++) {
      if (this.pages[i].url === url) {
        return this.pages[i]
      }
    }
  }
}

export default Application
