import ViewState from './view'

/**
 *  Represents the state of the sidebar.
 */
class SidebarState extends ViewState {
  constructor () {
    super()
    // List selections
    this.pages = []
    this.files = []
    this.media = []
  }
}

export default SidebarState
