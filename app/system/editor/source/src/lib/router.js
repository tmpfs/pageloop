/**
 *  Lightweight router implementation using regular expressions.
 *
 *  The Vue router does not appear to support matching the slash
 *  character and we need to match URL references (eg: /docs/help/document.html).
 */
class Router {
  constructor (href, strip) {
    this.defaultHref = href
    this.routes = []
    this.strip = strip
  }

  navigate (href, state) {
    let url = this.url(href)
    // TODO: work out error with setting state!
    // history.pushState({href: href, state: state}, '', url)
    history.pushState({href: href, state: null}, '', url)
    this.route(href)
  }

  url (href) {
    return this.pathname + '#' + href
  }

  get pathname () {
    return document.location.pathname
  }

  get hash () {
    let h = document.location.hash.replace(/^#/, '')
    if (this.strip) {
      h = h.replace(/\/$/, '')
    }
    return h
  }

  replace (href, trigger) {
    document.location.replace(this.url(href))
    if (trigger) {
      this.route(href)
    }
  }

  add (ptn, map, fn) {
    if (typeof map === 'function') {
      fn = map
      map = null
    }
    this.routes.push({ptn: ptn, fn: fn, map: map})
  }

  route (href, state) {
    function result (href, route) {
      let o = {
        state: state,
        href: href,
        route: route,
        parts: [],
        map: {}
      }

      let parts = href.replace(/^\//, '').replace(/\/$/, '').split('/')
      o.parts = parts

      if (route.map) {
        route.map.forEach((val, i) => {
          o.map[val] = parts[i]
        })
      }

      return o
    }

    let r, ptn, fn, res
    for (let i = 0; i < this.routes.length; i++) {
      r = this.routes[i]
      ptn = r.ptn
      fn = r.fn
      if (typeof ptn === 'string' && href === ptn) {
        res = fn(result(href, r))
        if (res !== true) {
          break
        }
      } else if (ptn instanceof RegExp && ptn.test(href)) {
        res = fn(result(href, r))
        if (res !== true) {
          break
        }
      }
    }
  }

  start () {
    window.addEventListener('popstate', (e) => {
      if (e.state && e.state.href) {
        this.route(e.state.href, e.state.state)
      } else {
        this.route(this.hash)
      }
    })
    if (!this.hash) {
      if (this.defaultHref) {
        this.replace(this.defaultHref, true)
      }
    } else {
      this.route(this.hash)
    }
  }
}

export default Router
