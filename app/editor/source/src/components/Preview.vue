<template>
  <div
    tabindex="-1"
    @click="focus"
    :class="{maximized: maximized === 'preview', minimized: maximized != '' && maximized !== 'preview'}"
    class="preview">
    <div class="column-header">
      <h2>Preview</h2>
        <nav class="toolbar clearfix">
          <a @click="refresh(path)"
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
    <iframe :src="src" sandbox="allow-same-origin allow-scripts" class="publish-preview"></iframe>
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
  computed: {
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
        this.refresh(this.path)
      }
      this.$store.commit('preview-refresh', false)
    }
  },
  mounted: function () {
    // This catches the case when switching main views
    // and a refresh is needed
    if (this.url) {
      this.refresh(this.url)
    }
  },
  methods: {
    refresh (url) {
      // TODO: work out how to stop the iframe interpreting
      // TODO: the pdf as an HTML Document
      let allowed = /\.(html?|pdf|svg|jpe?g|png|gif)$/
      if (url && !allowed.test(url)) {
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
      this.path = url
      this.src = this.getPreviewUrl(url)
    },
    getPreviewUrl (url) {
      if (url) {
        url = url.replace(/^\//, '')
      }

      let state = this.$store.state
      let host = state.host

      if (!host) {
        host = document.location.origin
      }

      return host + this.$store.state.app.url + (url || '')
    },
    focus: function () {
      this.$el.focus()
    }
  }
}
</script>

<style scoped>
  iframe.publish-preview {
    height: calc(100% - 4.6rem);
  }
</style>
