<template>
  <div class="sidebar"
    :class="{maximized: maximized === 'sidebar', minimized: maximized != '' && maximized !== 'sidebar'}">
    <div class="column-header">
      <nav class="toolbar">
        <a
          @click="confirmDelete"
          v-bind:class="{hidden: !canDelete}"
          title="Delete File">➖</a>
        <a
          @click="showNewFileView"
          v-bind:class="{disabled: currentView === 'new-file'}"
          title="New File">➕</a>
        <a
          @click="maximized = 'sidebar'"
          :class="{hidden: maximized === 'sidebar'}"
          title="Maximize">◩</a>
        <a
          @click="maximized = ''"
          :class="{hidden: maximized !== 'sidebar'}"
          title="Minimize">▣</a>
      </nav>
    </div>
    <div class="column-options">
      <nav class="tabs">
        <a v-bind:class="{selected: currentView === 'pages'}"
          @click="currentView = 'pages'"
          title="Show pages">Pages</a>
        <a v-bind:class="{selected: currentView === 'files'}"
          @click="currentView = 'files'"
          title="Show files">Files</a>
        <a v-bind:class="{selected: currentView === 'media'}"
          @click="currentView = 'media'"
          title="Show media files">Media</a>
      </nav>
    </div>
    <div
      class="scroll"
      @drop="drop"
      @dragover="noop"
      @dragend="noop"
      @dragleave="noop">
      <component v-bind:is="currentView"></component>
    </div>
    <div class="column-drag" :class="{hidden: maximized}" @mousedown="resizeColumn"></div>
  </div>
</template>

<script>

import NewFile from '@/components/NewFile'
import Pages from '@/components/Pages'
import Files from '@/components/Files'
import Media from '@/components/Media'

export default {
  name: 'sidebar',
  data: function () {
    return {
      shouldDelete: false
    }
  },
  computed: {
    maximized: {
      get: function () {
        return this.$store.state.columns.maximized
      },
      set: function (val) {
        this.$store.commit('maximize-column', val)
      }
    },
    canDelete: function () {
      return this.$store.state.hasFile()
    },
    currentFile: function () {
      return this.$store.state.current
    },
    currentView: {
      get: function () {
        return this.$store.state.sidebarView
      },
      set: function (val) {
        var values = [val]
        var file = this.$store.state.current
        if (file !== null) {
          if (val === 'files') {
            values.push(file.url)
          } else if (val === 'pages' && this.$store.state.isPage(file)) {
            values.push(file.url)
          }
        }
        let href = this.$store.state.getAppHref(...values)
        this.$store.dispatch('navigate', {href: href})
      }
    }
  },
  methods: {
    confirmDelete: function () {
      let details = {
        title: `Delete File (${this.currentFile.name})`,
        message: `Are you sure you want to delete the file ${this.currentFile.url}?`,
        note: 'Be careful file deletion is irreversible.',
        ok: () => {
          this.doDeleteFile()
        }
      }
      this.$store.commit('alert-show', details)
    },
    doDeleteFile: function () {
      return this.$store.dispatch('delete-file', this.currentFile)
        .catch((e) => console.error(e))
    },
    showNewFileView: function () {
      this.previousView = this.currentView
      this.currentView = 'new-file'
    },
    closeNewFileView: function () {
      this.currentView = this.previousView || 'pages'
    },
    resizeColumn: function (e) {
      this.$store.dispatch('resize-column', e)
    },
    drop: function (e) {
      e.preventDefault()
      const opts = {files: e.dataTransfer.files}
      this.$store.dispatch('upload', opts)
      return false
    },
    noop: function (e) {
      e.preventDefault()
      // console.log('noop called')
      return false
    }
  },
  components: {NewFile, Pages, Files, Media}
}
</script>

<style>
  .page, .file {
    display: block;
    padding: .2rem 0 .2rem 1rem;
    cursor: pointer;
    user-select: none;
    background: transparent;
    transition: all 0.3s ease-out;
  }

  .page, .file {
    font-size: 1.6rem;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .files-list a, .pages-list a {
    color: currentColor;
  }

  .files-list a:hover, .pages-list a:hover {
    color: var(--base2-color);
  }

  .file.selected, .page.selected {
    background: var(--base03-color);
    color: var(--base3-color);
  }
</style>
