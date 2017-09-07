<template>
  <div class="new-app">
    <div class="new-app-fields">
      <h2>New Application</h2>
      <form @submit="createApplication">
        <label class="small">Application name:</label>
        <input type="text" name="name"
          :value="applicationName" v-model="applicationName" />
        <label class="small">Short description:</label>
        <input type="text" name="description"
          :value="applicationDescription" v-model="applicationDescription" />
        <div class="form-actions">
          <input type="submit" value="Create Application" class="primary" />
        </div>
      </form>
    </div>
    <div class="new-app-templates">
      <h2>Choose Template</h2>
      <div class="app-template" v-for="tpl in templates">
        <input :id="tpl.name" type="radio" name="template"></input>
        <label :for="tpl.name">{{tpl.name}}</label>
        <p class="small">{{tpl.url}}</p>
        <iframe :src="tpl.url"></iframe>
      </div>
    </div>
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
  computed: {
    templates: function () {
      return this.$store.state.templates
    }
  },
  mounted: function () {
    this.$store.dispatch('list-templates')
  },
  methods: {
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
    display: flex;
    padding: 2rem;
    width: 100%;
    border-right: 1px solid var(--border-color)
  }

  .new-app-fields {
    max-width: 20%;
  }

  input[type="radio"] {
    display: inline-block;
    line-height: 3rem;
    padding-top: 0.4rem;
  }

  .new-app > * {
    flex: 1 0;
  }

  .new-app-templates {
    padding: 0 2rem;
  }

  iframe {
    width: 480px;
    height: 320px;
  }

  .app-template {
    padding: 1rem 0;
  }
</style>
