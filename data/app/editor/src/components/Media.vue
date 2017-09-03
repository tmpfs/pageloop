<template>
  <div class="media-list">
  </div>
</template>

<script>
export default {
  name: 'media',
  computed: {
    list: function () {
      return this.$store.state.app.media
    }
  }
}
</script>

<style scoped>
</style>
