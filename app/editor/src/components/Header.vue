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
      return this.$store.state.mainView
    }
  }
}
</script>

<style scoped>

</style>
