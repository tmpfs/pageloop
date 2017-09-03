<template>
  <footer>
    <p class="log" v-bind:class="{error: error}">{{prefix}}{{message}}</p>
  </footer>
</template>

<script>
export default {
  name: 'app-footer',
  computed: {
    message: function () {
      return this.$store.state.log.toString()
    },
    error: function () {
      return (this.$store.state.log.last instanceof Error)
    },
    prefix: function () {
      let msg = this.$store.state.log.last
      if (msg && this.error) {
        return '! '
      } else if (this.message && !this.error) {
        return '# '
      }
      return ''
    }
  }
}
</script>

<style scoped>

</style>
