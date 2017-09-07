<template>
  <div
    tabindex="7"
    :class="{maximized: maximized === 'preview', minimized: maximized != '' && maximized !== 'preview'}"
    class="preview">
    <div class="column-header">
      <h2>Preview</h2>
        <nav class="toolbar clearfix">
          <a @click="refresh(file)"
             title="Refresh preview"
             :class="{hidden: path == ''}">Refresh</a>
          <a
            @click="maximized = 'preview'"
            :class="{hidden: maximized === 'preview'}"
            title="Maximize">◩</a>
          <a
            @click="maximized = ''"
            :class="{hidden: maximized !== 'preview'}"
            title="Minimize">▣</a>
        </nav>
    </div>
    <div class="column-options">
      <h3>{{path}}</h3>
    </div>
    <iframe :src="src" class="publish-preview"></iframe>
  </div>
</template>

<script>
export default {
  name: 'preview',
  data: function () {
    return {
      path: '',
      src: ''
    }
  },
  updated: function () {
    const frame = this.$el.querySelector('iframe')
    if (/\.pdf/.test(this.url)) {
      frame.removeAttribute('sandbox')
    } else {
      // frame.setAttribute('sandbox', this.sandbox)
    }
  },
  computed: {
    sandbox: function () {
      // return 'allow-same-origin allow-scripts'
    },
    maximized: {
      get: function () {
        return this.$store.state.editor.columns.maximized
      },
      set: function (val) {
        this.$store.commit('maximize-column', val)
      }
    },
    previewRefresh: function () {
      return this.$store.state.preview.refresh
    },
    url: {
      get: function () {
        return this.$store.state.preview.url
      },
      set: function (val) {
        return this.$store.commit('preview-url', val)
      }
    }
  },
  watch: {
    url: function (url) {
      this.refresh(url)
    },
    previewRefresh: function (val) {
      if (val === true) {
        this.refresh(this.file)
      }
      this.$store.commit('preview-refresh', false)
    }
  },
  mounted: function () {
    // This catches the case when switching main views
    // and a refresh is needed
    if (this.file) {
      this.refresh(this.file)
    }
  },
  methods: {
    refresh (file) {
      const url = file.uri

      // TODO: work out how to stop the iframe interpreting
      // TODO: the pdf as an HTML Document
      let allowed = /\.(html?|pdf|svg|jpe?g|png|gif)$/
      if (!file.dir && url && !allowed.test(url)) {
        return
      }
      if (url === '') {
        this.path = ''
        this.src = ''
        return
      }
      // If the src attribute will not change the page
      // won't be refreshed so we need to call reload()
      if (url === this.path) {
        let frame = document.querySelector('.publish-preview')
        return frame.contentDocument.location.reload()
      }
      this.file = file
      this.path = url
      this.src = this.getPreviewUrl(url, file)
    },
    getPreviewUrl (url, file) {
      if (url) {
        url = url.replace(/^\//, '')
      }

      let state = this.$store.state
      let host = state.host
      let base = state.app.url

      if (!host) {
        host = document.location.origin
      }

      return host + base + (url || '')
    }
  }
}
</script>

<style scoped>
  iframe.publish-preview {
    height: calc(100% - 4.6rem);
  }
</style>
