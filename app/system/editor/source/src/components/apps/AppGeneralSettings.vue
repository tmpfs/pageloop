<template>
  <div class="app-general-settings">
    <h3>Details</h3>
    <ul class="details small">
      <li>
        <span>Name</span>
        <span>{{app && app.display ? app.display : na}}</span>
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
    <h3>Export</h3>
    <div class="export">
      <p class="small">Export your application files to create a backup or snapshot a point in time.</p>
      <nav>
        <a @click="exportSourceArchive()"
           class="download"
           title="Export source files as zip archive">🡇 Download source files</a>
        <a @click="exportPublicArchive()"
           class="download"
           title="Export public files as zip archive">🡇 Download public files</a>
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
  }

  h3 {
    margin: 0;
    padding: 0 0 0.5rem 0;
    border-bottom: 1px solid var(--border-color)
  }

  .details {
    /* margin: 1rem 0 2rem 0; */
    padding: 0;
    list-style-type: none;
    display: table;
    width: 100%;
  }

  .details > li {
    display: table-row;
    width: 100%;
  }

  .details > li > span:first-child {
    text-align: right;
    width: 33%;
  }

  .details > li > span:last-child {
    width: 67%;
    padding-left: 1rem;
  }

  .details > li > span:first-child::after {
    content: ':';
    display: inline-block;
    margin-left: 0.5rem;
  }

  .details > li > span {
    display: table-cell;
    padding: 0;
  }

  .export nav {
    margin-top: 1rem;
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
  }

</style>
