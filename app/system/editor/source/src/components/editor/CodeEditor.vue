<template>
  <div class="source-editor">
    <hint id="file-save" class="compact"></hint>
    <div class="text-editor"></div>
  </div>
</template>
<script>
/* globals CodeMirror */

import Hint from '@/components/Hint'

export default {
  name: 'code-editor',
  data: function () {
    return {
      mirror: null
    }
  },
  computed: {
    currentFile: function () {
      return this.$store.state.app.current
    },
    value: {
      get: function () {
        return this.$store.state.current.content
      },
      set: function (val) {
        this.$store.state.current.content = val
      }
    }
  },
  watch: {
    currentFile: function (file) {
      if (file && file.mime) {
        this.setCodeMirror(file, {mode: this.getModeForMime(file.mime)})
      }
    }
  },
  methods: {
    getModeForMime (mime) {
      // remove charset info
      mime = mime.replace(/;.*$/, '')
      switch (mime) {
        case 'text/html':
          return 'htmlmixed'
        case 'text/x-markdown':
          return 'yaml-frontmatter'
      }
      return mime
    },
    changes: function (cm, changes) {
      this.value = this.mirror.getValue()
      this.$parent.dirty = true
    },
    save: function () {
      this.$parent.save()
    },
    setCodeMirror: function (file, options) {
      options = options || {}
      let p = document.querySelector('.text-editor')

      if (this.mirror) {
        // TODO: verify listener is removed
        this.mirror.off('changes', this.changes)
        if (p.firstChild) {
          p.removeChild(p.firstChild)
        }
      }
      if (file.document) {
        file.document.cm = null
      }
      this.mirror = CodeMirror(p, {
        value: file.document || file.content || '',
        mode: options.mode || 'htmlmixed',
        theme: options.theme || 'solarized dark',
        lineNumbers: true,
        keyMap: 'vim'
      })
      this.mirror.on('changes', this.changes)
      this.mirror.setOption('extraKeys', {
        'Ctrl-S': (cm) => {
          this.save()
        },
        // Useful for debugging, need to refresh a lot
        // and when the code editor is focused refresh
        // does not trigger without this
        'Ctrl-R': (cm) => {
          document.location.reload()
        },
        Tab: (cm) => {
          var spaces = Array(cm.getOption('indentUnit') + 1).join(' ')
          cm.replaceSelection(spaces)
        }
      })

      // window.focus(this.$el)
      // console.log('window focus: ')

      /*
      let wait
      let opts = {column: 80}
      let changing = false
      this.mirror.on('change', function (cm, change) {
        if (changing) return
        clearTimeout(wait)
        wait = setTimeout(function () {
          changing = true
          cm.wrapParagraphsInRange(change.from, CodeMirror.changeEnd(change), opts)
          changing = false
        }, 200)
      })
      */

      // This gives us :w in vim mode
      CodeMirror.commands.save = (cm) => {
        this.save()
      }

      if (!file.document) {
        file.document = this.mirror.getDoc()
        this.$parent.changeGeneration = file.document.changeGeneration()
      }

      this.$parent.codeMirror = this.mirror
    }
  },
  mounted: function () {
    let file = this.currentFile
    // Handles setting file content when switching tabs
    if (file && file.mime) {
      this.setCodeMirror(file, {mode: this.getModeForMime(file.mime)})
    }
  },
  components: {Hint}
}
</script>

<style>

  .source-editor > .hint {
    height: 4.2rem;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .hint + .text-editor {
    height: calc(100% - 4.2rem)
  }

  .text-editor {
    height: 100%;
  }

  .CodeMirror {
    font-size: 1.6rem;
    height: 100%;
    /* flex: 1 0 auto; */
  }

  .CodeMirror-dialog {
    font-size: 1.4rem;
    font-family: inherit;
  }

  .CodeMirror-dialog input {
    padding: 0.2rem;
    border-radius: 0;
  }
</style>
