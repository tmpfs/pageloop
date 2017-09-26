<template>
  <header>
    <nav>
      <!--
      <a
        tabindex="1"
        @click="show = !show"
        :class="{selected: show}"
        class="menu"
        title="Navigation Menu"><span>Ꝏ</span>Menu</a>
      -->
      <a
        tabindex="1"
        @click="show = !show"
        :class="{selected: show}"
        class="menu"
        title="Navigation Menu">
        <i class="fa fa-navicon"></i>
        Menu</a>
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

  header a {
    display: inline-block;
    height: 3rem;
  }

  header a {
    font-size: 1.4rem;
    line-height: 3rem;
    vertical-align: middle;
    padding-left: 1rem;
    text-transform: uppercase;
    color: var(--base2-color) !important;
  }

  header a:hover {
    color: var(--base3-color) !important;
  }

  header a:hover, header a.selected {
    background: var(--base03-color);
  }

  header a.selected {
    background: var(--base03-color);
    pointer-events: auto;
  }

  a.menu {
    width: 16%;
  }

  a.menu i {
    margin-right: 0.5rem;
  }

</style>
