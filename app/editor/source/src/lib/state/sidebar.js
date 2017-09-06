import ViewState from './view'

// Represents the state of the sidebar.
class SidebarState extends ViewState {
  constructor () {
    super()
    // List selections
    this.pages = []
    this.files = []
    this.media = []
  }

  // Get a selection list for the current view.
  get selection () {
    if (!this.view) {
      return []
    }
    return this[this.view]
  }

  set selection (val) {
    if (!this.view) {
      return
    }
    if (!val) {
      val = []
    }
    this[this.view] = val
  }
}

export default SidebarState
