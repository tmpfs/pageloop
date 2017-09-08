<template>
  <div class="pages-list">
    <p class="small" v-if="!list.length">No pages found</p>
    <a
      @click="click($event, item)"
      class="page"
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
  name: 'pages',
  mixins: [SelectableFileList],
  computed: {
    list: function () {
      return this.$store.state.app.pages || []
    },
    selection: {
      get: function () {
        return this.$store.state.sidebar.pages
      },
      set: function (val) {
        this.$store.state.sidebar.pages = val
      }
    }
  },
  methods: {
    go: function (item) {
      return this.$store.dispatch('go-page', item)
    },
    getSelectionByUrl: function (url) {
      return this.$store.state.app.getPageByUrl(url)
    }
  }
}
</script>

<style scoped>
  .pages-list p.small {
    margin-left: 2rem;
  }
</style>
