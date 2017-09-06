<template>
  <footer>
    <!--<p class="log" v-bind:class="{error: error}">{{prefix}}{{message}}</p>-->
    <nav>
      <!--<a></a>-->
      <a
        class="reset"
        :class="{hidden: !canReset, disabled: !needsReset}"
        @click="reset"
        title="Reset columns"></a>

    </nav>
  </footer>
</template>

<script>
export default {
  name: 'app-footer',
  computed: {
    canReset: function () {
      let state = this.$store.state
      return state.main.view === 'edit' && state.hasApplication()
    },
    needsReset: function () {
      let state = this.$store.state
      return state.columns.custom || state.columns.maximized
    },
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
  },
  methods: {
    reset (e) {
      e.preventDefault()
      this.$store.commit('reset-column')
    }
  }
}
</script>

<style scoped>
  footer {
    font-size: 1.6rem;
    line-height: 1.6rem;
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

  nav {
    text-align: right;
    user-select: none;
  }

  nav a:hover {
    background: var(--base03-color);
    color: var(--base3-color);
  }

  nav a {
    display: inline-block;
    width: 3rem;
    height: 3rem;
    line-height: 3rem;
    color: currentColor;
    text-align: center;
  }

  a.reset {
    transform: rotate(90deg);
    font-size: 2.2rem;
  }
</style>
