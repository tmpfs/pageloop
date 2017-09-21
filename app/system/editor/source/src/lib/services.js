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

function getFileUrl (params, filter) {
  return API + `apps/${params.owner.container}/${params.owner.name}/${filter}${params.url}`
}

// Maps RPC service method names to REST request URLs and fetch options.
const services = {
  'Core.Meta': (rpc, params) => {
    return {
      url: API,
      options: getDefaultOptions(rpc)
    }
  },
  'Core.Stats': (rpc, params) => {
    return {
      url: API + 'stats/',
      options: getDefaultOptions(rpc)
    }
  },
  'Host.List': (rpc, params) => {
    return {
      url: API + 'apps/',
      options: getDefaultOptions(rpc)
    }
  },
  'Container.CreateApp': (rpc, params) => {
    return {
      url: API + `apps/${params.container}/`,
      options: getPutOptions(rpc)
    }
  },
  'Template.List': (rpc, params) => {
    return {
      url: API + 'templates/',
      options: getDefaultOptions(rpc)
    }
  },
  'Job.ActiveJobs': (rpc, params) => {
    return {
      url: API + 'jobs/',
      options: getDefaultOptions(rpc)
    }
  },
  'Application.Read': (rpc, params) => {
    return {
      url: API + `apps/${params.container}/${params.name}`,
      options: getDefaultOptions(rpc)
    }
  },
  'Application.ReadFiles': (rpc, params) => {
    return {
      url: API + `apps/${params.container}/${params.name}/files/`,
      options: getDefaultOptions(rpc)
    }
  },
  'Application.ReadPages': (rpc, params) => {
    return {
      url: API + `apps/${params.container}/${params.name}/pages/`,
      options: getDefaultOptions(rpc)
    }
  },
  'Application.DeleteFiles': (rpc, params) => {
    return {
      url: API + `apps/${params.container}/${params.name}/files/`,
      options: getDeleteOptions(rpc)
    }
  },
  'Application.RunTask': (rpc, params) => {
    return {
      url: API + `apps/${params.container}/${params.name}/tasks/${params.url}`,
      options: getPutOptions(rpc)
    }
  },
  'Application.Delete': (rpc, params) => {
    return {
      url: API + `apps/${params.container}/${params.name}`,
      options: getDeleteOptions(rpc)
    }
  },
  'File.Create': (rpc, params) => {
    return {
      url: getFileUrl(params, 'files'),
      options: getPutOptions(rpc)
    }
  },
  'File.CreateTemplate': (rpc, params) => {
    return {
      url: getFileUrl(params, 'files'),
      options: getPutOptions(rpc)
    }
  },
  'File.Save': (rpc, params) => {
    return {
      url: getFileUrl(params, 'files'),
      options: getPostOptions(rpc)
    }
  },
  'File.Move': (rpc, params) => {
    const o = getPostOptions(rpc)
    o.headers.Location = params.destination
    return {
      url: getFileUrl(params, 'files'),
      options: o
    }
  },
  'File.ReadSource': (rpc, params) => {
    return {
      url: getFileUrl(params, 'src'),
      options: getDefaultOptions(rpc)
    }
  },
  'File.ReadSourceRaw': (rpc, params) => {
    return {
      url: getFileUrl(params, 'raw'),
      options: getDefaultOptions(rpc)
    }
  }
}

function fetchFromRpc (rpc) {
  const o = {}
  const {url, opts} = services[rpc.method](rpc, rpc.parameters)
  o.url = url
  o.options = opts || {}

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
