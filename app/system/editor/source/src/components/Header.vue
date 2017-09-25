<template>
  <header class="clearfix">
    <nav class="home">
      <a
        tabindex="1"
        @click="show = !show"
        class="home"
        title="Navigation Menu">Ꝏ  Menu</a>
    </nav>
  </header>
</template>

<script>
import BuildTasks from '@/components/BuildTasks'

export default {
  name: 'app-header',
  computed: {
    buildable: function () {
      const state = this.$store.state
      return state.hasApplication() && state.app.hasTasks() && state.main.view === 'edit'
    },
    name: function () {
      return this.$store.state.application
    },
    selectedView: function () {
      return this.$store.state.main.view
    },
    show: {
      get: function () {
        return this.$store.state.showMainMenu
      },
      set: function (val) {
        this.$store.state.showMainMenu = val
      }
    }
  },
  components: {BuildTasks}
}
</script>

<style>
  nav.home > a.home {
    color: var(--base3-color) !important;
    text-decoration: none;
  }

  header a {
    font-size: 1.4rem;
    padding: .3rem 2rem;
    text-transform: uppercase;
    color: var(--base2-color) !important;
  }

  header a:hover, header a.selected {
    /* text-decoration: underline; */
    color: var(--base3-color) !important;
    background: var(--blue-color);
  }

  header a.selected {
    background: var(--base03-color);
  }
</style>

<style scoped>

  header {
    margin: 0;
    padding: 0;
    height: 3rem;
    user-select: none;
    cursor: default;
    position: relative;
    z-index: 100;
  }

  header nav {
    display: flex;
    flex-direction: row;
    float: right;
    margin: 0;
    padding: 0;
  }

  header nav.home {
    float: left;
  }

</style>
