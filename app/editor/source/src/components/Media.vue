<template>
  <div class="media-list">
    <p class="small" v-if="!list.length">No media files found</p>
    <a
      @click="click($event, item)"
      class="file"
      :data-url="item.url"
      :class="{selected: ~selection.indexOf(item)}"
      v-for="item in list">
      <span class="name">{{item.url}}</span>
    </a>
  </div>
</template>

<script>
import SelectableFileList from './SelectableFileList'

export default {
  name: 'media',
  mixins: [SelectableFileList],
  props: ['filter'],
  computed: {
    list: function () {
      return this.$store.state.app[this.filter] || []
    },
    selection: {
      get: function () {
        return this.$store.state.sidebar.media
      },
      set: function (val) {
        this.$store.state.sidebar.media = val
      }
    }
  },
  methods: {
    go: function (item) {
      return this.$store.dispatch('go-file', item)
    },
    getSelectionByUrl: function (url) {
      return this.$store.state.app.getFileByUrl(url)
    }
  }
}
</script>

<style scoped>
  .media-list p.small {
    margin-left: 2rem;
  }
</style>
