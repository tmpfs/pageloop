<template>
  <div
    tabindex="7"
    :class="{maximized: maximized === 'preview', minimized: maximized != '' && maximized !== 'preview'}"
    class="preview">
    <div class="column-header">
      <h2>Preview</h2>
        <nav class="toolbar clearfix">
          <a @click="reload()"
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
    <iframe @load="loaded" :src="src" class="publish-preview"></iframe>
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
    // TODO: figure out sandbox that allows pdfs to render
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
    file: {
      get: function () {
        return this.$store.state.preview.file
      },
      set: function (val) {
        return this.$store.commit('preview-change', val)
      }
    },
    url: function () {
      return this.file ? this.file.uri : undefined
    }
  },
  watch: {
    url: function (val) {
      this.refresh(this.file)
    },
    // We need this awkward toggle to ensure the property
    // actually changes each time
    previewRefresh: function (val) {
      if (val === true) {
        this.reload()
      }
      this.$store.commit('preview-refresh', false)
    }
  },
  mounted: function () {
    // This catches the case when switching main views
    // and a refresh is needed
    /*
    if (this.file) {
      this.refresh(this.file)
    }
    */
  },
  methods: {
    loaded: function (e) {
      const app = this.$store.state.app
      const base = app.publicUrl
      const win = e.currentTarget.contentWindow
      let pathname = win.location.pathname
      pathname = pathname.replace(base, '')
      this.path = pathname + win.location.hash

      if (/\.(txt)$/.test(pathname)) {
        console.log('show text file...')
      }

      // Kludge so we can show the anchor hash in the preview path
      win.addEventListener('click', (e) => {
        const win = e.view
        setTimeout(() => {
          this.path = win.location.pathname.replace(base, '') + win.location.hash
        }, 50)
      })
    },
    reload: function () {
      // If the src attribute will not change the page
      // won't be refreshed so we need to call reload()
      let frame = document.querySelector('.publish-preview')
      return frame.contentDocument.location.reload()
    },
    refresh (file) {
      if (!file) {
        this.path = ''
        this.src = ''
        return
      }

      const url = file.uri

      // console.log('binary: ' + file.binary)

      let allowed = /\.(html?|pdf|svg|jpe?g|png|gif|txt)$/
      if (!file.dir && url && !allowed.test(url)) {
        return
      }

      this.src = this.getPreviewUrl(url)
      console.log('preview src: ' + url)
    },
    getPreviewUrl (url) {
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
