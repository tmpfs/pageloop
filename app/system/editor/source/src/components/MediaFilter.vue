<template>
  <div class="media-filter">
    <div class="toggle">
      <a class="select" @click="toggle" title="Show filters">⏷</a>
    </div>
    <div class="filters" :class="{hidden: !show}">
      <a v-bind:class="{selected: filter === 'media'}"
        @click="select($event, 'media')"
        title="Show all media files">Media</a>
      <a v-bind:class="{selected: filter === 'images'}"
        @click="select($event, 'images')"
        title="Show image files">Images</a>
      <a v-bind:class="{selected: filter === 'text'}"
        @click="select($event, 'text')"
        title="Show text files">text</a>
      <a v-bind:class="{selected: filter === 'styles'}"
        @click="select($event, 'styles')"
        title="Show style files">Styles</a>
      <a v-bind:class="{selected: filter === 'scripts'}"
        @click="select($event, 'scripts')"
        title="Show script files">Scripts</a>
      <a v-bind:class="{selected: filter === 'audio'}"
        @click="select($event, 'audio')"
        title="Show audio files">audio</a>
      <a v-bind:class="{selected: filter === 'video'}"
        @click="select($event, 'video')"
        title="Show video files">video</a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'media-filter',
  data: function () {
    return {
      show: false,
      filter: 'media',
      label: 'media'
    }
  },
  methods: {
    toggle: function (e) {
      e.stopImmediatePropagation()
      this.show = !this.show
      const hide = () => {
        document.removeEventListener('click', hide)
        this.show = false
      }

      if (this.show) {
        document.addEventListener('click', hide)
      }
    },
    select: function (e, filter) {
      this.$parent.currentView = filter
      this.$parent.filter = filter
      this.filter = filter
      this.label = filter
    }
  }
}
</script>

<style scoped>

  .media-filter {
    position: relative;
  }

  .media-filter .toggle {
    display: inline-block;
    display: flex;
  }

  .media-filter .toggle > a:not(.select) {
    flex: 1 0;
  }

  .filters {
    position: absolute;
    top: 2.1rem;
    right: -2.2rem;
    display: flex;
    flex-direction: column;
    background: var(--base02-color);
    border: 1px solid var(--border-color);
    min-width: 10rem;
  }

  .select {
    display: inline-block;
    flex: none;
    background: var(--base02-color);
  }
</style>
