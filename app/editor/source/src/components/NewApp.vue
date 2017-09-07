<template>
  <div class="new-app">
    <h2>New Application</h2>
    <form @submit="chooseTemplate">
      <p class="small">Choose an application name:</p>
      <input type="text" name="name"
        :value="applicationName" v-model="applicationName" />
      <p class="small">Short description:</p>
      <input type="text" name="description"
        :value="applicationDescription" v-model="applicationDescription" />
      <div class="form-actions">
        <input type="submit" value="Next: Choose a template" class="primary" />
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'new-app',
  data: function () {
    return {
      applicationName: 'new-app',
      applicationDescription: 'New application'
    }
  },
  methods: {
    chooseTemplate: function (e) {
      e.preventDefault()
      console.log('choose template')
    },
    createApplication: function (e) {
      e.preventDefault()
      let app = {}
      if (this.applicationName) {
        app.name = this.applicationName
      }
      if (this.applicationDescription) {
        app.description = this.applicationDescription
      }
      this.$store.dispatch('new-app', app)
        .then(() => {
          return this.$store.dispatch('navigate', {href: `apps/user/${app.name}`})
        })
        .catch((e) => console.error(e))
    }
  }
}
</script>

<style scoped>
  .new-app {
    padding: 2rem;
    width: 20%;
    border-right: 1px solid var(--border-color)
  }
</style>
