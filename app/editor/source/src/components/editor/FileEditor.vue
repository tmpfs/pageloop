<template>
  <div class="file-editor">
    <div class="scroll panel">
      <h2 class="file-info"><span v-bind:class="{hidden: !file.dir}">🗀</span><span v-bind:class="{hidden: file.dir}">🗎</span>&nbsp;{{file.name}}</h2>

      <section v-if="file.dir">
        <h3>Transfer Files</h3>
        <div
          class="uploader">
          <input type="file" multiple name="files" value="files" @change="change" />
          <div class="file-upload-input">
            <p v-if="!files.length">Tap to select files or drag files here to upload</p>
            <ul v-if="files.length">
              <li v-for="file in files">
                <span>{{file.name}}</span>
              </li>
            </ul>
          </div>
        </div>

        <div class="form-actions">
          <input
            @click="resetFiles"
            type="reset"
            name="Reset"
            value="Reset"
            :class="{hidden: !files.length}" />
          <input
            @click="uploadFiles"
            type="submit"
            name="Upload"
            value="Upload"
            class="primary"
            :class="{disabled: !files.length}" />
        </div>
      </section>
      <section>
        <h3>Rename File</h3>
        <p>Choose a new name for your file.</p>
        <form @submit="rename" class="rename">
          <input type="text" name="fileName" v-model="newName" />
          <div class="form-actions">
            <input type="submit" name="Rename" value="Rename" />
          </div>
        </form>
      </section>
      <section>
        <h3>Delete File</h3>
        <p v-bind:class="{hidden: confirmDelete}">Danger zone: be careful!</p>
        <div class="form-actions">
          <button @click="confirmDelete = true"
            v-bind:class="{hidden: confirmDelete}"
            class="danger">Delete {{file.url}}</button>
        </div>
        <div v-bind:class="{hidden: !confirmDelete}">
          <p>Are you sure you want to delete {{file.url}}?<br />
          <small>
            Deleting a file is irreversible, it cannot be undone.
          </small>
          </p>
          <nav class="form-actions">
            <button @click="confirmDelete = false">Cancel</button>
            <button @click="doDelete" class="danger">Delete</button>
          </nav>
        </div>
      </section>
      <section>
        <h3>File Info</h3>
        <ul class="small compact-list">
          <li>Name: {{file.name}}</li>
          <li>URL : {{file.uri}}</li>
          <li v-bind:class="{hidden: !file.dir}">Directory: yes</li>
          <li v-bind:class="{hidden: file.dir}">Size: {{file.size}} bytes</li>
          <li v-bind:class="{hidden: file.dir}">Mime: {{file.mime}}</li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  name: 'file-editor',
  data: function () {
    return {
      confirmDelete: false,
      newName: '',
      files: []
    }
  },
  computed: {
    file: function () {
      return this.$store.state.current
    }
  },
  created: function () {
    this.newName = this.file.url
  },
  methods: {
    rename: function (e) {
      e.preventDefault()
      return this.$store.dispatch('rename-file', {file: this.file, newName: this.newName})
        .catch((e) => console.error(e))
    },
    doDelete: function () {
      this.confirmDelete = false
      return this.$store.dispatch('delete-file', this.file)
    },
    syncHeight: function () {
      setTimeout(() => {
        let mask = this.$el.querySelector('.file-upload-input')
        let input = this.$el.querySelector('input[type="file"]')
        let b = mask.getBoundingClientRect()
        let h = b.bottom - b.top
        input.setAttribute('style', `height: ${h}px`)
      }, 50)
    },
    change: function (e) {
      this.files = e.target.files
      this.syncHeight()
    },
    uploadFiles: function (e) {
      e.preventDefault()
      const info = {files: this.files, dir: this.file.url}
      this.$parent.$parent.upload(info)
      return false
    },
    resetFiles: function (e) {
      e.preventDefault()
      this.files = []
      this.syncHeight()
    }
  },
  watch: {
    file: function (file) {
      this.newName = file.url
    }
  }
}
</script>

<style scoped>
  .file-editor {
    height: 100%;
  }

  .uploader {
    position: relative;
  }

  input[type="file"], .file-upload-input {
    position: relative;
    width: 100%;
    min-height: 8.2rem;
    opacity: 0;
  }

  .file-upload-input {
    position: absolute;
    opacity: 1;
    background: var(--base03-color);
    padding: 2rem;
    border-radius: 0.6rem;
    left: 0;
    top: 0;
    pointer-events: none;
    text-align: center;
  }

  .input[type="file"]:hover + .file-upload-input {
    color: var(--base3-color);
  }

  .file-upload-input * {
    pointer-events: none;
  }

  .uploader ul {
    font-size: 1.4rem;
    text-align: left;
    margin: 0;
    list-style-type: none;
    padding: 0;
  }
</style>
