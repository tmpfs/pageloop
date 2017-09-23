<template>
  <div class="storage-preferences" v-bind="keys">
    <p class="small">We store the state of the application and your preferences in the browser's local storage. Clearing the local storage will revert the interface to it's default state.</p>

    <div class="form-actions">
      <button
        @click="clearLocalStorage"
        :class="{disabled: count === 0}"
        class="primary">Reset to defaults</button>
    </div>
    <p class="small title">{{count}} items in storage</p>
    <ul class="storage">
      <li v-for="_, key in keys" :class="{disabled: localStorage[key] === undefined}">
        <div class="storage-key">{{key}}</div>
        <div class="storage-value">{{settings.get(key, true)}}</div>
      </li>
    </ul>
  </div>
</template>

<script>

export default {
  name: 'storage-preferences',
  data: function () {
    return {
      localStorage: window.localStorage,
      count: window.localStorage.length,
      settings: this.$store.state.settings,
      keys: this.$store.state.settings.keys
    }
  },
  updated: function () {
    this.count = localStorage.length
  },
  methods: {
    clearLocalStorage: function () {
      this.$store.commit('clear-local-storage')
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

  .title {
    border-bottom: 1px solid var(--border-color);
  }

  .storage {
    display: table;
    width: 100%;
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
