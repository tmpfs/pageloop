<template>
  <header class="clearfix">
    <nav class="home">
      <a
        @click="$store.dispatch('navigate', {href: 'edit'})"
        :class="{selected: selectedView === 'edit', hidden: !this.$store.state.hasApplication()}"
        title="Edit Current Application">{{name}}</a>
    </nav>
    <nav class="main">
      <a
        @click="$store.dispatch('navigate', {href: 'home'})"
        class="home"
        :class="{selected: selectedView === 'home'}"
        title="Home page">Ꝏ </a>
      <a
        @click="$store.dispatch('navigate', {href: 'apps'})"
        :class="{selected: selectedView === 'apps'}"
        title="View and edit applications">Apps</a>
      <a
        @click="$store.dispatch('navigate', {href: 'docs'})"
        :class="{selected: selectedView === 'docs'}"
        title="Documentation">Docs</a>
      <a
        @click="$store.dispatch('navigate', {href: 'settings'})"
        :class="{selected: selectedView === 'settings'}"
        title="Settings">Settings</a>
    </nav>
  </header>
</template>

<script>
export default {
  name: 'app-header',
  computed: {
    name: function () {
      return this.$store.state.application
    },
    selectedView: function () {
      return this.$store.state.main.view
    }
  }
}
</script>

<style scoped>

  header {
    margin: 0;
    padding: 0;
    height: 3rem;
    border-bottom: 1px solid var(--border-color);
    user-select: none;
    cursor: default;
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

  nav.home > a.home {
    color: var(--base3-color) !important;
    text-decoration: none;
  }

  header > nav > a {
    font-size: 1.4rem;
    padding: .3rem 2rem;
    text-transform: uppercase;
    color: var(--base2-color) !important;
  }

  header > nav > a:hover, header > nav > a.selected {
    text-decoration: underline;
    color: var(--base3-color) !important;
    background: var(--blue-color);
  }

  header > nav > a.selected {
    background: var(--base03-color);
  }
</style>
