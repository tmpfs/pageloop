<template>
  <div class="apps-list">
   <div class="app" v-for="app in apps">
        <span :class="{hidden: !app.protected}">🔒&nbsp;</span>
        <span class="name">{{app.display || app.name}}</span>
        <p class="small">URL: {{app.url}}<br />{{app.description}}
          <p class="app-actions">
            <a class="name"
              @click="editApplication(app)"
              :title="title(app, 'Edit')">Edit</a>
            <a class="name"
              :href="linkify(app, true)"
              :title="title(app, 'Open')">Open</a>
            <a v-if="!app.protected" class="name"
              @click="confirmDeleteApplication(app)"
              :title="title(app, 'Delete')">Delete</a>
          </p>
        </p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'apps-list',
  props: {
    apps: {
      type: Array
    }
  },
  methods: {
    getContainer: function (app) {
      return this.$store.state.getContainerByName(app.container)
    },
    editApplication: function (app) {
      const container = this.getContainer(app)
      return this.$store.dispatch('edit-app', {container: container, application: app})
    },
    confirmDeleteapp: function (app) {
      const container = this.getContainer(app)
      let details = {
        title: `Delete app (${app.name})`,
        message: `Are you sure you want to permanently delete ${app.name}?`,
        note: 'Be careful deleting an app will remove all app files forever.',
        ok: () => {
          this.deleteapp(container, app)
        }
      }
      this.$store.commit('alert-show', details)
    },
    deleteapp: function (app) {
      const container = this.getContainer(app)
      this.$store.dispatch('del-app', {container: container.name, app: app.name})
        .catch((e) => console.error(e))
      return false
    },
    linkify: function (app, open) {
      const container = this.getContainer(app)
      if (open) {
        return app.url
      }
      return `apps/${container.name}/${app.name}`
    },
    title: function (app, prefix) {
      return `${prefix} ${app.name}`
    }
  }
}
</script>

<style scoped>

  .apps-list {
    padding-top: 1rem;
  }

  .app > p.small {
    margin-bottom: 0;
  }

  .app > .app-actions {
    margin-top: 0;
    font-size: 1.5rem;
  }

  .app-actions > *:not(:last-child) {
    margin-right: 1rem;
  }

  .name.container, .new-app h2 {
    font-size: 1.4rem;
    text-decoration: underline;
    text-transform: uppercase;
  }

  .name + p.small {
    margin-top: 0.2rem;
  }

</style>
