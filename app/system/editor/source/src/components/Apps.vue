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
              :class="{selected: appListView === 'all'}"
              @click="listApps('all')">All</a>
            <a
              title="Show template applications"
              :class="{selected: appListView === 'template'}"
              @click="listApps('template')">Templates</a>
            <a
              title="Show open applications"
              :class="{selected: appListView === 'open'}"
              @click="listApps('open')">Open</a>
          </nav>
        </div>
        <div class="scroll">
          <component :type="appListView" is="apps-list"></component>
        </div>
      </div>
      <div class="content-column activity">
        <div class="column-header">
          <h2>Settings</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a
              :class="{selected: appSettingsView === 'app-general-settings'}"
              @click="appSettingsView = 'app-general-settings'">General</a>
            <a
              :class="{selected: appSettingsView === 'app-archive-settings'}"
              @click="appSettingsView = 'app-archive-settings'">Archive</a>
            <a
              :class="{selected: appSettingsView === 'app-publish-settings'}"
              @click="appSettingsView = 'app-publish-settings'">Publish</a>
          </nav>
        </div>
        <div class="scroll">
          <component :app="selectedApp" na="--" v-bind:is="appSettingsView"></component>
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

import AppGeneralSettings from '@/components/apps/AppGeneralSettings'
import AppArchiveSettings from '@/components/apps/AppArchiveSettings'
import AppPublishSettings from '@/components/apps/AppPublishSettings'

export default {
  name: 'apps',
  data: function () {
    return {
      appsList: [],
      appListView: 'all',
      appSettingsView: 'app-general-settings'
    }
  },
  computed: {
    selectedApp: {
      get: function () {
        return this.$store.state.apps.selected
      },
      set: function (val) {
        this.$store.commit('app-list-selected', val)
      }
    },
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
      if (type === 'all') {
        this.appsList = this.apps
      } else if (type === 'template') {
        this.appsList = this.templateApps
      } else if (type === 'open') {
        this.appsList = this.openApps
      }
      this.appListView = type
    }
  },
  components: {NewAppInfo, NewAppTemplate, NewAppCreate, AppsList, AppGeneralSettings, AppArchiveSettings, AppPublishSettings}
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

  .settings .scroll, .preferences .scroll {
    padding: 0 1rem;
  }

</style>
