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
    // Slightly nicer for users to see than true/false
    if (typeof (val) === 'boolean') {
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
