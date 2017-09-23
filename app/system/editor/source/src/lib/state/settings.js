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
      const nm = this.propName(key)
      this.keys[key] = nm

      // Set up properties
      Object.defineProperty(this, nm, {
        enumerable: true,
        configurable: false,
        get: () => {
          let val = this.get(key)
          if (val === undefined) {
            // Coerce defaults so booleans are displayed as 0 or 1
            return this.coerce('' + defaults[key])
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
        console.log('got local storage value: ' + localStorage[key])
        this.storage[key] = this.coerce(localStorage[key])
      }
    }
  }

  propName (key) {
    return camel(key.replace(/^[^:]+:/, ''))
  }

  coerce (val) {
    if (val === 'null') {
      return null
    } else if (val === 'true') {
      return 1
    } else if (val === 'false') {
      return 0
    } else if (parseInt(val).toString() === val) {
      return parseInt(val)
    } else if (!isNaN(Number(val))) {
      return Number(val)
    }
    return val
  }

  get (key) {
    return localStorage[key]
  }

  getDefault (key) {
    // Go via the property for coercion and default value
    return this[this.propName(key)]
  }

  del (key) {
    localStorage[key] = null
    delete this.storage[key]
  }

  set (key, value) {
    // Update the backing storage
    localStorage[key] = value

    // Sparse storage for reactive values by key
    this.storage[key] = value

    // All keys for reactive properties need to be mutated
    // so that data bindings fire
    this.keys[key] = value
  }

  reset () {
    let k, v
    for (k in this.keys) {
      v = this.coerce('' + defaults[k])
      this.set(k, v)
      this[this.propName(k)] = v
    }

    this.storage = {}

    // Clear all local storage items
    for (k in localStorage) {
      localStorage.removeItem(k)
    }
  }

  get length () {
    return localStorage.length
  }
}

export default Settings
