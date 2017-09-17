<template>
  <ul class="log-activity">
    <p class="small" v-if="!activityLogs.length">No log activity detected.</p>
    <li class="item"
      :class="{error: item.error}"
      v-for="item in activityLogs">
      <h5 v-if="!item.error">{{item.title}}</h5>
      <h5 v-else>Error</h5>
      <p v-if="!item.error" class="small">{{item.message}}</p>
      <p v-else class="small">{{item.error.message}}</p>
    </li>
  </ul>
</template>

<script>

export default {
  name: 'log-activity',
  computed: {
    activityLogs: function () {
      return this.$store.state.activity.log
    }
  }
}
</script>

<style scoped>
  ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
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
