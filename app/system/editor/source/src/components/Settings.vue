<template>
  <div class="content-main">
    <div class="content">
      <div class="content-column settings">
        <div class="column-header">
          <h2>Settings</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a
              :class="{selected: settingsView === 'profile-settings'}"
              @click="settingsView = 'profile-settings'">Profile</a>
            <a
              :class="{selected: settingsView === 'organization-settings'}"
              @click="settingsView = 'organization-settings'">Organization</a>
            <a
              :class="{selected: settingsView === 'about-settings'}"
              @click="settingsView = 'about-settings'">About</a>
          </nav>
        </div>
        <div class="scroll">
          <component v-bind:is="settingsView"></component>
        </div>
      </div>
      <div class="content-column preferences">
        <div class="column-header">
          <h2>Preferences</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a
              :class="{selected: preferencesView === 'editor-preferences'}"
              @click="preferencesView = 'editor-preferences'">Editor</a>
            <a
              :class="{selected: preferencesView === 'template-preferences'}"
              @click="preferencesView = 'template-preferences'">Templates</a>
            <a
              :class="{selected: preferencesView === 'storage-preferences'}"
             @click="preferencesView = 'storage-preferences'">Storage</a>
          </nav>
        </div>
        <div class="scroll">
          <component v-bind:is="preferencesView"></component>
        </div>
      </div>
      <div class="content-column activity">
        <div class="column-header">
          <h2>Activity</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a
              :class="{selected: activityView === 'notification-activity'}"
              @click="activityView = 'notification-activity'">Notifications</a>
            <a
              :class="{selected: activityView === 'log-activity'}"
              @click="activityView = 'log-activity'">Logs</a>
            <a
              :class="{selected: activityView === 'network-activity'}"
              @click="activityView = 'network-activity'">Network</a>
          </nav>
        </div>
        <div class="scroll">
          <component v-bind:is="activityView"></component>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ProfileSettings from '@/components/settings/ProfileSettings'
import OrganizationSettings from '@/components/settings/OrganizationSettings'
import AboutSettings from '@/components/settings/AboutSettings'

import EditorPreferences from '@/components/settings/EditorPreferences'
import TemplatePreferences from '@/components/settings/TemplatePreferences'
import StoragePreferences from '@/components/settings/StoragePreferences'

import NotificationActivity from '@/components/settings/NotificationActivity'
import LogActivity from '@/components/settings/LogActivity'
import NetworkActivity from '@/components/settings/NetworkActivity'

export default {
  name: 'settings',
  data: function () {
    return {
      settingsView: 'profile-settings',
      preferencesView: 'editor-preferences',
      activityView: 'notification-activity'
    }
  },
  components: {
    ProfileSettings,
    OrganizationSettings,
    AboutSettings,
    EditorPreferences,
    TemplatePreferences,
    StoragePreferences,
    NotificationActivity,
    LogActivity,
    NetworkActivity
  }
}
</script>

<style scoped>
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
