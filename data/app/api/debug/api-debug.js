/* globals document fetch */

function send (url, options, cb) {
  return fetch(url, options)
    .then((res) => res.json().then((doc) => cb(null, res, doc)))
    .catch((err) => cb(err))
}

/*
var link = document.querySelector('.api-link')

link.addEventListener('click', (e) => {
  e.preventDefault()
  get(API, {}, renderers.containers.bind(null, ROOT))
})
*/

function onSubmit (e) {
  e.preventDefault()
  var form = e.currentTarget
  console.log(e)
  var url = document.getElementById('url').value
  var group = form.querySelectorAll('input[type="radio"]')
  var radio
  for (let i = 0; i < group.length; i++) {
    if (group[i].checked) {
      radio = group[i]
      break
    }
  }
  var method = radio.value
  var options = {method: method}

  console.log(options)

  send(url, options, (err, res, doc) => {
    if (err) {
      return console.error(err)
    }

    console.log(res)

    var response = document.getElementById('response')
    response.classList.remove('hidden')
    var pre = response.querySelector('pre')
    var status = response.querySelector('.status')
    status.innerText = res.status

    //console.log(doc)

    var str = JSON.stringify(doc, undefined, 2)
    pre.innerText = str
  })
}

var form = document.getElementById('debug')
form.addEventListener('submit', onSubmit)
