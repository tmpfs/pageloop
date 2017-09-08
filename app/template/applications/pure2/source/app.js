'use strict'

;(function () {
  function easeOutQuad (iteration, start, diff, total) {
    return -diff * (iteration /= total) * (iteration - 2) + start
  }

  class Scroll {
    constructor (opts = {}) {
      let id = this.onScrollToLink.bind(this)
      let elements = opts.id || []
      elements.forEach((el) => {
        el.addEventListener('click', id)
      })

      this.duration = opts.duration || 50
    }

    getScrollHeight () {
      return Math.max(
        document.body.scrollHeight,
        document.body.offsetHeight,
        document.documentElement.clientHeight,
        document.documentElement.scrollHeight,
        document.documentElement.offsetHeight)
    }

    getScrollLimit () {
      return this.getScrollHeight() - window.innerHeight
    }

    onScrollTop (e) {
      e.preventDefault()
      this.scrollToTop(0)
    }

    getScrollPosition () {
      let doc = document.documentElement
      let left = (window.pageXOffset || doc.scrollLeft) - (doc.clientLeft || 0)
      let top = (window.pageYOffset || doc.scrollTop) - (doc.clientTop || 0)
      return {left: left, top: top}
    }

    scrollToTop (val) {
      let start = this.getScrollPosition().top
      let iteration = 0
      let duration = this.duration
      let diff = val === 0 ? -start : val
      let requestAnimationFrame = window.requestAnimationFrame ||
                                  window.mozRequestAnimationFrame ||
                                  window.webkitRequestAnimationFrame ||
                                  window.msRequestAnimationFrame

      // perform the animation
      function doScroll () {
        const value = easeOutQuad(iteration, start, diff, duration)
        const amount = value < 0 ? -value : value
        window.scrollTo(0, amount)
        if (iteration >= duration) {
          window.scrollTo(0, Math.floor(amount))
          return
        }
        requestAnimationFrame(doScroll)
        iteration++
      }

      doScroll()
    }

    onScrollToLink (e) {
      e.preventDefault()
      const id = e.currentTarget.getAttribute('href').replace(/^#/, '')
      this.scrollToId(id)
    }

    scrollToId (id) {
      const el = document.getElementById(id)
      if (!el) {
        return false
      }
      const bounds = el.getBoundingClientRect()
      this.scrollToTop(bounds.top)
    }

  }

  return new Scroll({id: document.querySelectorAll('[href^="#"]')})
})()
