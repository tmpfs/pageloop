class Activity {
  constructor () {
    this.notifications = []
  }

  addNotificationActivity (info) {
    info.time = Date.now()
    this.notifications.unshift(info)
  }
}

export default Activity
