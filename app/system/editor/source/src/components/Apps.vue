<template>
  <div class="content-main">
    <div class="content">
      <div class="content-column settings">
        <div class="column-header">
          <h2>New App<span v-if="newAppId">/ {{newAppId}}</span></h2>
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
          <h2>App List</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a
              @click="listApplications(container, $event)"
              :title="getLinkTitle(container)"
              :class="{selected: isSelected(container)}"
              v-for="container in list" v-if="enabled[container.name]">{{container.name}}</a>
          </nav>
        </div>
        <div class="scroll">
          <component :containerName="containerName" is="apps-list"></component>
        </div>
      </div>
      <div class="content-column activity">
        <div class="column-header">
          <h2>App Settings</h2>
        </div>
        <div class="column-options">
          <!--
          <nav class="tabs">
            <a
              :class="{selected: activityView === 'notification-activity'}"
              @click="activityView = 'notification-activity'">Notifications</a>
            <a
              :class="{selected: activityView === 'job-activity'}"
              @click="activityView = 'job-activity'">Jobs</a>
            <a
              :class="{selected: activityView === 'log-activity'}"
              @click="activityView = 'log-activity'">Logs</a>
            <a
              :class="{selected: activityView === 'network-activity'}"
              @click="activityView = 'network-activity'">Network</a>
          </nav>
          -->
        </div>
        <!--
        <div class="scroll">
          <component v-bind:is="activityView"></component>
        </div>
        -->
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
      currentView: 'apps-list',
      containerName: 'user',
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

      o.user = this.system || this.template

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
  methods: {
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
  },
  components: {NewAppInfo, NewAppTemplate, NewAppCreate, AppsList}
}
</script>

<style>

  .apps-list {
    /*
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    padding: 2rem;
    width: 100%;
    */
  }

  .app {
    /* flex: 1 0;
    margin-right: 2rem;
    min-width: 24rem; */
  }

  /*
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
  */

  .apps-nav {
    border-top: 1px solid var(--border-color);
    border-bottom: 1px solid var(--border-color);
  }

  .apps-nav .tabs {
    max-width: 40%;
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
