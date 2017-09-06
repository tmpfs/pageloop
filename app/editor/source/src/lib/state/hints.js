class Hints {
  constructor () {
    this.messages = {
      'new-file': `
        Use <code>/path/to/file/document.md</code> to create directories when adding new files.
        To create a directory end the name with a slash (<code>/path/to/dir/</code>).
        `
    }

    // Test localStorage so hints are only displayed until dismissed
    // and not displayed again on subsequent visits
    for (let k in this.messages) {
      if (localStorage['hint:' + k]) {
        delete this.messages[k]
      }
    }
  }

  dismiss (id) {
    delete this.messages[id]
    localStorage['hint:' + id] = true
  }
}

export default Hints
