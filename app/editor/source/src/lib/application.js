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