/* globals document fetch */

var API = '/api/'
var ROOT = document.getElementById('containers')

var renderers = {
  containers: containers,
  applications: applications
}

function containers (parent, doc) {
  parent.innerHTML = ''
  let ul = document.createElement('ul')
  let i, item
  let html = ''
  for (i = 0; i < doc.containers.length; i++) {
    item = doc.containers[i]
    html += `<li><h4>${item.name}</h4><p>${item.description}</p>`
    html += `<a class="api-link" href="#" data-renderer="applications" data-method="get", data-url="${API}${item.name}/">GET ${API}${item.name}/</a>`
    // html += `<nav><a href="#raw">Raw</a></nav>`
    // html += `<pre>${JSON.stringify(item, undefined, 2)}</pre>`
    html += `<div></div>`
    html += `</li>`
  }
  ul.innerHTML = html
  parent.appendChild(ul)

  let links = ul.querySelectorAll('.api-link')
  for (i = 0; i < links.length; i++) {
    links[i].addEventListener('click', (e) => {
      e.preventDefault()
      let el = e.currentTarget
      let url = el.getAttribute('data-url')
      let method = el.getAttribute('data-method')
      let renderer = el.getAttribute('data-renderer')
      get(url, {method: method}, renderers[renderer].bind(null, el.parentNode.lastChild, url))
    })
  }
}

function applications (parent, url, doc) {
  parent.innerHTML = ''
  let ul = document.createElement('ul')
  let i, item
  let html = ''
  console.log(doc)
  for (i = 0; i < doc.length; i++) {
    item = doc[i]
    html += `<li><h4>${item.name}</h4><p>${item.description}</p>`
    html += `<a class="api-link" href="#" data-method="get", data-url="${url}${item.name}/">GET ${url}${item.name}/</a>`
    // html += `<nav><a href="#raw">Raw</a></nav>`
    // html += `<pre>${JSON.stringify(item, undefined, 2)}</pre>`
    html += `</li>`
  }
  ul.innerHTML = html
  parent.appendChild(ul)
}

function get (url, options, renderer) {
  return fetch(url, options)
    .then((res) => res.json().then((doc) => renderer(doc)))
    .catch((err) => console.error(err))
}

var link = document.querySelector('.api-link')

link.addEventListener('click', (e) => {
  e.preventDefault()
  get(API, {}, renderers.containers.bind(null, ROOT))
})
