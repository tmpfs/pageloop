class Hints {
  constructor () {
    this.prefix = 'hint:'
    this.messages = {
      'new-file': `
        Use <code>/path/to/file/document.md</code> to create directories when adding new files.
        To create a directory end the name with a slash (<code>/path/to/dir/</code>).
        `,
      'drop-upload': `
        Drag and drop files here to upload.
        `
    }

    // Test localStorage so hints are only displayed until dismissed
    // and not displayed again on subsequent visits
    for (let k in this.messages) {
      if (localStorage[this.prefix + k]) {
        // delete this.messages[k]
      }
    }
  }

  dismiss (id) {
    this.messages[id] = null
    localStorage[this.prefix + id] = true
  }
}

export default Hints
