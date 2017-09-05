<template>
  <div class="files-list">
    <a
      @click="click(item)"
      @dragover="dragover"
      @dragleave="dragleave"
      :data-dir="item.dir ? item.url : ''"
      class="file"
      :class="{selected: currentFile.url === item.url}"
      v-for="item in list">
      <span class="name">{{item.url}}</span>
    </a>
  </div>
</template>

<script>
export default {
  name: 'files',
  computed: {
    currentFile: function () {
      return this.$store.state.app.current
    },
    list: function () {
      return this.$store.state.app.files
    }
  },
  methods: {
    click: function (item) {
      return this.$store.dispatch('go-file', item)
    },
    dragover: function (e) {
      if (!e.currentTarget.getAttribute('data-dir')) {
        return
      }
      e.preventDefault()
      e.stopImmediatePropagation()
      if (e.dataTransfer.files.length) {
        e.currentTarget.classList.add('drop-target')
      } else {
        e.currentTarget.classList.add('drop-disabled')
      }
      return false
    },
    dragleave: function (e) {
      if (!e.currentTarget.getAttribute('data-dir')) {
        return
      }
      e.preventDefault()
      e.stopImmediatePropagation()
      e.currentTarget.classList.remove('drop-target')
      e.currentTarget.classList.remove('drop-disabled')
      return false
    }
  }
}
</script>

<style scoped>
  .drop-target {
    border-top: 1px solid var(--base2-color);
  }

  .drop-disabled {
    cursor: no-drop;
  }
</style>
