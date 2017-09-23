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

      // All keys for data bindings
      this.keys[key] = defaults[key]

      // Set up properties
      Object.defineProperty(this, nm, {
        enumerable: true,
        configurable: true,
        get: () => {
          return this.get(key)
        },
        set: (v) => {
          this.set(key, v)
        }
      })

      // Import from storage on load
      if (localStorage[key] !== undefined) {
        this.storage[key] = this.coerce(localStorage[key])
      }
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

  get (key) {
    let val = this.storage[key]
    if (val === undefined) {
      return defaults[key]
    }
    val = this.coerce(val)
    return val
  }

  /*
  del (key) {
    localStorage[key] = null
    delete this.storage[key]
  }
  */

  set (key, value) {
    // Update the backing storage
    localStorage[key] = JSON.stringify(value)

    // Sparse storage for reactive values by key
    this.storage[key] = value

    // All keys for reactive properties need to be mutated
    // so that data bindings fire
    this.keys[key] = value
  }

  reset () {
    console.log('reset to defaults')
    let k
    for (k in this.keys) {
      // v = this.coerce('' + defaults[k])
      console.log('setting key : ' + k)
      console.log('setting key value : ' + defaults[k])
      console.log('setting prop : ' + this.propName(k))

      this[this.propName(k)] = defaults[k]
      this.set(k, defaults[k])

      console.log('prop value: ' + this[this.propName(k)])
    }

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
