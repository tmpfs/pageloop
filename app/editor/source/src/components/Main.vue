<template>
  <component v-bind:is="currentView"></component>
</template>

<script>

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
  components: {NotFound, Home, Apps, Edit, Docs, Settings}
}
</script>

<style scoped>

</style>
