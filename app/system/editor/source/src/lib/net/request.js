const TypeJson = 0
const TypeByte = 1

// Message identifier counter
let id = 0

class Request {
  constructor (id, method, params) {
    this.id = id
    this.method = method
    if (params && !Array.isArray(params)) {
      params = [params]
    }
    this.params = params || []

    Object.defineProperty(this, 'request', {
      enumerable: false,
      configurable: false,
      writable: true
    })
  }

  // Set an object to be serialized to JSON and
  // sent as the request body (REST only).
  json (value) {
    this.request = {
      type: TypeJson,
      value: value
    }
  }

  // Set a raw body for the request (REST only).
  body (value, mime) {
    this.request = {
      type: TypeByte,
      value: value,
      mime: mime
    }
  }

  get parameters () {
    return this.params[0]
  }

  // Get a JSON RPC request object.
  static rpc (method, params) {
    return new Request(++id, method, params)
  }
}

export {Request, TypeJson, TypeByte}
