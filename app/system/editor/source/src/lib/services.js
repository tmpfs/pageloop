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
  'Container.List': function () {
    return API + 'apps/'
  },
  'Container.CreateApp': function (rpc) {
    return API + `apps/${rpc.parameters.context}/`
  },
  'Template.ReadApplications': function () {
    return API + 'templates/'
  },
  'Jobs.ReadActiveJobs': function () {
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
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files/${rpc.parameters.item}`
  },
  'File.CreateTemplate': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files/${rpc.parameters.item}`
  },
  'File.Save': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files/${rpc.parameters.item}`
  },
  'File.Move': function (rpc) {
    return API + `apps/${rpc.parameters.context}/${rpc.parameters.target}/files/${rpc.parameters.item}`
  }
}

const options = {
  'Core.Meta': getDefaultOptions,
  'Core.Stats': getDefaultOptions,
  'Container.List': getDefaultOptions,
  // TODO: restore create app from template
  'Container.CreateApp': getPutOptions,
  'Template.ReadApplications': getDefaultOptions,
  'Jobs.ReadActiveJobs': getDefaultOptions,
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
  }
}

export {urls, options}
