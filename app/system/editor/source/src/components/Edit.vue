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
    upload: function (info) {
      const state = this.$store.state

      // Set up file transfer data
      this.$store.commit('transfers', info)

      const uploader = () => {
        // Start the file upload
        this.$store.dispatch('upload', info)
          .then((transfers) => {
            let names = transfers.filter((f) => {
              return f.handle
            }).map((f) => {
              return f.handle.url
            })
            names = names.join(' ')
            state.notify({title: `Transfer (${info.files.length})`, message: `Uploaded ${names}`})
          })
          .catch((e) => {
            console.error(e)
            state.notify({error: e})
          })
      }

      // Check for existing files which need POST
      const existing = state.transfer.transfers.filter((f) => {
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
            return uploader()
          }
        }

        existing.forEach((f) => {
          details.note += f.url + '\n'
        })

        return this.$store.commit('alert-show', details)
      }

      // All new files - upload them
      uploader()
    },
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
      return this.upload(info)
    }
  }
}
</script>

<style scoped>
  /* Column focus */
  > :focus {
    border-top: 1px solid var(--blue-color);
  }

  .content-main {
    border-top: 0;
  }
</style>
