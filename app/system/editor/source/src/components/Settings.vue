<template>
  <div class="content-main">
    <div class="content">
      <div class="content-column scroll">
        <h2>Settings</h2>
      </div>
      <div class="content-column scroll">
        <h2>Preferences</h2>
        <h3>Local Storage</h3>
        <p class="small">We store the state of the application and your preferences for the User Interface in the browser's local storage. Clearing the local storage will revert the interface to it's default state.</p>
        <div class="form-actions"><button class="primary">Clear Local Storage</button></div>
        <ul class="storage">
          <li v-for="v, k in localStorage">
            <div class="storage-key">{{k}}</div>
            <div class="storage-value">{{v}}</div>
          </li>
        </ul>
      </div>
      <div class="content-column scroll">
        <h2>Activity</h2>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'settings',
  computed: {
    localStorage: function () {
      return localStorage
    }
  }
}
</script>

<style scoped>
  h2 {
    margin: 2rem 0;
  }

  .content-main {
    border-top: 1px solid var(--border-color);
  }

  .content {
    padding: 0 2rem;
  }

  .content > * {
    width: 33.3%;
  }

  .content-column:not(:last-child) {
    margin-right: 2rem;
  }

  .storage {
    display: table;
    width: 100%;
    list-style-type: none;
    padding: 0;
    margin: 1rem 0;
    font-size: 1.4rem;
  }

  .storage li {
    display: table-row;

  }

  .storage li > * {
    display: table-cell;
  }

  .storage .storage-value {
    text-align: right;
  }
</style>
