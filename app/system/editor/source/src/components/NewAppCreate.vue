<template>
  <div class="new-app-create small">
    <form @submit="createApplication">
      <label>Name: {{applicationName}}</label>
      <label>Description: {{applicationDescription}}</label>
      <div v-if="template">
        <label>Template: {{template.name}} ({{template.description}})</label>
      </div>
      <div class="form-actions">
        <input type="submit" value="Create Application" class="primary" />
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'new-app-create',
  data: function () {
    return {
      applicationId: this.$store.state.newApp.id,
      applicationName: this.$store.state.newApp.name,
      applicationDescription: this.$store.state.newApp.description,
      applicationTemplate: this.$store.state.newApp.templateUrl,
      template: this.$store.state.newApp.template
    }
  },
  methods: {
    createApplication: function (e) {
      e.preventDefault()

      let app = {}
      if (this.applicationId) {
        app.name = this.applicationId
      }
      if (this.applicationName) {
        app['display-name'] = this.applicationName
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
  }
}
</script>

<style scoped>
  .new-app-create {
    padding: 1rem;
  }

  label {
    display: block;
  }
</style>
