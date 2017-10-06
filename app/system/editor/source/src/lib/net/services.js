// REST API endpoint
const API = '/api/'

function getBodyOptions (rpc, options) {
  options.headers = options.headers || {}
  // Use input raw request body
  if (rpc.request) {
    let body = rpc.request.value

    // Sending JSON body
    if (rpc.request.type === 0) {
      body = JSON.stringify(body)
    }

    options.body = body

    if (rpc.request.type === 0) {
      options.headers['Content-Type'] = 'application/json; charset=utf-8'
    } else {
      options.headers['Content-Type'] =
        rpc.request.mime || 'application/octet-stream'
    }

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

// TODO: deprecate this function: use getFileRefUrl() for everything
function getFileUrl (params, filter) {
  const u = API + `apps/${params.container}/${params.application}/${filter}${params.url}`
  delete params.container
  delete params.application
  return u
}

function parseFileRef (ref) {
  if (!ref) {
    throw new Error('Asset reference is empty, cannot create URL')
  }
  const u = new URL(ref)
  const parts = u.pathname.replace(/^\//, '').split('/')
  const container = parts[0]
  const application = parts[1]
  const url = u.hash.replace(/^#/, '')
  return {
    container: container,
    application: application,
    url: url
  }
}

function getFileRefUrl (params, filter) {
  const {container, application, url} = parseFileRef(params.ref)
  return API + `apps/${container}/${application}/${filter}${url}`
}

function getArchiveUrl (params) {
  let url = API + `apps/${params.container}/${params.name}/zip/`
  if (params.filter) {
    url += `${params.filter}/`
  }
  return url
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
  'Container.Read': (rpc, params) => {
    return {
      url: API + `apps/${params.name}/`,
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
  'Service.List': (rpc, params) => {
    return {
      url: API + 'services/',
      options: getDefaultOptions(rpc)
    }
  },
  'Service.Read': (rpc, params) => {
    return {
      url: API + `services/${params.service}/`,
      options: getDefaultOptions(rpc)
    }
  },
  'Service.ReadMethod': (rpc, params) => {
    return {
      url: API + `services/${params.service}/${params.method}/`,
      options: getDefaultOptions(rpc)
    }
  },
  'Service.ReadMethodCalls': (rpc, params) => {
    return {
      url: API + `services/${params.service}/${params.method}/calls/`,
      options: getDefaultOptions(rpc)
    }
  },
  'Job.List': (rpc, params) => {
    return {
      url: API + 'jobs/',
      options: getDefaultOptions(rpc)
    }
  },
  'Job.Read': (rpc, params) => {
    return {
      url: API + `jobs/${params.id}/`,
      options: getDefaultOptions(rpc)
    }
  },
  'Job.Delete': (rpc, params) => {
    return {
      url: API + `jobs/${params.id}/`,
      options: getDeleteOptions(rpc)
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
      url: API + `apps/${params.container}/${params.name}/tasks/${params.task}`,
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
      url: getFileRefUrl(params, 'files'),
      options: o
    }
  },
  'File.Read': (rpc, params) => {
    return {
      url: getFileRefUrl(params, 'files'),
      options: getDefaultOptions(rpc)
    }
  },
  'File.ReadSource': (rpc, params) => {
    return {
      url: getFileRefUrl(params, 'src'),
      options: getDefaultOptions(rpc)
    }
  },
  'File.ReadSourceRaw': (rpc, params) => {
    return {
      url: getFileRefUrl(params, 'raw'),
      options: getDefaultOptions(rpc)
    }
  },
  'Archive.Export': (rpc, params) => {
    return {
      url: getArchiveUrl(params),
      options: getDefaultOptions(rpc)
    }
  }
}

function fetchFromRpc (rpc) {
  const o = {}
  const fn = services[rpc.method]
  if (typeof (fn) !== 'function') {
    throw new Error(`No client definition for service method ${rpc.method}`)
  }
  const {url, options} = fn(rpc, rpc.parameters)
  o.url = url
  o.options = options || {}
  o.options.headers = o.options.headers || {}

  // Hint for optimized route lookup
  o.options.headers['X-Method-Name'] = rpc.method
  o.options.headers['X-Method-Seq'] = rpc.id
  return o
}

export {fetchFromRpc}
