<template>
  <div class="visual-editor">
    <component v-bind:is="currentView"></component>
  </div>
</template>

<script>

import BitmapEditor from '@/components/editor/BitmapEditor'

export default {
  name: 'visual-editor',
  computed: {
    currentView: function () {
      // console.log(this.currentFile)

      const file = this.currentFile

      if (/^image\/(gif|jpeg|png)/.test(file.mime)) {
        return 'bitmap-editor'
      }
    },
    currentFile: function () {
      return this.$store.state.current
    }
  },
  components: {BitmapEditor}
}
</script>

<style scoped>
</style>
