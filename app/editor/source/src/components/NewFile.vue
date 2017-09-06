<template>
  <div class="new-file">
    <section>
      <h3>File Name</h3>
      <form @submit="createNewFile">
        <input v-model="fileName" type="text" name="name" :value="fileName" />
        <p class="hint">

          Use <code>/path/to/file/document.md</code> to create directories when adding new files.

          To create a directory end the name with a slash (<code>/path/to/dir/</code>).
          <a class="close minimal" title="Dismiss this hint"></a>
        </p>

        <div class="template-select">
          <h3>Template</h3>
          <p class="small">Select an optional file template:</p>
          <ul class="small compact-list">
            <li>
              <input type="radio" v-model="template"
                id="empty-file" name="template" value="" checked />
              <label for="empty-file">Empty File</label>
            </li>
            <li>
              <input type="radio" @change="extension = '.md'" v-model="template"
                id="markdown-partial" name="template" value="markdown-partial" />
              <label for="markdown-partial">Markdown Partial</label>
            </li>
            <li>
              <input type="radio" @change="extension = '.md'" v-model="template"
                id="markdown-standalone" name="template" value="markdown-standalone" />
              <label for="markdown-standalone">Markdown Standalone</label>
            </li>
            <li>
              <input type="radio" v-model="template"
                id="html-layout" @change="extension = '.html'" name="template" value="html-layout" />
              <label for="html-layout">HTML Layout</label>
            </li>
            <li>
              <input type="radio" v-model="template"
                id="html-partial" @change ="extension = '.html'" name="template" value="html-partial" />
              <label for="html-partial">HTML Partial</label>
            </li>
            <li>
              <input type="radio" v-model="template"
                id="html-standalone" @change="extension = '.html'" name="template" value="html-standalone" />
              <label for="html-standalone">HTML Standalone</label>
            </li>
          </ul>
        </div>
        <nav class="form-actions">
          <input @click="cancel" type="reset" name="Reset" value="Cancel" />
          <input type="submit" name="Create" value="Create" class="primary" />
        </nav>
      </form>
    </section>
  </div>
</template>

<script>
export default {
  name: 'new-file',
  data: function () {
    return {
      templateMap: {
        'markdown-partial': {
          container: 'template',
          application: 'documents',
          file: '/markdown-partial.md'
        },
        'markdown-standalone': {
          container: 'template',
          application: 'documents',
          file: '/markdown-standalone.md'
        },
        'html-layout': {
          container: 'template',
          application: 'documents',
          file: '/layout.html'
        },
        'html-partial': {
          container: 'template',
          application: 'documents',
          file: '/html-partial.html'
        },
        'html-standalone': {
          container: 'template',
          application: 'documents',
          file: '/html-standalone.html'
        }
      },
      fileName: '/untitled.md',
      template: '',
      extension: ''
    }
  },
  watch: {
    extension: function (val) {
      this.displayExtension = val
    }
  },
  computed: {
    displayExtension: {
      get: function () {
        return this.extension
      },
      set: function (val) {
        if (val) {
          let current = this.fileName
          if (/[^.]+\.[^.]*$/.test(current)) {
            current = current.replace(/\.[^.]*$/, val)
            this.fileName = current
          }
        }
      }
    }
  },
  methods: {
    cancel: function (e) {
      e.preventDefault()
      this.$parent.closeNewFileView()
    },
    createNewFile: function (e) {
      e.preventDefault()
      let action = 'go-file'
      if (this.$parent.previousView === 'pages') {
        action = 'go-page'
      }
      return this.$store.dispatch(
        'new-file', {name: this.fileName, template: this.templateMap[this.template], action: action})
        .then(() => {
          this.$parent.closeNewFileView()
        })
        .catch((e) => console.error(e))
    }
  }
}
</script>

<style scoped>
  .new-file {
    padding: 1rem;
  }

  h3 {
    display: inline-block;
    font-size: 1.4rem;
    border-bottom: 1px solid currentColor;
    margin: 0 0 1rem 0;
  }

</style>
