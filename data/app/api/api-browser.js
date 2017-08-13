/* globals document fetch */

var API = '/api/'
var ROOT = document.getElementById('containers')

var renderers = {
  containers: containers,
  applications: applications,
  application: application
}

function events (target) {
  let i
  let links = target.querySelectorAll('.api-link')
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

function containers (parent, doc) {
  parent.innerHTML = ''
  let ul = document.createElement('ul')
  let i, item
  let html = ''
  for (i = 0; i < doc.containers.length; i++) {
    item = doc.containers[i]
    html += `<li>`
    html += `<h4>${item.name}</h4>`
    html += `<span class="label">Container</span>`
    if (item.description) {
      html += `<p>${item.description}</p>`
    }
    html += `<a class="api-link" href="#" data-renderer="applications" data-method="get", data-url="${API}${item.name}/">GET ${API}${item.name}/</a>`
    // html += `<nav><a href="#raw">Raw</a></nav>`
    // html += `<pre>${JSON.stringify(item, undefined, 2)}</pre>`
    html += `<div></div>`
    html += `</li>`
  }
  ul.innerHTML = html
  parent.appendChild(ul)

  events(ul)
}

function applications (parent, url, doc) {
  parent.innerHTML = ''
  let ul = document.createElement('ul')
  let i, item
  let html = ''
  for (i = 0; i < doc.length; i++) {
    item = doc[i]
    html += `<li>`
    html += `<h4>${item.name}</h4>`
    html += `<span class="label">Application</span>`
    if (item.description) {
      html += `<p>${item.description}</p>`
    }
    html += `<a class="api-link" href="#" data-renderer="application" data-method="get", data-url="${url}${item.name}/files/">GET ${url}${item.name}/files/</a>`
    // html += `<nav><a href="#raw">Raw</a></nav>`
    // html += `<pre>${JSON.stringify(item, undefined, 2)}</pre>`
    html += `<div></div>`
    html += `</li>`
  }
  ul.innerHTML = html
  parent.appendChild(ul)

  events(ul)
}

function application (parent, url, doc) {
  parent.innerHTML = ''
  let ul = document.createElement('ul')
  let i, item, next
  let html = ''
  for (i = 0; i < doc.length; i++) {
    item = doc[i]
    next = url.replace(/\/$/, '') + item.url
    html += `<li>`
    html += `<h4>${item.name} ${item.dir ? '' : '(' + item.size + ' bytes)'}</h4>`
    html += `<span class="label">${item.dir ? 'Dir' : 'File'}</span>`
    html += `<a class="api-link" href="#" data-renderer="files" data-method="get", data-url="${next}">GET ${next}</a>`
    // html += `<nav><a href="#raw">Raw</a></nav>`
    // html += `<pre>${JSON.stringify(item, undefined, 2)}</pre>`
    html += `<div></div>`
    html += `</li>`
  }
  ul.innerHTML = html
  parent.appendChild(ul)

  // events(ul)
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
