<template>
  <div class="content-main">
    <div class="content">
      <div class="content-column settings">
        <div class="column-header">
          <h2>Settings</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a>User</a>
            <a>Organization</a>
            <a>About</a>
          </nav>
        </div>
        <div class="scroll">
          <p class="small">Lorem ipsum, blah</p>
        </div>
      </div>
      <div class="content-column preferences">
        <div class="column-header">
          <h2>Preferences</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a>Templates</a>
            <a>Storage</a>
          </nav>
        </div>
        <div class="scroll">
          <p class="small">We store the state of the application and your preferences for the User Interface in the browser's local storage. Clearing the local storage will revert the interface to it's default state.</p>

          <p class="small">{{count}} items in storage</p>
          <div class="form-actions">
            <button
              @click="clearLocalStorage"
              :class="{disabled: count == 0}"
              class="primary">Clear Local Storage</button>
          </div>
          <ul class="storage">
            <li v-for="v, k in storage">
              <div class="storage-key">{{k}}</div>
              <div class="storage-value">{{v}}</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="content-column activity">
        <div class="column-header">
          <h2>Activity</h2>
        </div>
        <div class="column-options">
          <nav class="tabs">
            <a>Notifications</a>
            <a>Logs</a>
            <a>Network</a>
          </nav>
        </div>
        <div class="scroll">
          <ul class="activity">
            <p class="small" v-if="!activityNotifications.length">No notifications found.</p>
            <li class="item"
              :class="{error: item.error}"
              v-for="item in activityNotifications">
              <h5 v-if="!item.error">{{item.title}}</h5>
              <h5 v-else>Error</h5>
              <p v-if="!item.error" class="small">{{item.message}}</p>
              <p v-else class="small">{{item.error.message}}</p>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'settings',
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
    },
    activityNotifications: function () {
      return this.$store.state.activity.notifications
    }
  },
  mounted: function () {
    this.storage = window.localStorage
    // console.log('this storage: ')
    // console.log(this.storage)
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
  .settings .scroll {
    padding-left: 2rem;
  }

  h3 {
    margin: 0;
  }

  .content-main {
    border-top: 1px solid var(--border-color);
  }

  .content {
    padding: 0;
  }

  .content > * {
    width: 33.3%;
  }

  .scroll {
    padding: 0 2rem 0 0;
  }

  ul {
    list-style-type: none;
    padding: 0;
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

  h5 {
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
