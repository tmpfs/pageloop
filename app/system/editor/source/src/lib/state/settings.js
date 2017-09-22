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
    for (let k in defaults) {
      const key = k
      const nm = camel(k.replace(/^[^:]+:/, ''))
      console.log(nm)

      // Set up defaults
      if (localStorage[key] === undefined) {
        this.set(key, defaults[k])
      } else {
        this.storage[key] = this.get(key)
      }

      // Set up properties
      Object.defineProperty(this, nm, {
        enumerable: true,
        configurable: false,
        get: () => {
          let val = this.get(key)
          console.log(val)
          if (typeof (defaults[key]) === 'boolean') {
            return val === '1'
          }
          return val === '1'
        },
        set: (v) => {
          if (typeof (v) === 'boolean') {
            v = v ? '1' : '0'
          }
          this.set(key, v)
        }
      })
    }
  }

  get (key) {
    // Local storage is string and we need booleans
    return localStorage[key]
  }

  del (key) {
    delete localStorage[key]
    delete this.storage[key]
  }

  set (key, value) {
    localStorage[key] = value
    // Store for reactive values by key
    this.storage[key] = value
  }

  get length () {
    return localStorage.length
  }

  get keys () {
    const defs = {}
    for (let k in defaults) {
      defs[k] = camel(k.replace(/^[^:]+:/, ''))
    }
    return defs
  }
}

export default Settings
