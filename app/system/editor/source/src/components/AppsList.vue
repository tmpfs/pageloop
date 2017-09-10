<template>
  <div class="apps-list scroll">
   <div class="app" v-for="app in container.apps">
        <span :class="{hidden: !app.protected}">🔒&nbsp;</span>
        <span class="name">{{app.name}}</span>
        <p class="small">URL: {{app.url}}<br />{{app.description}}
          <p class="app-actions">
            <a class="name"
              @click="$store.dispatch('navigate', {href: linkify(container, app)})"
              :title="title(app, 'Edit')">Edit</a>
            <a class="name"
              :href="linkify(container, app, true)"
              :title="title(app, 'Open')">Open</a>
            <a v-if="!app.protected" class="name"
              @click="confirmDeleteApplication(container, app)"
              :title="title(app, 'Delete')">Delete</a>
          </p>
        </p>
    </div>
  </div>
</template>

<script>

export default {
  name: 'apps-list',
  props: ['containerName'],
  computed: {
    list: function () {
      return this.$store.state.containers
    },
    container: function () {
      return this.$store.state.getContainerByName(this.containerName)
    }
  },
  methods: {
    confirmDeleteApplication: function (container, application) {
      let details = {
        title: `Delete Application (${application.name})`,
        message: `Are you sure you want to permanently delete ${application.name}?`,
        note: 'Be careful deleting an application will remove all application files forever.',
        ok: () => {
          this.deleteApplication(container, application)
        }
      }
      this.$store.commit('alert-show', details)
    },
    deleteApplication: function (container, application) {
      this.$store.dispatch('del-app', {container: container.name, application: application.name})
        .catch((e) => console.error(e))
      return false
    },
    linkify: function (c, a, open) {
      if (open) {
        return a.url
      }
      return `apps/${c.name}/${a.name}`
    },
    title: function (a, prefix) {
      return `${prefix} ${a.name}`
    }
  }
}
</script>

<style scoped>

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
