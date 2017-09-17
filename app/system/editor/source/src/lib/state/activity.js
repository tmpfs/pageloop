class Activity {
  constructor (log) {
    this.notification = []
    this.network = []
    this.log = log
  }

  add (info) {
    info.time = Date.now()
    this.notification.unshift(info)
  }
}

export default Activity
