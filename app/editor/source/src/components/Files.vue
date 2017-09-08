<template>
  <div class="files-list">
    <p class="small" v-if="!list.length">No files found</p>
    <a
      @click="click($event, item)"
      @dragover="dragover"
      @dragleave="dragleave"
      :data-dir="item.dir ? item.url : ''"
      :data-url="item.url"
      class="file"
      :class="{selected: ~selection.indexOf(item)}"
      v-for="item in list">
      <span class="name">{{item.url}}</span>
    </a>
  </div>
</template>

<script>

import SelectableFileList from './SelectableFileList'

export default {
  name: 'files',
  mixins: [SelectableFileList],
  computed: {
    list: function () {
      return this.$store.state.app.files || []
    },
    selection: {
      get: function () {
        return this.$store.state.sidebar.files
      },
      set: function (val) {
        this.$store.state.sidebar.files = val
      }
    }
  },
  methods: {
    go: function (item) {
      return this.$store.dispatch('go-file', item)
    },
    getSelectionByUrl: function (url) {
      return this.$store.state.app.getFileByUrl(url)
    },
    dragover: function (e) {
      if (!e.currentTarget.getAttribute('data-dir')) {
        return
      }
      e.preventDefault()
      e.stopImmediatePropagation()
      e.currentTarget.classList.add('drop-target')
      return false
    },
    dragleave: function (e) {
      if (!e.currentTarget.getAttribute('data-dir')) {
        return
      }
      e.preventDefault()
      e.stopImmediatePropagation()
      e.currentTarget.classList.remove('drop-target')
      return false
    }
  }
}
</script>

<style scoped>
  .drop-target {
    border-top: 1px solid var(--base2-color);
  }

  .files-list p.small {
    margin-left: 2rem;
  }
</style>
