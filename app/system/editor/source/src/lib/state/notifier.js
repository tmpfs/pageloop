/**
 *  Maintains a list of notifications used by the UI to
 *  render notifications which timeout by default.
 */
class Notifier {
  constructor () {
    this.notifications = []
  }

  notify (info, del) {
    if (del) {
      for (let i = 0; i < this.notifications.length; i++) {
        if (info === this.notifications[i]) {
          this.notifications.splice(i, 1)
          break
        }
      }
      return
    }

    info.reveal = true

    this.notifications.unshift(info)
  }
}

export default Notifier
