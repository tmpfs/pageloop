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
    return API + `apps/${rpc.parameters.container}/`
  },
  'Template.List': function () {
    return API + 'templates/'
  },
  'Job.ActiveJobs': function () {
    return API + 'jobs/'
  },
  'Application.Read': function (rpc) {
    return API + `apps/${rpc.parameters.container}/${rpc.parameters.name}`
  },
  'Application.ReadFiles': function (rpc) {
    return API + `apps/${rpc.parameters.container}/${rpc.parameters.name}/files/`
  },
  'Application.ReadPages': function (rpc) {
    return API + `apps/${rpc.parameters.container}/${rpc.parameters.name}/pages/`
  },
  'Application.DeleteFiles': function (rpc) {
    return API + `apps/${rpc.parameters.container}/${rpc.parameters.name}/files/`
  },
  'Application.RunTask': function (rpc) {
    return API + `apps/${rpc.parameters.container}/${rpc.parameters.name}/tasks/${rpc.parameters.url}`
  },
  'Application.Delete': function (rpc) {
    return API + `apps/${rpc.parameters.container}/${rpc.parameters.name}`
  },
  'File.Create': function (rpc) {
    return API + `apps/${rpc.parameters.owner.container}/${rpc.parameters.owner.name}/files${rpc.parameters.url}`
  },
  'File.CreateTemplate': function (rpc) {
    return API + `apps/${rpc.parameters.owner.container}/${rpc.parameters.owner.name}/files${rpc.parameters.url}`
  },
  'File.Save': function (rpc) {
    return API + `apps/${rpc.parameters.owner.container}/${rpc.parameters.owner.name}/files${rpc.parameters.url}`
  },
  'File.Move': function (rpc) {
    return API + `apps/${rpc.parameters.owner.container}/${rpc.parameters.owner.name}/files${rpc.parameters.url}`
  },
  'File.ReadSource': function (rpc) {
    return API + `apps/${rpc.parameters.owner.container}/${rpc.parameters.owner.name}/src${rpc.parameters.url}`
  },
  'File.ReadSourceRaw': function (rpc) {
    return API + `apps/${rpc.parameters.owner.container}/${rpc.parameters.owner.name}/raw${rpc.parameters.url}`
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
  'Core.Meta': getDefaultOptions,                 // v3
  'Core.Stats': getDefaultOptions,                // v3
  'Host.List': getDefaultOptions,                 // v3
  'Container.CreateApp': getPutOptions,           // v3
  'Template.List': getDefaultOptions,             // v3
  'Job.ActiveJobs': getDefaultOptions,            // v3 - requires testing
  'Application.Read': getDefaultOptions,          // v3
  'Application.ReadFiles': getDefaultOptions,     // v3
  'Application.ReadPages': getDefaultOptions,     // v3
  'Application.DeleteFiles': getDeleteOptions,    // v3
  'Application.RunTask': getPutOptions,           // v3
  'Application.Delete': getDeleteOptions,         // v3
  'File.Create': getPutOptions,                   // v3
  'File.CreateTemplate': getPutOptions,           // v3
  'File.Save': getPostOptions,                    // v3
  'File.Move': (rpc) => {                         // v3
    const o = getPostOptions(rpc)
    o.headers.Location = rpc.args[0]
    return o
  },
  'File.ReadSource': getBinaryOptions,            // v3
  'File.ReadSourceRaw': getBinaryOptions          // v3
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
