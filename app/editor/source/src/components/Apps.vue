<template>
  <div class="content-main">
    <div class="content scroll">
      <new-app></new-app>
      <div class="containers" v-for="container in list">
        <span :class="{hidden: !container.protected}">🔒&nbsp;</span>
        <span class="name container">{{container.name}}</span>
        <p class="small">{{container.description}}</p>
        <ul class="compact-list">
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
        </ul>
      </div>
    </div>
  </div>
</template>

<script>

import NewApp from '@/components/NewApp'

export default {
  name: 'apps',
  data: function () {
    return {
      applicationName: 'new-app',
      applicationDescription: 'New application'
    }
  },
  computed: {
    list: function () {
      return this.$store.state.containers
    }
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
    },
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
  },
  components: {NewApp}
}
</script>

<style scoped>

  .containers {
    flex: 1 0;
    padding: 2rem;
  }

  .containers > ul {
    margin-left: 2rem;
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

</style>
