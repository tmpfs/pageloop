<template>
  <div class="content-main not-found">
    <h2>Not Found</h2>
    <p>Oops, we could not find a page that matched your request <code>{{flash}}</code></p>
  </div>
</template>

<script>
export default {
  name: 'not-found',
  computed: {
    flash: function () {
      let flash = this.$store.state.flash
      if (flash) {
        flash = `#${flash}`
      }
      return flash
    }
  }
}
</script>

<style scoped>

</style>
