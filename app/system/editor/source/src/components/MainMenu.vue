<template>
  <div class="main-menu">
    <div class="column-header">
      <h2>Navigation</h2>
    </div>
    <nav>
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
    }
  }
}
</script>

<style scoped>
  .main-menu {
    width: 18rem;
    border-right: 1px solid var(--border-color);
    font-size: 1.5rem;
  }

  .main-menu > nav {
    padding: 1rem;
  }

  .main-menu > nav > a {
    display: block;
    transition: all 0.4s ease-out;
    color: var(--base00-color);
    margin-bottom: 0.5rem;
  }

  .main-menu > nav > a:hover {
    color: var(--base2-color);
  }

  .main-menu > nav > a.selected {
    color: var(--base2-color);
    border-bottom: 1px solid var(--base2-color);
  }

</style>
