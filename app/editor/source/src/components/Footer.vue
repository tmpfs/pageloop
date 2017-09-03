<template>
  <footer>
    <!--<p class="log" v-bind:class="{error: error}">{{prefix}}{{message}}</p>-->
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
  footer {
    font-size: 1.6rem;
    line-height: 1.6rem;
    padding: 0.8rem 0;
    height: 3rem;
    /*min-height: 3.2rem;*/
    border-top: 1px solid var(--border-color);
  }

  footer p {
    margin: 0 2rem;
  }

  .log.error {
    color: var(--red-color);
  }
</style>
