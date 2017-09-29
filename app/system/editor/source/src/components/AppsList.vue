<template>
  <div class="apps-list">
    <div class="small" v-if="!apps.length">No apps found.</div>
    <transition-group appear name="reveal">
      <div
          key="app"
          @click="selectedApp = app"
          class="app"
          v-for="app in apps"
          :class="{selected: app.url === selectedApp.url}">
          <span v-if="app.protected"><i class="fa fa-lock"></i>&nbsp;</span>
          <span class="name">{{app.display}}</span>
          <p class="small">
            {{app.description}}
          </p>
          <p class="app-actions small">
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
      </div>
    </transition-group>
  </div>
</template>

<script>
export default {
  name: 'apps-list',
  props: {
    type: {
      type: String
    }
  },
  computed: {
    apps: function () {
      return this.$store.state.apps[this.type]
    },
    selectedApp: {
      get: function () {
        return this.$store.state.apps.selected
      },
      set: function (val) {
        this.$store.commit('app-list-selected', val)
      }
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
    confirmDeleteApplication: function (app) {
      let details = {
        title: `Delete app (${app.name})`,
        message: `Are you sure you want to permanently delete ${app.name}?`,
        note: 'Be careful deleting an app will remove all app files forever.',
        ok: () => {
          this.deleteApp(app)
        }
      }
      this.$store.commit('alert-show', details)
    },
    deleteApp: function (app) {
      return this.$store.dispatch('del-app', {container: app.container, application: app.name})
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

  .app {
    background: var(--base03-color);
    margin-bottom: 1rem;
    padding: 1rem 1rem 0.8rem 1rem;
    border-bottom: 2px solid var(--base03-color);
    transition: all 0.4s ease-out;
  }

  .app.selected {
    border-bottom: 2px solid var(--base3-color);
    pointer-events: auto;
  }

  .app > p.small {
    margin-bottom: 0;
  }

  .app > .app-actions {
    margin: 0;
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

  .reveal-enter {
		opacity: 0;
  }

  .reveal-enter-active, .reveal-leave-active {
		transition: all 0.4s ease-in;
    opacity: 1;
  }

  .reveal-enter, .reveal-leave-to {
		opacity: 0;
  }

  /*
  .reveal-move {
    transition: transform 1s;
  }
  */

</style>
