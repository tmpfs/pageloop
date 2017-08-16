/* global fetch */
(function () {
  var url = '/api/'
  function request (url, options, cb) {
    return fetch(url, options)
      .then((res) => res.json().then((doc) => cb(null, doc)))
      .catch((err) => cb(err))
  }

  function template (selector) {
    var tpl = document.querySelector('template')
    tpl = tpl && tpl.content ? tpl.content : tpl
    tpl = tpl.querySelector(selector)
    tpl = tpl.cloneNode(true)
    return tpl
  }

  request(url, null, (err, doc) => {
    if (err) {
      return console.error(err)
    }
    var parent = document.getElementById('containers')
    var tpl
    var list
    var app
    doc.forEach((val) => {
      tpl = template('.container')

      // render containers
      tpl.querySelector('.name').appendChild(
        document.createTextNode(val.name))
      tpl.querySelector('.description').appendChild(
        document.createTextNode(val.description))
      parent.appendChild(tpl)

      // render applications
      list = tpl.querySelector('.applications')
      if (val.apps) {
        val.apps.forEach((application) => {
          app = template('.app')
          app.querySelector('.name').appendChild(
            document.createTextNode(application.name))
          app.querySelector('.url').appendChild(
            document.createTextNode(application.url))

          app.querySelector('nav .view').setAttribute('href', application.url)

          list.appendChild(app)
        })
      }
    })
  })
})()
