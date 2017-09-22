<template>
  <div class="small">
    <div class="info">
      <span class="meta">{{meta.info.name}} v{{meta.info.version}}</span>
      <span class="uptime">{{meta.stats.uptime.value}}</span>
    </div>
    <div class="stats">
      <div v-for="group, key in meta.stats">
        <h4>{{key}}</h4>
        <ul class="stat">
          <li v-for="v, k in group">
            <span>{{k}}</span>
            <span>{{v}}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'about-settings',
  computed: {
    meta: function () {
      return this.$store.state.meta
    }
  },
  mounted: function () {
    this.$store.dispatch('get-meta')
  }
}
</script>

<style scoped>

  .info {
    clear: both;
    height: 4.2rem;
    margin-bottom: 1rem;
    padding: 1rem 0;
    border-bottom: 1px solid var(--border-color);
  }

  .meta, .uptime {
    display: inline-block;
  }

  .meta {
    float: left;
  }

  .uptime {
    float: right;
  }

  .stats h4 {
    text-transform: uppercase;
    background: var(--base03-color);
    padding-left: 0.5rem;
  }

  ul {
    margin: 0;
    padding: 0;
    list-style-type: none;
    display: table;
    width: 100%;
  }

  li {
    display: table-row;
    width: 100%;
  }

  li > span {
    display: table-cell;
    width: 50%;
  }
</style>
