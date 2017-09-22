<template>
  <div class="storage-preferences">
    <p class="small">We store the state of the application and your preferences in the browser's local storage. Clearing the local storage will revert the interface to it's default state.</p>

    <div class="form-actions">
      <button
        @click="clearLocalStorage"
        :class="{disabled: count == 0}"
        class="primary">Reset to defaults</button>
    </div>
    <p class="small title">{{count}} items in storage</p>
    <ul class="storage">
      <li v-for="v, k in storage">
        <div class="storage-key">{{k}}</div>
        <div class="storage-value">{{v}}</div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'storage-preferences',
  data: function () {
    return {
      storage: window.localStorage
    }
  },
  computed: {
    count: function () {
      if (!this.storage) {
        return 0
      }
      return Object.keys(this.storage).length
    }
  },
  watch: {},
  mounted: function () {
    this.storage = window.localStorage
  },
  methods: {
    clearLocalStorage: function () {
      this.$store.commit('clear-local-storage')
      this.storage = null
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
