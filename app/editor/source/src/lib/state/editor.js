import ViewState from './view'

class EditorState extends ViewState {
  constructor () {
    super()
    this.defaultView = 'code-editor'
  }
}

export default EditorState
