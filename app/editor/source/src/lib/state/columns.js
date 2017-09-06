class ColumnManager {
  constructor () {
    this.state = null
    this.styles = null
    this.maximized = ''
    // Keep track of whether a custom column has been set
    // via a drag operation so we can determine if a reset is needed
    this.custom = false

    this.doDrag = (e) => {
      e.stopImmediatePropagation()
      let target = this.state.target
      let parent = this.state.parent
      let index = this.state.index
      let maximum = this.state.maximum
      let tb = target.getBoundingClientRect()

      // Out of bounds cursor
      if (e.clientX < 0 || e.clientX > maximum || e.clientX < tb.left) {
        return
      }

      // Try to stop other events interfering
      // TODO: work out why cursor is not displaying
      document.querySelector('body').setAttribute('style', 'user-select: none; pointer-events: none; cursor: ew-resize;')

      // Resize target column by percentage
      let percent = Math.round(((e.clientX - tb.left) / maximum) * 100)
      target.setAttribute('style', 'max-width: none; width:' + percent + '%')

      // Columns after the target being resized, they are
      // compacted but maintain the aspect ratio
      let taken = this.state.widths + percent
      let remainder = 100 - taken
      let i, n
      for (i = index + 1; i < parent.childNodes.length; i++) {
        n = parent.childNodes[i]
        if (n.nodeType === 1) {
          n.setAttribute('style', 'max-width: none; width: ' + (remainder * this.state.ratios[i]) + '%')
        }
      }

      this.custom = true
    }

    this.stopDrag = (e) => {
      e.stopImmediatePropagation()
      document.querySelector('body').removeAttribute('style')
      window.removeEventListener('mousemove', this.doDrag)
      this.state = null
    }
  }

  startDrag (e) {
    e.stopImmediatePropagation()

    // Target to reize
    let target = e.currentTarget.parentNode

    // Parent gives overall available width
    let parent = target.parentNode

    // Width of all columns to calculate percentage
    let pb = parent.getBoundingClientRect()
    let maximum = pb.right - pb.left

    this.state = {
      target: target,
      parent: parent,
      maximum: maximum,
      index: undefined,
      widths: 0,
      ratios: []
    }

    // Used to track remaining available pixels
    let total = 0

    let i, n, b, w, ratio, percent
    for (i = 0; i < parent.childNodes.length; i++) {
      n = parent.childNodes[i]
      if (n.nodeType !== 1) {
        continue
      }

      b = n.getBoundingClientRect()
      w = b.right - b.left

      // Get ratios of subsequent columns
      if (this.state.index !== undefined) {
        // How much of the remaining space is used by this column
        ratio = w / (maximum - total)
        // Sparse array!
        this.state.ratios[i] = ratio
      }

      ratio = w / maximum
      percent = Math.round(ratio * 100)

      if (n === target) {
        this.state.index = i
        total += w
      }

      // Fix widths of previous columns
      if (this.state.index === undefined) {
        total += w
        this.state.widths += percent
        n.setAttribute('style', 'max-width: none; width:' + percent + '%')
      }
    }

    // Start the drag operation
    window.addEventListener('mousemove', this.doDrag)

    // Need to capture on the window for mouse up outside
    window.addEventListener('mouseup', this.stopDrag)
  }

  maximize (className) {
    const el = document.querySelector('.' + className)
    const parent = el.parentNode
    this.styles = {}
    parent.childNodes.forEach((n, index) => {
      if (n.nodeType !== 1) {
        return
      }
      this.styles[index] = n.getAttribute('style')
      if (n === el) {
        n.setAttribute('style', 'max-width: none; width: 100%;')
      } else {
        n.setAttribute('style', 'max-width: none; width: 0%;')
      }
    })
    this.maximized = className
  }

  minimize (className) {
    if (!className) {
      return
    }
    const el = document.querySelector('.' + className)
    const parent = el.parentNode
    parent.childNodes.forEach((n, index) => {
      if (n.nodeType !== 1) {
        return
      }
      if (this.styles[index]) {
        n.setAttribute('style', this.styles[index])
      } else {
        n.removeAttribute('style')
      }
    })
    this.styles = null
    this.maximized = ''
  }

  // Remove inline styles from columns will restore columns
  // to the defaults declared in the stylesheet.
  reset () {
    let parent = document.querySelector('.content-main > .content')
    let i, n
    this.styles = {}
    for (i = 0; i < parent.childNodes.length; i++) {
      n = parent.childNodes[i]
      if (n.nodeType !== 1) {
        continue
      }
      n.removeAttribute('style')
    }
    this.custom = false
  }
}

export default ColumnManager
