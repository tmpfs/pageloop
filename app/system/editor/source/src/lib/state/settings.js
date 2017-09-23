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
    return this.storage[key]
  }

  getDefault (key) {
    return this.get(key) || defaults[key]
  }

  del (key) {
    localStorage[key] = null
    delete this.storage[key]
  }

  set (key, value) {
    console.log(`set ${key} to ${value}`)

    localStorage[key] = value
    // Store for reactive values by key
    this.storage[key] = value

    // mutate keys for bindings
    this.keys[key] = value
  }

  reset () {
    console.log('settings reset')
    let k
    for (k in this.keys) {
      // Trigger properties for bindings to fire
      this[this.propName(k)] = defaults[k]
      console.log('settings prop name ' + this.propName(k))
      console.log('settings prop name ' + defaults[k])
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
