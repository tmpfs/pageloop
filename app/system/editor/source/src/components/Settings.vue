<template>
  <div class="content-main">
    <div class="content">
      <div class="content-column settings">
        <div class="column-header">
          <h2>Settings</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a>Profile</a>
            <a>Organization</a>
            <a>About</a>
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
            <a>Editor</a>
            <a>Templates</a>
            <a>Storage</a>
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
            <a>Notifications</a>
            <a>Logs</a>
            <a>Network</a>
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

import NotificationsActivity from '@/components/settings/NotificationsActivity'
import LogActivity from '@/components/settings/LogActivity'
import NetworkActivity from '@/components/settings/NetworkActivity'

export default {
  name: 'settings',
  data: function () {
    return {
      settingsView: 'profile-settings',
      preferencesView: 'editor-preferences',
      activityView: 'notifications-activity'
    }
  },
  components: {
    ProfileSettings,
    OrganizationSettings,
    AboutSettings,
    EditorPreferences,
    TemplatePreferences,
    StoragePreferences,
    NotificationsActivity,
    LogActivity,
    NetworkActivity
  }
}
</script>

<style scoped>
  .settings .scroll {
    padding-left: 2rem;
  }

  h3 {
    margin: 0;
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

  .scroll {
    padding: 0;
  }

  ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
  }

  .storage {
    display: table;
    width: 100%;
    margin: 1rem 0;
    font-size: 1.4rem;
  }

  .storage li {
    display: table-row;

  }

  .storage li > * {
    display: table-cell;
  }

  .storage .storage-value {
    text-align: right;
  }

  .item {
    background: var(--base03-color);
    display: flex;
    padding: 1rem 0;
  }

  .item > :last-child {
    flex: 1 0;
  }

  .item :first-child {
    min-width: 8rem;
    text-align: right;
  }

  .item h5 {
    display: inline-block;
    background: transparent;
    font-size: 1.4rem;
    border-right: 2px solid currentColor;
    padding: 0 1rem;
  }

  .item p {
    padding: 0 1rem;
    margin: 0;
  }

  .item.error {
    background: var(--red-color);
    color: var(--base3-color);
  }
</style>
