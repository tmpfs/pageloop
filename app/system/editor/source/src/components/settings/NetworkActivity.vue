<template>
  <div class="network-activity">
    <p class="small empty" v-if="!activity.length">No network activity.</p>
    <ul v-else>
      <li class="item"
        :class="{error: item.error}"
        v-for="item in activity">
        <div>
          <h5>{{item.level}}</h5>
          <p class="small">{{item.message}}</p>
        </div>
        <div v-for="child in item.messages">
          <h5>{{child.level}}</h5>
          <p class="small">{{child.message}}</p>
        </div>
      </li>
    </ul>
  </div>
</template>

<script>

export default {
  name: 'network-activity',
  computed: {
    activity: function () {
      return this.$store.state.network.messages
    }
  }
}
</script>

<style scoped>
  .empty {
    padding: 0 1rem;
  }

  ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
  }

  .item {
    background: var(--base03-color);
    padding: 1rem 0;
  }

  .item > div {
    display: flex;
  }

  .item > div > :last-child {
    flex: 1 0;
  }

  .item > div > :first-child {
    min-width: 8rem;
    text-align: right;
  }

  .item h5 {
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
