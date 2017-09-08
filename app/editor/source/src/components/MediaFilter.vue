<template>
  <div class="media-filter">
    <div class="toggle">
      <a v-bind:class="{selected: $parent.currentView === 'media'}"
        @click="toggle"
         title="Show media files">{{label}} <span>⏷</span></a>
    </div>
    <div class="filters" :class="{hidden: !show}">
      <a v-bind:class="{selected: filter === 'images'}"
        @click="select($event, 'images')"
        title="Show images">Images</a>
      <a v-bind:class="{selected: filter === 'styles'}"
        @click="select($event, 'styles')"
        title="Show styles">Styles</a>
      <a v-bind:class="{selected: filter === 'scripts'}"
        @click="select($event, 'scripts')"
        title="Show scripts">Scripts</a>
      <a v-bind:class="{selected: filter === 'media'}"
        @click="select($event, 'media')"
        title="Show all media files">Media</a>
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
      label: 'filter'
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
      this.$parent.currentView = 'media'
      this.filter = filter
      this.label = filter
      console.log('show media filter: ' + filter)
    }
  }
}
</script>

<style scoped>

  .media-filter {
    position: relative;
    flex: 1 0;
  }

  .media-filter .toggle {
    display: inline-block;
    display: flex;
  }

  .media-filter .toggle > a {
    flex: 1 0;
  }

  .filters {
    position: absolute;
    top: 2.2rem;
    left: 0;
    right: 0;
    display: flex;
    flex-direction: column;
    background: var(--base02-color);
    border: 1px solid var(--border-color);
  }
</style>
