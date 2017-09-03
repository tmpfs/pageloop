<template>
  <div class="files-list">
    <a @click="click(item)" class="file" :class="{selected: currentFile.url === item.url}" v-for="item in list">
      <span class="name">{{item.url}}</span>
    </a>
  </div>
</template>

<script>
export default {
  name: 'files',
  computed: {
    currentFile: function () {
      return this.$store.state.app.current
    },
    list: function () {
      return this.$store.state.app.files
    }
  },
  methods: {
    click: function (item) {
      return this.$store.dispatch('go-file', item)
    }
  }
}
</script>

<style scoped>
</style>
