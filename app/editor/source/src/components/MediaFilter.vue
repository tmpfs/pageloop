<template>
  <div class="media-filter">
    <a v-bind:class="{selected: $parent.currentView === 'media'}"
      @click="$parent.currentView = 'media'"
      title="Show media files">Media</a>
  </div>
</template>

<script>
export default {
  name: 'media-filter',
  props: ['parentView'],
  computed: {
    currentView: {
      get: function () {
        return this.parentView
      },
      set: function (val) {
        this.parentView = val
      }
    }
  }
}
</script>

<style scoped>
  .media-filter {
    display: inline-block;
    display: flex;
  }

  .media-filter > a {
    flex: 1 0;
  }
</style>
