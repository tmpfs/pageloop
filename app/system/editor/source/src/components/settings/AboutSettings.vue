<template>
  <p class="small">
    <span>{{meta.info.name}} v{{meta.info.version}}</span>
    {{meta.stats}}
  </p>
</template>

<script>
export default {
  name: 'about-settings',
  computed: {
    meta: function () {
      return this.$store.state.meta
    }
  },
  mounted: function () {
    this.$store.dispatch('get-meta')
  }
}
</script>

<style scoped>
</style>
