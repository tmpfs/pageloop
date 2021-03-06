/**
 *  Maintains a list of notifications used by the UI to
 *  render notifications which timeout by default.
 */

let counter = 0

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

    info.id = (++counter)

    this.notifications.unshift(info)
  }

  getById (id) {
    for (let i = 0; i < this.notifications.length; i++) {
      if (id === this.notifications[i].id) {
        return this.notifications[i]
      }
    }
  }
}

export default Notifier
