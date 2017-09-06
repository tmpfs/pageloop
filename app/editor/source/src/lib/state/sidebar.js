/**
 *  Represents the selections in the sidebar file lists.
 */
class SidebarState {
  constructor () {
    // Curent sidebar view
    this.view = ''

    // List selections
    this.pages = []
    this.files = []
    this.media = []
  }
}

export default SidebarState
