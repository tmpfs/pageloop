/* globals document fetch */

(function () {
  const PUT = 'PUT'
  const POST = 'POST'

  function send (url, req, cb) {
    req.url = url
    return fetch(url, req)
      .then((res) => res.json().then((doc) => cb(null, req, res, doc)))
      .catch((err) => cb(err))
  }

  function log (o) {
    let s = o.toString()
    let logger = document.getElementById('log')
    logger.innerText += s
  }

  function onSubmit (e) {
    e.preventDefault()
    let form = e.currentTarget
    let url = document.getElementById('url').value
    let group = form.querySelectorAll('input[type="radio"]')
    let radio
    for (let i = 0; i < group.length; i++) {
      if (group[i].checked) {
        radio = group[i]
        break
      }
    }
    let method = radio.value
    let options = {method: method}

    let data = document.getElementById('data').value
    if ((method === PUT || method === POST) && data !== '') {
      try {
        // check JSON is valid before sending
        data = JSON.parse(data)
        data = JSON.stringify(data)
      } catch (e) {
        return log(e)
      }

      options.body = data
    }

    log(`${options.method} ${url}\n`)

    send(url, options, (err, req, res, doc) => {
      if (err) {
        return console.error(err)
      }

      log(`${res.status} ${req.url}\n`)

      let response = document.getElementById('response')
      response.classList.remove('hidden')
      let pre = response.querySelector('pre')
      let status = response.querySelector('.status')
      status.innerText = res.status

      let str = JSON.stringify(doc, undefined, 2)
      pre.innerText = str
    })
  }

  let form = document.getElementById('debug')
  form.addEventListener('submit', onSubmit)
})()
