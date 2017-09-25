<template>
  <div class="main-view">
    <main-menu></main-menu>
    <component v-bind:is="currentView"></component>
  </div>
</template>

<script>

import MainMenu from '@/components/MainMenu'
import NotFound from '@/components/NotFound'
import Home from '@/components/Home'
import Apps from '@/components/Apps'
import Edit from '@/components/Edit'
import Docs from '@/components/Docs'
import Settings from '@/components/Settings'

export default {
  name: 'app-main',
  computed: {
    currentView: function () {
      return this.$store.state.main.view
    }
  },
  components: {MainMenu, NotFound, Home, Apps, Edit, Docs, Settings}
}
</script>

<style scoped>
  .main-view {
    display: flex;
    flex-direction: row;
    height: 100%;
  }
</style>
