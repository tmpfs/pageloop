<template>
  <div class="content-main not-found">
    <h2>Not Found</h2>
    <p>Oops, we could not find a page that matched your request <code v-if="flash">{{flash}}</code></p>
  </div>
</template>

<script>
export default {
  name: 'not-found',
  computed: {
    flash: function () {
      let flash = this.$store.state.flash.message
      if (flash) {
        flash = `#${flash}`
      }
      return flash
    }
  }
}
</script>

<style scoped>
  .not-found {
    padding: 4rem;
  }

  .not-found h2 {
    display: block;
    padding: 0;
  }
</style>
