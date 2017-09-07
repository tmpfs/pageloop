<template>
  <div class="content-main">
    <nav class="tabs">
      <a
        @click="currentView = 'new-app'"
        title="Create a new application">New Application</a>
      <a
        @click="currentView = 'apps-list'"
        title="List application">Show Applications</a>
    </nav>
    <div class="content scroll">
      <component v-bind:is="currentView"></component>
    </div>
  </div>
</template>

<script>

import NewApp from '@/components/NewApp'
import AppsList from '@/components/AppsList'

export default {
  name: 'apps',
  data: function () {
    return {
      currentView: 'apps-list'
    }
  },
  components: {NewApp, AppsList}
}
</script>

<style scoped>

  .containers-list {
    display: flex;
  }

  .containers {
    flex: 1 0;
    padding: 2rem;
  }

  .containers > ul {
    margin-left: 2rem;
  }

  .app > p.small {
    margin-bottom: 0;
  }

  .app > .app-actions {
    margin-top: 0;
    font-size: 1.5rem;
  }

  .app-actions > *:not(:last-child) {
    margin-right: 1rem;
  }

  .name.container, .new-app h2 {
    font-size: 1.4rem;
    text-decoration: underline;
    text-transform: uppercase;
  }

  .name + p.small {
    margin-top: 0.2rem;
  }

</style>
