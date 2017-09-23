const defaults = {
  'setting:show-notifications': true,
  'setting:show-system-applications': false,
  'setting:show-template-applications': false,
  'hint:new-file': `
    Use <code>/path/to/file/document.md</code> to create directories when adding new files.
    To create a directory end the name with a slash (<code>/path/to/dir/</code>).
  `,
  'hint:drop-upload': `
    Drag and drop files here to upload.
  `,
  'hint:file-save': `
    Hit <kbd>Ctrl+s</kbd> (or <kbd>:w</kbd> in vim mode) to save and preview your changes.
  `
}

function camel (k) {
  return k.split('-').map((p, i) => {
    if (i) {
      return p.charAt(0).toUpperCase() + p.substr(1)
    }
    return p
  }).join('')
}

class Settings {
  constructor () {
    this.keys = {}
    for (let key in defaults) {
      // All keys for data bindings
      this.keys[key] = defaults[key]

      // Set up properties
      Object.defineProperty(this, this.propName(key), {
        enumerable: true,
        configurable: true,
        get: () => {
          return this.get(key)
        },
        set: (v) => {
          this.set(key, v)
        }
      })
    }
  }

  propName (key) {
    return camel(key.replace(/^[^:]+:/, ''))
  }

  coerce (val) {
    if (typeof val === 'string') {
      return JSON.parse(val)
    }
    return val
  }

  prettify (val) {
    const type = typeof (val)
    // Slightly nicer for users to see than true/false
    if (type === 'boolean') {
      return val ? 1 : 0
    }
    return val
  }

  get (key, pretty) {
    let val = localStorage[key]
    if (val === undefined) {
      return pretty ? this.prettify(defaults[key]) : defaults[key]
    }
    val = this.coerce(val)
    if (pretty) {
      val = this.prettify(val)
    }
    return val
  }

  del (key) {
    localStorage.removeItem(key)
  }

  // Returns the raw native value with no default lookup.
  getRaw (key) {
    return this.keys[key]
  }

  set (key, value) {
    // Update the backing storage
    localStorage[key] = JSON.stringify(value)
    // All keys for reactive properties need to be mutated
    // so that data bindings fire
    this.keys[key] = value
  }

  reset () {
    for (let k in this.keys) {
      this[this.propName(k)] = defaults[k]
      localStorage.removeItem(k)
    }
  }
}

export default Settings
