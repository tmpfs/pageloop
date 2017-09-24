<template>
  <div class="new-app-info">
    <p class="small">Step 1/3: Enter application info</p>
    <form @submit="nextStep">
      <label class="small">Name:</label>
      <input type="text" name="name" placeholder="Enter an app name"
        :value="applicationName" v-model="applicationName" />
      <label class="small">Description:</label>
      <input type="text" name="description" placeholder="Enter an app description"
        :value="applicationDescription" v-model="applicationDescription" />
      <div class="form-actions">
        <input type="submit"
          value="Next: Select a template"
          :class="{disabled: !newAppValid}"
          class="primary" />
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'new-app-info',
  computed: {
    newAppValid: function () {
      return this.$store.state.newApp.valid
    },
    applicationName: {
      get: function () {
        return this.$store.state.newApp.name
      },
      set: function (val) {
        this.$store.state.newApp.name = val
      }
    },
    applicationDescription: {
      get: function () {
        return this.$store.state.newApp.description
      },
      set: function (val) {
        this.$store.state.newApp.description = val
      }
    }
  },
  methods: {
    nextStep: function (e) {
      e.preventDefault()
      if (!this.newAppValid) {
        // TODO: highlight fields
        return
      }
      this.$store.commit('new-app-view', 'new-app-template')
    }
  }
}
</script>

<style scoped>
</style>
