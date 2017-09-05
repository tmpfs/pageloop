<template>
  <div class="content-main">
    <div class="content">
      <sidebar></sidebar>
      <editor></editor>
      <preview></preview>
    </div>
  </div>
</template>

<script>

import Sidebar from '@/components/Sidebar'
import Editor from '@/components/Editor'
import Preview from '@/components/Preview'

export default {
  name: 'edit',
  components: {Sidebar, Editor, Preview},
  methods: {
    transfer: function (e) {
      // Check if drop occured on a directory
      let dir
      if (e.target) {
        dir = e.target.getAttribute('data-dir')
        // TOOD: handle deeply nested children in files list
        if (!dir && e.target.parentElement && e.target.parentElement.getAttribute('data-dir')) {
          dir = e.target.parentElement.getAttribute('data-dir')
        }
      }

      const info = {files: e.dataTransfer.files, dir: dir}
      const state = this.$store.state

      // Set up file transfer data
      this.$store.commit('transfers', info)

      const upload = () => {
        // Start the file upload
        this.$store.dispatch('upload', info)
          .then((transfers) => {
            let names = transfers.map((f) => {
              return f.handle.url
            }).join(' ')
            state.notify({title: `Transfer Complete (${info.files.length})`, message: `Uploaded ${names}`})
          })
          .catch((e) => {
            state.notify({error: e})
          })
      }

      // Check for existing files which need POST
      const existing = state.transfers.filter((f) => {
        return f.exists
      }).map((f) => {
        return f.exists
      })

      if (existing.length) {
        let details = {
          title: `Overwrite`,
          message: `Are you sure you want to overwrite files on upload?`,
          note: '',
          ok: () => {
            return upload()
          }
        }

        existing.forEach((f) => {
          details.note += f.url + '\n'
        })

        return this.$store.commit('alert-show', details)
      }

      // All new files - upload them
      upload()
    }
  }
}
</script>

<style scoped>
</style>
