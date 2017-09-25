<template>
  <div class="app-general-settings" v-if="app">
    <ul class="details small">
      <li v-if="app.display">
        <span>Name</span>
        <span>{{app.display}}</span>
      </li>
      <li>
        <span>Description</span>
        <span>{{app.description}}</span>
      </li>
      <li>
        <span>Id</span>
        <span>{{app.name}}</span>
      </li>
      <li>
        <span>URL</span>
        <span>{{app.url}}</span>
      </li>
      <li>
        <span>Container</span>
        <span>{{app.container}}</span>
      </li>
      <li>
        <span>Protected</span>
        <span>{{app.protected ? 'yes' : 'no'}}</span>
      </li>
      <li>
        <span>Template</span>
        <span>{{app['is-template'] ? 'yes' : 'no'}}</span>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'app-general-settings',
  props: {
    app: {
      type: Object
    }
  }
}
</script>

<style scoped>
  .details {
    margin: 1rem 0 2rem 0;
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
</style>
