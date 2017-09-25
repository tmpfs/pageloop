<template>
  <div class="content-main">
    <div class="content">
      <div class="content-column settings">
        <div class="column-header">
          <h2>New<span v-if="newAppId">/ {{newAppId}}</span></h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a
              :class="{selected: newAppView === 'new-app-info'}"
              @click="newAppView = 'new-app-info'">Info</a>
            <a
              :class="{selected: newAppView === 'new-app-template', disabled: !newAppValid}"
              @click="newAppView = 'new-app-template'">Template</a>
            <a
              :class="{selected: newAppView === 'new-app-create', disabled: !newAppValid}"
              @click="newAppView = 'new-app-create'">Create</a>
          </nav>
        </div>
        <div class="scroll">
          <component v-bind:is="newAppView"></component>
        </div>
      </div>
      <div class="content-column preferences">
        <div class="column-header">
          <h2>List</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a
              title="Show all applications"
              :class="{selected: listAppView === 'all'}"
              @click="listApps('all')">All</a>
            <a
              title="Show template applications"
              :class="{selected: listAppView === 'templates'}"
              @click="listApps('templates')">Templates</a>
            <a
              title="Show open applications"
              :class="{selected: listAppView === 'open'}"
              @click="listApps('open')">Open</a>

            <!--
            <a
              @click="listApplications(container, $event)"
              :title="getLinkTitle(container)"
              :class="{selected: isSelected(container)}"
              v-for="container in list" v-if="enabled[container.name]">{{container.name}}</a>
            -->
          </nav>
        </div>
        <div class="scroll">
          <!--<component :containerName="containerName" is="apps-list"></component>-->
          <component :apps="apps" is="apps-list"></component>
        </div>
      </div>
      <div class="content-column activity">
        <div class="column-header">
          <h2>Settings</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a
              :class="{selected: appSettingsView === 'general-app-settings'}"
              @click="appSettingsView = 'general-app-settings'">General</a>
            <a
              :class="{selected: appSettingsView === 'archive-app-settings'}"
              @click="appSettingsView = 'archive-app-settings'">Archive</a>
            <a
              :class="{selected: appSettingsView === 'build-app-settings'}"
              @click="appSettingsView = 'build-app-settings'">build</a>
          </nav>
        </div>
        <div class="scroll">
          <component v-bind:is="appSettingsView"></component>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

import NewAppInfo from '@/components/NewAppInfo'
import NewAppTemplate from '@/components/NewAppTemplate'
import NewAppCreate from '@/components/NewAppCreate'
import AppsList from '@/components/AppsList'

export default {
  name: 'apps',
  data: function () {
    return {
      apps: [],
      currentView: 'apps-list',
      appListView: 'all',
      appSettingsView: '',
      /* containerName: 'user', */
      user: true
    }
  },
  computed: {
    newAppView: {
      get: function () {
        return this.$store.state.newApp.view
      },
      set: function (val) {
        this.$store.commit('new-app-view', val)
      }
    },
    newAppValid: function () {
      return this.$store.state.newApp.valid
    },
    newAppId: function () {
      return this.$store.state.newApp.id
    },
    list: function () {
      return this.$store.state.containers
    },
    enabled: function () {
      const o = {
        system: this.system,
        template: this.template
      }

      // o.user = this.system || this.template

      return o
    },
    system: {
      get: function () {
        return this.$store.state.settings.showSystemApplications
      },
      set: function (val) {
        this.$store.state.settings.showSystemApplications = val
      }
    },
    template: {
      get: function () {
        return this.$store.state.settings.showTemplateApplications
      },
      set: function (val) {
        this.$store.state.settings.showTemplateApplications = val
      }
    }
  },
  mounted: function () {
    this.listApps('all')
  },
  methods: {
    listApps: function (type) {
      let apps = []
      this.list.forEach((container) => {
        if (this.enabled[container.name] !== undefined && !this.enabled[container.name]) {
          return
        }
        apps = apps.concat(container.apps)
      })

      if (type === 'templates') {
        apps = apps.filter((app) => {
          return app['is-template']
        })
      } else if (type === 'open') {
        apps = apps.filter((app) => {
          return app.open
        })
      }
      this.appListView = type
      this.apps = apps
    }

    /*
    isSelected: function (container) {
      return this.currentView === 'apps-list' && this.containerName === container.name
    },
    listApplications: function (container, e) {
      e.preventDefault()
      this.containerName = container.name
      this.currentView = 'apps-list'
    },
    getLinkTitle: function (container) {
      return `Show applications in ${container.name}`
    }
    */
  },
  components: {NewAppInfo, NewAppTemplate, NewAppCreate, AppsList}
}
</script>

<style>
  .new-app-step {
    display: table;
    width: 100%;
    margin: 1rem 0 2rem 0;
    padding: 0;
  }

  .new-app-step li {
    display: table-row;
    width: 100%;
    background: var(--base03-color);
  }

  .new-app-step > li > span {
    display: table-cell;
    padding: 1rem;
  }

  .new-app-step > li > span:last-child {
    text-align: right;
  }
</style>

<style scoped>
  .scroll {
    height: calc(100% - 4.6rem);
  }

  .content-main {
    border-top: 1px solid var(--border-color);
  }

  .content {
    padding: 0;
  }

  .content > * {
    width: 33.3%;
  }

  .content-column:not(:last-child) {
    border-right: 1px solid var(--border-color);
  }

  .settings .scroll, .preferences .scroll {
    padding: 0 1rem;
  }

</style>
