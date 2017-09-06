// Mapping of key combinations to functions.
//
// If element is given the element must be focused
// for mappings in this key map to be invoked.
class KeyMap {
  constructor (map, element) {
    this.map = map
    this.element = element
  }
}

// Encapsulates a collection of key maps and
// triggers the corresponding map function when
// a matching key stroke combination is detected.
class KeyManager {

  constructor () {
    this.maps = []

    this.add(new KeyMap({
      'Meta+n': () => console.log('Meta + n'),
      'Shift+j': () => console.log('Shift + j')
    }))

    window.addEventListener('keyup', (e) => {
      e.preventDefault()
      e.stopImmediatePropagation()
      console.log(e)
      const fn = this.find(e)
      if (typeof (fn) === 'function') {
        fn(e)
      }
      return false
    })
  }

  find (e) {
    const key = e.key
    const code = e.code

    console.log(key)
    console.log(code)

    let i, map, k
    for (i = 0; i < this.maps.length; i++) {
      map = this.maps[i]
      for (k in map) {
        console.log(k)
      }
    }
  }

  // Add a key map
  add (map) {
    // TODO: normalize map to use keycodes from string keys
    this.maps.unshift(map)
  }

  // Remove a key map
  remove (map) {
    const ind = this.maps.indexOf(map)
    if (ind > -1) {
      this.maps.splice(ind, 1)
    }
  }
}

export {KeyManager, KeyMap}
