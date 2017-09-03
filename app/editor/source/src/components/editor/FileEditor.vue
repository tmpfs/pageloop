<template>
  <div class="file-editor">
    <div class="scroll panel">
      <h2 class="file-info"><span v-bind:class="{hidden: !file.dir}">🗀</span><span v-bind:class="{hidden: file.dir}">🗎</span>&nbsp;{{file.name}}</h2>
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
      newName: ''
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
</style>
