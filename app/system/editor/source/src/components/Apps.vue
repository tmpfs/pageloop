<template>
  <div class="content-main">
    <div class="apps-nav">
      <nav class="tabs">
        <a
          @click="currentView = 'new-app'"
          :class="{selected: currentView === 'new-app'}"
          title="Create a new application">➕ New Application</a>

        <a
          @click="listApplications(container, $event)"
          :title="getLinkTitle(container)"
          :class="{selected: isSelected(container)}"
          v-for="container in list">{{container.name}} Apps</a>
      </nav>
    </div>
    <div class="content">
      <component v-bind:is="currentView" :containerName="containerName"></component>
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
      currentView: 'apps-list',
      containerName: 'user'
    }
  },
  computed: {
    list: function () {
      return this.$store.state.containers
    }
  },
  methods: {
    isSelected: function (container) {
      return this.currentView === 'apps-list' && this.containerName === container.name
    },
    listApplications: function (container, e) {
      e.preventDefault()
      this.containerName = container.name
      this.currentView = 'apps-list'
    },
    getLinkTitle: function (container) {
      return `Show applications in ${container.name}`
    }
  },
  components: {NewApp, AppsList}
}
</script>

<style>

  .apps-list {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    padding: 2rem;
    width: 100%;
  }

  .app {
    flex: 1 0;
    margin-right: 2rem;
    min-width: 24rem;
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

  .apps-nav {
    border-top: 1px solid var(--border-color);
    border-bottom: 1px solid var(--border-color);
  }

  .apps-nav .tabs {
    max-width: 40%;
  }

</style>

<style scoped>
  .scroll {
    height: calc(100% - 2.4rem);
  }

</style>
