import keyNames from './keynames'

// Represents a mapping from a raw definition.
//
// Encapsulates the parsed keyName and any modifier keys
// specified in the mapping string,
class Mapping {
  constructor (key, fn) {
    this.key = key
    this.fn = fn
    this.keyName = undefined
    this.altKey = false
    this.ctrlKey = false
    this.metaKey = false
    this.shiftKey = false
  }

  // Are any of the modifier flags set.
  hasModifier () {
    return this.altKey || this.ctrlKey || this.metaKey || this.shiftKey
  }

  // Do modifier keys match, expects
  // a KeyEvent argument.
  modifiers (e) {
    if (!this.hasModifier()) {
      return true
    }
    return e.altKey === this.altKey &&
      e.ctrlKey === this.ctrlKey &&
      e.metaKey === this.metaKey &&
      e.shiftKey === this.shiftKey
  }

  // Does the key name match a given keyCode,
  // expects a KeyEvent argument.
  hasKeyName (e) {
    return keyNames[e.keyCode] === this.keyName
  }
}

// Mapping of key combinations to functions.
class KeyMap {
  constructor (keys) {
    this.keys = keys
    this.normalized = this.normalize(keys)
  }

  // Normalize creating Mapping instances for each key map
  // entry.
  normalize (keys) {
    const o = []
    const ptn = '+'
    let k, parts, mapping
    for (k in keys) {
      parts = k.split(ptn)

      mapping = new Mapping(k, keys[k])

      parts.forEach((p) => {
        const altKey = /^Alt/i.test(p)
        const ctrlKey = /^(Ctrl|Control)/i.test(p)
        const metaKey = /^(Meta|Mod)/i.test(p)
        const shiftKey = /^(Shift)/i.test(p)

        if (altKey || ctrlKey || metaKey || shiftKey) {
          if (altKey) mapping.altKey = altKey
          if (ctrlKey) mapping.ctrlKey = ctrlKey
          if (metaKey) mapping.metaKey = metaKey
          if (shiftKey) mapping.shiftKey = shiftKey
          return
        }

        if (mapping.keyName) {
          throw new Error(`Invalid key mapping ${k}, multiple key names`)
        }

        mapping.keyName = p
      })

      o.push(mapping)
    }
    return o
  }
}

// Encapsulates a collection of key maps and
// triggers the corresponding map function when
// a matching key stroke combination is detected.
//
// Key names are specified using a + delimiter
// between modifiers and the key name.
//
// Shift+Ctrl+N
//
// Alphabetic characters should be specified in
// uppercase.
//
// Only a single key name is permitted.
//
// The last key map added is the first to execute
// and key handler functions may defer to previously
// declared mappings by returning true.
class KeyManager {
  constructor () {
    this.maps = []

    // Test map
    /*
    this.add(new KeyMap({
      'Meta+N': () => console.log('Deferred called Meta+N'),
      'Shift+J': () => console.log('Shift+J')
    }))
    this.add(new KeyMap({
      'Meta+N': () => {
        return true
      }
    }))
    */

    /*
    this.add(new KeyMap({
      'Enter': () => {
        console.log('Enter called')
      },
      'Esc': () => {
        console.log('Esc called')
      }
    }))
    */

    // NOTE: must use keyup *not* keypress
    window.addEventListener('keyup', (e) => {
      e.preventDefault()
      e.stopImmediatePropagation()
      let i, fn, val
      const maps = this.find(e)
      for (i = 0; i < maps.length; i++) {
        fn = maps[i].fn
        if (typeof (fn) === 'function') {
          val = fn(e)
          if (val !== true) {
            break
          }
        }
      }
      return false
    })
  }

  find (e) {
    let i, j, map, list, mapping
    const maps = []
    for (i = 0; i < this.maps.length; i++) {
      map = this.maps[i]
      list = map.normalized
      for (j = 0; j < list.length; j++) {
        mapping = list[j]
        if (mapping.modifiers(e) && mapping.hasKeyName(e)) {
          maps.push(mapping)
          break
        }
      }
    }
    return maps
  }

  // Add a key map.
  add (map) {
    // Support plain object definitions
    if (map && !(map instanceof KeyMap)) {
      map = new KeyMap(map)
    }
    this.maps.unshift(map)
    return map
  }

  // Remove a key map.
  remove (map) {
    const ind = this.maps.indexOf(map)
    if (ind > -1) {
      return this.maps.splice(ind, 1)
    }
  }
}

export {KeyManager, KeyMap}
