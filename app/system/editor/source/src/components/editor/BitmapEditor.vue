<template>
  <div class="bitmap-editor">
    <img :src="src" />
  </div>
</template>

<script>

export default {
  name: 'bitmap-editor',
  computed: {
    src: function () {
      return this.$store.state.app.url + this.$store.state.current.uri
    }
  }
}
</script>

<style scoped>
</style>
