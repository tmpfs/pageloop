<template>
  <div
    :class="{maximized: maximized === 'editor', minimized: maximized != '' && maximized !== 'editor'}"
    class="editor">
    <div class="column-header">
      <div>
        <h2>Editor</h2>
      </div>
      <nav class="toolbar clearfix">
        <a @click="save"
          v-bind:class="{hidden: currentView != 'code-editor'}" href="#" title="Save & Run">Save & Run</a>
        <a
          @click="maximized = 'editor'"
          :class="{hidden: maximized === 'editor'}"
          title="Maximize">◩</a>
        <a
          @click="maximized = ''"
          :class="{hidden: maximized !== 'editor'}"
          title="Minimize">▣</a>
      </nav>
    </div>
    <div class="column-options">
      <nav class="tabs">
        <h3><span class="status-dirty" :class="{hidden: !isDirty}">✺</span>{{currentFile.name}}</h3>
        <a v-bind:class="{selected: currentView === 'file-editor', disabled: fileHidden}"
          @click="currentView = 'file-editor'"
          title="Show file editor">File</a>
        <a v-bind:class="{selected: currentView === 'data-editor', disabled: dataHidden}"
          @click="currentView = 'data-editor'"
          title="Show data editor">Data</a>
        <a v-bind:class="{selected: currentView === 'code-editor', disabled: hidden}"
          @click="currentView = 'code-editor'"
          title="Show source editor">Code</a>
        <a v-bind:class="{selected: currentView === 'visual-editor', disabled: hidden}"
          @click="currentView = 'visual-editor'"
          title="Show visual editor">Visual</a>
      </nav>
    </div>
    <component v-bind:is="currentView"></component>
    <div class="column-drag" :class="{hidden: maximized}" @mousedown="resizeColumn">&nbsp;</div>
  </div>
</template>

<script>

import Welcome from '@/components/editor/Welcome'
import FileEditor from '@/components/editor/FileEditor'
import DataEditor from '@/components/editor/DataEditor'
import CodeEditor from '@/components/editor/CodeEditor'
import VisualEditor from '@/components/editor/VisualEditor'

export default {
  name: 'editor',
  computed: {
    dirty: {
      get: function () {
        return this.currentFile.dirty
      },
      set: function (val) {
        if (this.changeGeneration > -1 && this.currentFile.document) {
          if (val === true && this.currentFile.document.isClean(this.changeGeneration)) {
            val = false
          }
        }
        this.$store.commit('current-file-dirty', val)
        this.isDirty = val
      }
    },
    maximized: {
      get: function () {
        return this.$store.state.columns.maximized
      },
      set: function (val) {
        this.$store.commit('maximize-column', val)
      }
    },
    dataHidden: function () {
      return this.hidden || !this.$store.state.current.page
    },
    fileHidden: function () {
      return !this.$store.state.hasFile()
    },
    hidden: function () {
      return !this.$store.state.hasFile() || this.$store.state.isDirectory()
    },
    currentFile: function () {
      return this.$store.state.app.current
    },
    currentView: {
      get: function () {
        return this.$store.state.editorView
      },
      set: function (view) {
        this.$store.commit('editor-view', view)
      }
    }
  },
  watch: {
    currentFile: function (file) {
      this.title = file.url
      this.dirty = file.dirty
      if (file && file.dir) {
        this.currentView = 'file-editor'
      }
    }
  },
  data: function () {
    return {
      title: '',
      isDirty: false,
      changeGeneration: -1,
      codeMirror: null
    }
  },
  methods: {
    save: function (e) {
      if (e) {
        e.preventDefault()
      }
      this.$store.dispatch('save-file')
        .then(() => {
          this.dirty = false
          this.changeGeneration = this.currentFile.document.changeGeneration()
        })
        .catch((e) => console.error(e))
    },
    resizeColumn: function (e) {
      this.$store.dispatch('resize-column', e)
    }
  },
  components: {Welcome, FileEditor, DataEditor, CodeEditor, VisualEditor}
}
</script>

