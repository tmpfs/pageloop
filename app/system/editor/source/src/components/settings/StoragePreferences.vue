<template>
  <div class="storage-preferences">
    <p class="small">We store the state of the application and your preferences in the browser's local storage. Clearing the local storage will revert the interface to it's default state.</p>

    <div class="form-actions">
      <button
        @click="clearLocalStorage"
        v-bind:length="settings.length"
        :class="{disabled: count === 0}"
        class="primary">Reset to defaults</button>
    </div>
    <p class="small title" v-bind="settings.length">{{count}} items in storage</p>
    <ul class="storage">
      <li v-for="_, key in keys" :class="{disabled: localStorage[key] === undefined}" v-bind="storage">
        <div class="storage-key">{{key}}</div>
        <div class="storage-value" v-bind:val="storage[key]">{{settings.get(key)}}</div>
      </li>
    </ul>
  </div>
</template>

<script>

export default {
  name: 'storage-preferences',
  data: function () {
    return {
      count: this.$store.state.settings.length,
      localStorage: window.localStorage,
      settings: this.$store.state.settings,
      storage: this.$store.state.settings.storage,
      keys: this.$store.state.settings.keys
    }
  },
  methods: {
    clearLocalStorage: function () {
      this.$store.commit('clear-local-storage')
      const cache = this.$store.state.settings.storage
      this.storage = null
      this.storage = cache
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
