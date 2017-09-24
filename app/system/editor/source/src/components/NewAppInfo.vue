<template>
  <div class="new-app-info">
    <form @submit="nextStep">
      <label class="small">Name:</label>
      <input type="text" name="name" placeholder="Enter an app name"
        :value="applicationName" v-model="applicationName" />
      <label class="small">Description:</label>
      <input type="text" name="description" placeholder="Enter an app description"
        :value="applicationDescription" v-model="applicationDescription" />
      <div class="form-actions">
        <input type="submit" value="Next: Select a template" class="primary" />
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'new-app-info',
  data: function () {
    return {
      applicationName: '',
      applicationDescription: ''
    }
  },
  methods: {
    nextStep: function (e) {

    }
    /*
    createApplication: function (e) {
      e.preventDefault()

      let app = {}
      if (this.applicationName) {
        app.name = this.applicationName
      }
      if (this.applicationDescription) {
        app.description = this.applicationDescription
      }
      if (this.template) {
        app.template = {
          container: this.template.container,
          application: this.template.name
        }
      }
      this.$store.dispatch('new-app', app)
        .then(() => {
          return this.$store.dispatch('navigate', {href: `apps/user/${app.name}`})
        })
        .catch((e) => console.error(e))
    }
    */
  }
}
</script>

<style scoped>
  .new-app-info {
    padding: 1rem;
  }
</style>
