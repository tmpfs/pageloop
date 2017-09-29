<template>
  <div class="app-general-settings" :class="{disabled: !app}">
    <h3>Details</h3>
    <ul class="details small">
      <li>
        <span>Name</span>
        <span>{{app ? app.display : na}}</span>
      </li>
      <li>
        <span>Description</span>
        <span>{{app ? app.description : na}}</span>
      </li>
      <li>
        <span>Id</span>
        <span>{{app ? app.name : na}}</span>
      </li>
      <li>
        <span>URL</span>
        <span>{{app ? app.url : na}}</span>
      </li>
      <li>
        <span>Container</span>
        <span>{{app ? app.container : na}}</span>
      </li>
      <li>
        <span>Protected</span>
        <span>{{app ? (app.protected ? 'yes' : 'no') : na}}</span>
      </li>
      <li>
        <span>Template</span>
        <span>{{app ? (app['is-template'] ? 'yes' : 'no') : na}}</span>
      </li>
    </ul>
    <div class="build-tasks small" v-if="app && app.build">
      <h3>Build Tasks</h3>
      <div class="task" v-for="task, name in app.build.tasks">
        <span>{{name}}</span>
        <span class="command">$ {{task}}</span>
      </div>
    </div>
    <div class="export">
      <h3>Export</h3>
      <p class="small">Export your application files as a .zip archive to create a backup or snapshot a point in time.</p>
      <nav>
        <a @click="exportFullArchive()"
           class="download"
           :class="{disabled: !app}"
           title="Export all files as zip archive">
          <i class="fa fa-download"></i>
          All files</a>
        <a @click="exportSourceArchive()"
           class="download"
           :class="{disabled: !app}"
           title="Export source files as zip archive">
          <i class="fa fa-download"></i>
          Source files</a>
        <a @click="exportPublicArchive()"
           class="download"
           :class="{disabled: !app}"
           title="Export public files as zip archive">
          <i class="fa fa-download"></i>
          Public files</a>
      </nav>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app-general-settings',
  props: {
    app: {
      type: Object
    },
    na: {
      type: String
    }
  },
  methods: {
    exportFullArchive: function () {
      this.$store.dispatch('app-export', {app: this.app})
    },
    exportSourceArchive: function () {
      this.$store.dispatch('app-export', {app: this.app, filter: 'source'})
    },
    exportPublicArchive: function () {
      this.$store.dispatch('app-export', {app: this.app, filter: 'public'})
    }
  }
}
</script>

<style scoped>

  .app-general-settings {
    padding: 1rem;
    transition: all 0.3s ease-out;
  }

  .app-general-settings:not(.disabled) {
    opacity: 1;
  }

  h3 {
    margin: 0 0 1rem 0;
    padding: 0 0 0.5rem 0;
    border-bottom: 1px solid var(--border-color);
    font-size: 1.5rem;
  }

  .export nav {
    display: flex;
  }

  .export nav a:first-child {
    margin-right: 0.5rem;
  }

  .export nav a:last-child {
    margin-left: 0.5rem;
  }

  .export nav a {
    flex: 1 0;
    background: var(--base03-color);
    padding: 1rem;
    color: var(--base2-color);
    text-align: center;
    font-size: 1.6rem;
    transition: all 0.4s ease-out;
  }

  .build-tasks {
    margin-bottom: 2rem;
  }

  .task {
    padding: 0.5rem;
    background: var(--base03-color);
  }

  .command {
    display: block;
  }

</style>
