const mediaPattern = /\.(jpe?g|png|gif|aac|mp3|mp4|pdf)$/
const scriptPattern = /\.(jsx?|ts|coffee|es6)$/
const stylePattern = /\.(css|sss|scss|less)$/

class Application {
  constructor () {
    this.defaultFile = {content: ''}
    this.url = ''
    this.identifier = ''
    this.owner = ''
    this.pages = []
    this.files = []
    // current selected file
    this.current = this.defaultFile
  }

  get media () {
    let list = this.files.filter((f) => {
      return mediaPattern.test(f.name)
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
