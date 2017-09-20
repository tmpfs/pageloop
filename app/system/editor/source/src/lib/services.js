// REST API endpoint
const API = '/api/'

function getBodyOptions (rpc, options) {
  options.headers = options.headers || {}
  if (rpc.body) {
    let body = rpc.body
    if (!rpc.raw) {
      body = JSON.stringify(rpc.body)
    }
    options.body = body
    options.headers['Content-Type'] =
      rpc.mime || 'application/json; charset=utf-8'
    options.headers['Content-Length'] = body.length
  }
  return options
}

function getDefaultOptions () {
  return {method: 'GET'}
}

function getPutOptions (rpc) {
  const o = {method: 'PUT'}
  return getBodyOptions(rpc, o)
}

function getPostOptions (rpc) {
  const o = {method: 'POST'}
  return getBodyOptions(rpc, o)
}

function getDeleteOptions (rpc) {
  const o = {method: 'DELETE'}
  return getBodyOptions(rpc, o)
}

// Maps RPC function names to REST request URLs
const urls = {
  'Core.Meta': function () {
    return API
  },
  'Core.Stats': function () {
    return API + 'stats/'
  },
  'Host.List': function () {
    return API + 'apps/'
  },
  'Container.CreateApp': function (rpc) {
    return API + `apps/${rpc.parameters.context}/`
  },
  'Template.List': function () {
    return API + 'templates/'
  },
  'Job.ActiveJobs': function () {
    return API + 'jobs/'
  },
  'Application.Read': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}`
  },
  'Application.ReadFiles': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files/`
  },
  'Application.ReadPages': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/pages/`
  },
  'Application.DeleteFiles': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files/`
  },
  'Application.RunTask': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/tasks/${rpc.parameters.item}`
  },
  'Application.Delete': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}`
  },
  'File.Create': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files${rpc.parameters.item}`
  },
  'File.CreateTemplate': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files${rpc.parameters.item}`
  },
  'File.Save': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files${rpc.parameters.item}`
  },
  'File.Move': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files${rpc.parameters.item}`
  },
  'File.ReadSource': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/src${rpc.parameters.item}`
  },
  'File.ReadSourceRaw': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/raw${rpc.parameters.item}`
  }
}

// Get a binary file response.
function getBinaryOptions () {
  const o = getDefaultOptions()
  o.headers = {
    'Accept': 'application/octet-stream'
  }
  return o
}

const options = {
  'Core.Meta': getDefaultOptions,
  'Core.Stats': getDefaultOptions,
  'Host.List': getDefaultOptions,
  // TODO: restore create app from template
  'Container.CreateApp': getPutOptions,
  'Template.List': getDefaultOptions,
  'Job.ActiveJobs': getDefaultOptions,
  'Application.Read': getDefaultOptions,
  'Application.ReadFiles': getDefaultOptions,
  'Application.ReadPages': getDefaultOptions,
  'Application.DeleteFiles': getDeleteOptions,
  'Application.RunTask': getPutOptions,
  'Application.Delete': getDeleteOptions,
  'File.Create': getPutOptions,
  'File.CreateTemplate': getPutOptions,
  'File.Save': getPostOptions,
  'File.Move': (rpc) => {
    const o = getPostOptions(rpc)
    o.headers.Location = rpc.args[0]
    return o
  },
  'File.ReadSource': getBinaryOptions,
  'File.ReadSourceRaw': getBinaryOptions
}

function fetchFromRpc (rpc) {
  const o = {}
  o.url = urls[rpc.method](rpc)
  o.options = options[rpc.method](rpc)
  o.options.headers = o.options.headers || {}

  // Hint for optimized route lookup
  o.options.headers['X-Method-Name'] = rpc.method
  o.options.headers['X-Method-Seq'] = rpc.id

  if (rpc.fetch) {
    o.options.raw = true
  }
  return o
}

export {fetchFromRpc}
