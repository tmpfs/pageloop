<template>
  <div
    :class="{maximized: maximized === 'editor', minimized: maximized != '' && maximized !== 'editor'}"
    class="editor">
    <div class="column-header">
      <h2>Editor</h2>
      <div class="column-options">
        <nav class="tabs">
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
    </div>
    <nav class="toolbar clearfix">
      <h2><span class="status-dirty" :class="{hidden: !isDirty}">✺</span>{{currentFile.name}}</h2>
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

<style scoped>
  .column-header {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    border-bottom: 1px solid var(--border-color);
    user-select: none;
    overflow: hidden;
  }

  .column-header > * {
    flex: 1 0;
  }

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

  .status-dirty {
    margin-right: 0.5rem;
    color: var(--orange-color);
  }

  .file-editor, .source-editor, .visual-editor {
    height: calc(100% - 4.6rem);
  }

</style>
