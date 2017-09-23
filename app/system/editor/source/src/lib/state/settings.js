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
    let val = localStorage[key]
    if (val === undefined) {
      return defaults[key]
    }
    val = this.coerce(val)
    return val
  }

  del (key) {
    localStorage.removeItem(key)
  }

  set (key, value) {
    // Update the backing storage
    localStorage[key] = JSON.stringify(value)

    // All keys for reactive properties need to be mutated
    // so that data bindings fire
    this.keys[key] = value
  }

  reset () {
    let k
    for (k in this.keys) {
      this[this.propName(k)] = defaults[k]
      localStorage.removeItem(k)
    }
  }

  get length () {
    return localStorage.length
  }
}

export default Settings
