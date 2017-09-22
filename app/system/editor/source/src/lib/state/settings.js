const prefix = 'setting:'

const defaults = {
  'show-notifications': 1,
  'show-system-applications': 0,
  'show-template-applications': 0
}

class Settings {
  constructor () {
    const storage = localStorage
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
          // Local storage is string and we need booleans
          return storage[key] === '1'
        },
        set: (v) => {
          v = v ? '1' : '0'
          storage[key] = v
        }
      })
    }
  }
}

export default Settings
