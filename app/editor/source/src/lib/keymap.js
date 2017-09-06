const ALT = {test: 'Alt'}
const SHIFT = {test: 'Shift'}
const META = {test: 'Meta'}
const CONTROL = {test: 'Control', value: 'Ctrl'}

const normalizations = [
  ALT,
  SHIFT,
  META,
  CONTROL
]

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
    this.pressed = {}

    this.add(new KeyMap({
      'Meta+n': () => console.log('Meta + n'),
      'Shift+j': () => console.log('Shift + j')
    }))

    document.addEventListener('keydown', (e) => {
      console.log('GOT KEY DOWN: ' + this.normalize(e))
      console.log(e)
      this.pressed[this.normalize(e)] = {key: e.key, code: e.code, keyCode: e.keyCode}
    })

    document.addEventListener('keyup', (e) => {
      console.log('GOT KEY UP: ' + this.normalize(e))
      this.pressed[this.normalize(e)] = {key: e.key, code: e.code, keyCode: e.keyCode}
      delete this.pressed[this.normalize(e)]
    })

    document.addEventListener('keypress', (e) => {
      e.preventDefault()
      e.stopImmediatePropagation()
      console.log('keypress: ' + this.normalize(e))
      delete this.pressed[this.normalize(e)]
      console.log(e)
      const fn = this.find(e)
      if (typeof (fn) === 'function') {
        fn()
      }
      return false
    })
  }

  normalize (e) {
    let i, n
    const key = e.key
    for (i = 0; i < normalizations.length; i++) {
      n = normalizations[i]
      if (key.indexOf(n.test) === 0) {
        return n.value || n.test
      }
    }

    return key
  }

  find (e) {
    const key = e.key
    const code = e.code

    console.log(key)
    console.log(code)

    console.log(this.pressed)

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
