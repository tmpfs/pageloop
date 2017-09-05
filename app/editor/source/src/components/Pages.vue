<template>
  <div class="pages-list">
    <a
      @click="click(item)"
      class="page"
      :class="{selected: currentFile.url === item.url, uploader: item.dir}"
      v-for="item in list">
      <span class="name">{{item.url}}</span>
    </a>
  </div>
</template>

<script>
export default {
  name: 'pages',
  computed: {
    currentFile: function () {
      return this.$store.state.app.current
    },
    list: function () {
      return this.$store.state.app.pages
    }
  },
  methods: {
    click: function (item) {
      return this.$store.dispatch('go-page', item)
    }
  }
}
</script>

<style scoped>
</style>
