const defaults = {
  'setting:show-notifications': true,
  'setting:show-system-applications': false,
  'setting:show-template-applications': false
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
    this.storage = {}
    this.keys = {}
    for (let key in defaults) {
      const nm = camel(key.replace(/^[^:]+:/, ''))
      this.keys[key] = nm

      // Set up properties
      Object.defineProperty(this, nm, {
        enumerable: true,
        configurable: false,
        get: () => {
          let val = this.get(key)
          if (val === undefined) {
            return defaults[key]
          }
          val = this.coerce(val)
          return val
        },
        set: (v) => {
          this.set(key, v)
        }
      })

      // Import from storage on load
      if (localStorage[key] !== undefined) {
        // this[nm] = localStorage[key]
        this.storage[key] = this.coerce(localStorage[key])
      }
    }
  }

  coerce (val) {
    if (val === 'null') {
      return null
    } else if (val === 'true') {
      return true
    } else if (val === 'false') {
      return false
    } else if (parseInt(val).toString() === val) {
      return parseInt(val)
    } else if (!isNaN(Number(val))) {
      return Number(val)
    }
    return val
  }

  get (key) {
    // Local storage is string and we need booleans
    return this.storage[key]
  }

  del (key) {
    localStorage[key] = null
    delete this.storage[key]
  }

  set (key, value) {
    localStorage[key] = value
    // Store for reactive values by key
    this.storage[key] = value
  }

  saved (key) {
    return localStorage[key] !== undefined
  }

  get length () {
    return localStorage.length
  }
}

export default Settings
