<template>
  <div class="new-app-info">
    <ul class="new-app-step small">
      <li>
        <span>Application info</span>
        <span>Step 1/3</span>
      </li>
    </ul>
    <form @submit="nextStep">
      <input
        type="text" v-focus
        name="name" placeholder="Enter an app name"
        :value="applicationName" v-model="applicationName" />
      <input
        type="text"
        name="description" placeholder="Enter an app description"
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
  input[type="text"] {
    margin-bottom: 2rem;
  }
</style>
