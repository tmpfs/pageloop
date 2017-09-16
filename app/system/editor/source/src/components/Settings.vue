<template>
  <div class="content-main">
    <div class="content">
      <div class="content-column">
        <h2>Settings</h2>
        <div class="scroll"></div>
      </div>
      <div class="content-column">
        <h2>Preferences</h2>
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
            <li v-for="v, k in localStorage">
              <div class="storage-key">{{k}}</div>
              <div class="storage-value">{{v}}</div>
            </li>
          </ul>
        </div>
      </div>
      <div class="content-column">
        <h2>Activity</h2>
        <div class="scroll">
          <h3>Notifications</h3>
          <ul class="activity">
            <li class="item"
              v-for="item in activityNotifications">
              <div v-if="!item.error">
                <h5>{{item.title}}</h5>
                <p class="small">{{item.message}}</p>
              </div>
              <div class="error"
                v-else>
                <h5>Error</h5>
                <p class="small">{{item.error.message}}</p>
              </div>
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
  computed: {
    localStorage: function () {
      return localStorage
    },
    count: function () {
      return Object.keys(this.localStorage).length
    },
    activityNotifications: function () {
      return this.$store.state.activity.notifications
    }
  },
  methods: {
    clearLocalStorage: function () {
      this.$store.commit('clear-local-storage')
    }
  }
}
</script>

<style scoped>
  h2 {
    margin: 0;
    height: 4.6rem;
    padding: 1.3rem 0;
    font-size: 2rem;
  }

  h3 {
    margin: 0 0 1rem 0;
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

  .content-column:not(:last-child) {
    margin-right: 2rem;
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
  }

  .item :first-child {
    position: relative;
  }

  h5 {
    display: inline-block;
    position: absolute;
    top: 0;
    left: 1rem;
    background: transparent;
    font-size: 1.5rem;
  }

  .item p {
    padding: 0 1rem;
    margin: 0;
  }

  .item .error {
    background: var(--red-color);
    color: var(--base3-color);
  }
</style>
