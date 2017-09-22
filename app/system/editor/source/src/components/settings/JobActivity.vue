<template>
  <div class="job-activity">
    <p class="small empty" v-if="!activity.length">No active jobs.</p>
  </div>
</template>

<script>

export default {
  name: 'job-activity',
  computed: {
    activity: function () {
      return this.$store.state.activity.jobs
    }
  }
}
</script>

<style scoped>

  .empty {
    padding: 0 1rem;
  }

</style>
