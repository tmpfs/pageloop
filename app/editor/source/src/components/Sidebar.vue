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
      @dragover="dragover"
      @dragend="dragleave"
      @dragleave="dragleave">
      <div class="uploads" :class="{hidden: !transfers.length}">
        <div class="upload" v-for="file in transfers">
          <div class="info">
            <span v-if="!file.complete">Uploading {{file.name}}</span>
            <span class="percent" v-if="!file.complete">{{Math.round(file.info.ratio * 100)}}%</span>
            <span class="complete" v-if="file.complete">Uploaded {{file.name}}</span>
          </div>
          <div class="progress" :class="{complete: file.complete}" :style="progress(file)"></div>
        </div>
      </div>
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
    return {}
  },
  computed: {
    transfers: function () {
      return this.$store.state.currentTransfer
    },
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
        return this.$store.state.sidebar.view
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
    progress: function (file) {
      return `width: ${Math.round(file.info.ratio * 100)}%`
    },
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
      // TODO: notify on error
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

      // Remove drop-target highlights
      const targets = document.querySelectorAll('.drop-target')
      targets.forEach((n) => {
        n.classList.remove('drop-target')
      })

      // We only accept file transfers
      if (!e.dataTransfer.files.length) {
        return false
      }

      this.$parent.transfer(e)

      return false
    },
    dragover: function (e) {
      e.preventDefault()
      e.currentTarget.classList.add('drop-target')
      return false
    },
    dragleave: function (e) {
      e.preventDefault()
      e.currentTarget.classList.remove('drop-target')
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
    border-top: 1px solid transparent;
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
    /* Need pointer events for file drag / drop */
    pointer-events: auto;
  }

</style>

<style scoped>
  .scroll {
    border-top: 1px solid transparent;
    transition: all 0.3s ease-out;
  }

  .scroll.drop-target {
    border-top: 1px solid var(--base2-color);
  }

  .uploads {
    border-bottom: 1px solid var(--border);
    font-size: 1.4rem;
    margin-bottom: 1rem;
  }

  .upload {
    padding: 0 1rem;
    height: 2.8rem;
    overflow: hidden;
  }

  .upload > .info {
    display: flex;
    justify-content: flex-end;
  }

  .upload > .info:first-child {
    align-self: flex-start;
  }

  .upload .complete {
    color: var(--green-color);
  }

  .upload > .progress {
    height: 0.6rem;
    width: 100%;
    background: var(--cyan-color);
  }

  .upload > .progress.complete {
    background: var(--green-color);
  }

  .percent {
    margin-left: auto;
  }


</style>
