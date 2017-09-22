// REST API endpoint
const API = '/api/'

function getBodyOptions (rpc, options) {
  options.headers = options.headers || {}

  console.log('get body options')
  console.log(rpc.rawRequest)

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

function getFileUrl (params, filter) {
  const u = API + `apps/${params.container}/${params.application}/${filter}${params.url}`
  delete params.container
  delete params.application
  return u
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
  const {url, options} = services[rpc.method](rpc, rpc.parameters)
  o.url = url
  o.options = options || {}
  o.options.headers = o.options.headers || {}

  // Hint for optimized route lookup
  o.options.headers['X-Method-Name'] = rpc.method
  o.options.headers['X-Method-Seq'] = rpc.id
  return o
}

export {fetchFromRpc}
