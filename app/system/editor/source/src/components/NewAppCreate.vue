<template>
  <div class="new-app-create small">
    <ul class="new-app-step small">
      <li>
        <span>Confirm app details</span>
        <span>Step 3/3</span>
      </li>
    </ul>
    <form @submit="createApplication">
      <ul class="details">
        <li>
          <span>Name</span>
          <span>{{applicationName}}</span>
        </li>
        <li>
          <span>Description</span>
          <span>{{applicationDescription}}</span>
        </li>
        <li>
          <span>Template</span>
          <span v-if="!template">None</span>
          <span v-else>{{template.name}}<br>{{template.description}}</span>
        </li>
      </ul>
      <div class="form-actions">
        <input type="submit" value="Create Application" class="primary" v-focus />
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
        app.display = this.applicationName
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
          this.$store.state.newApp.reset()
          return this.$store.dispatch('navigate', {href: `apps/user/${app.name}`})
        })
        .catch((e) => console.error(e))
    }
  }
}
</script>

<style scoped>
  .details {
    margin: 0 0 2rem 0;
    padding: 0;
    list-style-type: none;
    display: table;
    width: 100%;
  }

  .details > li {
    display: table-row;
    width: 100%;
  }

  .details > li > span:first-child {
    text-align: right;
  }


  .details > li > span:first-child::after {
    content: ':';
    display: inline-block;
  }

  .details > li > span {
    display: table-cell;
    width: 50%;
    padding: 0 1rem;
  }
</style>
