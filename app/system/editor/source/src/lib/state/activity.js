class Activity {
  constructor () {
    this.notification = []
    this.log = []
    this.network = []
  }

  addNotificationActivity (info) {
    info.time = Date.now()
    this.notification.unshift(info)
  }
}

export default Activity