<style>
  /* Editor global styles */
  .sidebar, .editor, .preview {
    flex: 1 1 auto;
    opacity: 1;
    transition: all 0.5s ease-out;
    background: var(--background);
  }

  .sidebar {
    max-width: 32rem;
    width: 20%;
  }

  .sidebar, .editor {
    position: relative;
  }

  .sidebar:not(.minimized), .editor:not(.minimized) {
    min-width: var(--drag-size);
  }

  .sidebar:not(.maximized) > :not(.column-drag), .editor:not(.maximized) > :not(.column-drag) {
    /*
     * Setting padding on sidebar/editor causes issues when maximizing
     * so we set on the child elements instead.
     */
    margin-right: var(--drag-size);
  }

  .sidebar > .column-drag, .editor > .column-drag {
    position: absolute;
    top: 0;
    left: calc(100% - var(--drag-size));
    bottom: 0;
    right: 0;
    width: var(--drag-size);
    height: 100%;
    cursor: ew-resize;
    background: var(--border-color);
  }

  .editor, .preview {
    width: 40%;
  }

  .file-editor, .source-editor, .visual-editor {
    height: calc(100% - 4.6rem);
  }

  .column-header {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    border-bottom: 1px solid var(--border-color);
    user-select: none;
    overflow: hidden;
    height: 2.2rem;
  }

  .column-header > * {
    flex: 1 0;
  }

  .column-header h2 {
    display: inline-block;
    font-size: 1.4rem;
    text-transform: uppercase;
    padding: 0 1rem;
  }

  .column-options {
    border-bottom: 1px solid var(--border-color);
  }

  h2.file-info {
    font-size: 2.2rem;
    padding: 0 0 2rem 0;
  }

  .column-options > .toolbar {
    display: inline-block;
  }

  .toolbar {
    height: 2.3rem;
    font-size: 1.3rem;
    text-transform: uppercase;
    text-align: right;
    /*font-size: 1.4rem;*/
    text-align: right;
    user-select: none;
    overflow: hidden;
  }

  .toolbar a {
    display: inline-block;
    text-align: center;
    color: var(--toolbar-link);
    /*
    width: 2.3rem;
    height: 2.3rem;
    */
    padding: 0 0.5rem;
  }

  .toolbar a:hover {
    background: var(--base03-color);
    color: var(--base3-color);
  }

  h3 {
    font-size: 1.4rem;
    text-transform: none;
    float: left;
    border: 0;
    padding: 0;
    margin: 0 0 0 1rem;
  }

  .editor h3 {
    width: 50%;
  }

  .tabs {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    padding: 0;
    font-size: 1.4rem;
    text-transform: uppercase;
    user-select: none;
  }

  .tabs > a {
    flex: 1 0;
    color: currentColor;
    text-align: center;
    background: transparent;
  }

  .tabs > .selected {
    text-decoration: underline;
  }

  .tabs > a:not(:first-child) {
    border-left: 1px solid var(--border-color);
  }

  .tabs > a:hover, .tabs > a.selected {
    background: var(--base03-color);
    color: var(--base3-color);
    transition: all 0.5s ease-out;
  }

</style>

<style scoped>
  /* Editor scoped styles */
  .editor > .column-header .tabs > a, .preview > .column-header .tabs > a {
    border-left: 1px solid var(--border-color);
  }

  .preview {
    position: relative;
  }

  .sidebar .scroll {
    /*padding-top: 1rem;*/
    height: calc(100% - 4.3rem);
  }

  .welcome.scroll {
    padding: 0 2rem;
  }

  .maximized {
    width: 100%;
    max-width: none;
  }

  .minimized {
    width: 0%;
    opacity: 0;
    pointer-events: none;
  }

  .minimized > .column-drag {
    display: none;
  }

  .status-dirty {
    margin-right: 0.5rem;
    color: var(--orange-color);
  }

</style>
