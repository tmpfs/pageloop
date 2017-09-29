<template>
  <div class="profile-settings small">
    <div>
      <input
        id="show-notifications"
        type="checkbox"
        v-model="showNotifications" />
      <label for="show-notifications">Show notifications</label>
      <p class="small">When notifications are disabled you will not be notified as events occur but can view them in the notification activity (errors are always shown).</p>
    </div>
    <div>
      <input
        id="show-system"
        type="checkbox"
        v-model="showSystemApplications" />
      <label for="show-system">List system applications</label>
    </div>
    <div>
      <input
        id="show-template"
        type="checkbox"
        v-model="showTemplateApplications" />
      <label for="show-template">List template applications</label>
    </div>
  </div>
</template>

<script>
export default {
  name: 'profile-settings',
  computed: {
    showNotifications: {
      get: function () {
        return this.$store.state.settings.showNotifications
      },
      set: function (val) {
        this.$store.state.settings.showNotifications = val
      }
    },
    showSystemApplications: {
      get: function () {
        return this.$store.state.settings.showSystemApplications
      },
      set: function (val) {
        this.$store.commit('show-system-applications', val)
      }
    },
    showTemplateApplications: {
      get: function () {
        return this.$store.state.settings.showTemplateApplications
      },
      set: function (val) {
        this.$store.commit('show-template-applications', val)
      }
    }
  }
}
</script>

<style scoped>
  .profile-settings {
    margin-top: 1rem;
  }

  p.small {
    margin-top: 0;
    margin-left: 2.5rem;
  }

  label {
    user-select: none;
  }

</style>
