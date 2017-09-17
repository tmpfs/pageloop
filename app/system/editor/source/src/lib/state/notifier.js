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
      console.log('trying to delete notification')
      for (let i = 0; i < this.notifications.length; i++) {
        if (info === this.notifications[i]) {
          this.notifications.splice(i, 1)
          console.log('deleting notification')
          break
        }
      }
      return
    }

    info.id = (++counter)

    this.notifications.unshift(info)
  }
}

export default Notifier
