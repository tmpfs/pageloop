const prefix = 'setting:'

const defaults = {
  'show-notifications': true,
  'show-system-applications': false,
  'show-template-applications': false
}

class Settings {
  constructor (storage) {
    storage = storage || localStorage
    for (let k in defaults) {
      const key = prefix + k
      let nm = k.split('-').map((p, i) => {
        if (i) {
          return p.charAt(0).toUpperCase() + p.substr(1)
        }
        return p
      }).join('')

      // Set up defaults
      if (storage[key] === undefined) {
        storage[key] = defaults[k]
      }

      // Set up properties
      Object.defineProperty(this, nm, {
        enumerable: true,
        configurable: false,
        get: () => {
          return storage[key]
        },
        set: (v) => {
          storage[key] = v
        }
      })
    }
  }
}

export default Settings
