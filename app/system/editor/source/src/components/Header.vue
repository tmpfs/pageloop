<template>
  <header class="clearfix">
    <!--
    <nav class="home">
      <a
        tabindex="0"
        @click="$store.dispatch('navigate', {href: 'edit'})"
        title="Applications">
          <span>Apps</span>
          <span v-if="name"> / {{name}}</span>
        </a>
        <build-tasks v-if="buildable"></build-tasks>
    </nav>
    -->
    <nav class="main">
      <a
        tabindex="1"
        @click="$store.dispatch('navigate', {href: 'home'})"
        class="home"
        :class="{selected: selectedView === 'home'}"
        title="Home page">Ꝏ </a>
      <a
        tabindex="2"
        @click="$store.dispatch('navigate', {href: 'apps'})"
        :class="{selected: selectedView === 'apps'}"
        title="View and edit applications">Apps</a>
      <a
        tabindex="3"
        @click="$store.dispatch('navigate', {href: 'docs'})"
        :class="{selected: selectedView === 'docs'}"
        title="Documentation">Docs</a>
      <a
        tabindex="4"
        @click="$store.dispatch('navigate', {href: 'settings'})"
        :class="{selected: selectedView === 'settings'}"
        title="Settings">Settings</a>
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
