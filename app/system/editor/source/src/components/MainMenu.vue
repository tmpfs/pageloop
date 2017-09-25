<template>
  <div class="main-menu" v-if="show">
    <nav>
      <a
        tabindex="1"
        @click="$store.dispatch('navigate', {href: 'home'})"
        class="home"
        :class="{selected: selectedView === 'home'}"
        title="Home page">Home</a>
      <a
        tabindex="2"
        @click="$store.dispatch('navigate', {href: 'apps'})"
        :class="{selected: selectedView === 'apps'}"
        title="View and edit applications">Applications</a>
      <a
        tabindex="3"
        @click="$store.dispatch('navigate', {href: 'docs'})"
        :class="{selected: selectedView === 'docs'}"
        title="Documentation">Documentation</a>
      <a
        tabindex="4"
        @click="$store.dispatch('navigate', {href: 'settings'})"
        :class="{selected: selectedView === 'settings'}"
        title="Settings">Settings</a>
    </nav>
  </div>
</template>

<script>
export default {
  name: 'main-menu',
  computed: {
    selectedView: function () {
      return this.$store.state.main.view
    },
    show: function () {
      return this.$store.state.showMainMenu
    }
  }
}
</script>

<style scoped>
  .main-menu {
    /*
    position: absolute;
    left: 0;
    top: 2.8rem;
    bottom: 0;
    right: 0;
    z-index: 1000;
    */
    width: 16%;
    border-top: 1px solid var(--border-color);
    border-right: 1px solid var(--border-color);
    font-size: 1.5rem;
    background: var(--base03-color);
  }

  .main-menu > nav {
    padding: 1rem;
  }

  .main-menu > nav > a {
    display: block;
    transition: all 0.4s ease-out;
    color: var(--base00-color);
    margin-bottom: 0.5rem;
    border-bottom: 1px solid transparent;
  }

  .main-menu > nav > a:hover {
    color: var(--base2-color);
  }

  .main-menu > nav > a.selected {
    color: var(--base2-color);
    border-bottom: 1px solid var(--base2-color);
  }
</style>
