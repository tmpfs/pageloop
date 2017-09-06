import ColumnManager from './columns'
import ViewState from './view'

class EditorState extends ViewState {
  constructor () {
    super()
    this.defaultView = 'code-editor'

    // State for edit mode columns
    this.columns = new ColumnManager()
  }
}

export default EditorState
