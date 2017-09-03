<template>
  <div class="alert" :class="{hidden: !alert.visible}">
    <div class="background"></div>
    <div class="dialog">
      <a class="close" @click="dismiss"></a>
      <h2>{{alert.title}}</h2>
      <div class="dialog-panel">
        <p v-if="alert.message">{{alert.message}}</p>
        <small v-if="alert.note">{{alert.note}}</small>
        <div class="form-actions">
          <button class="sml" @click="dismiss">Cancel</button>
          <button class="sml primary" @click="ok">OK</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app-alert',
  computed: {
    alert: function () {
      return this.$store.state.alert
    }
  },
  methods: {
    dismiss: function () {
      this.$store.commit('alert-hide')
    },
    ok: function () {
      this.alert.ok()
      this.dismiss()
    }
  }
}
</script>

<style scoped>

</style>
