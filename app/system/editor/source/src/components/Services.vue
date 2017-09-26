<template>
  <div class="content-main">
    <div class="content">
      <div class="content-column service-list">
        <div class="column-header">
          <h2>Services</h2>
        </div>
        <div class="scroll">
          <div v-for="srv, id in services">
            {{srv.name}}
            <ul class="methods">
              <li v-for="fn in srv.methods">
                <a
                  @click="showServiceMethod(srv, fn)">
                  <span class="fn">{{fn.method}}</span>
                  <span class="calls">{{fn.calls}}</span>
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!--{{services}}-->
    </div>
  </div>
</template>

<script>
export default {
  name: 'services',
  computed: {
    services: function () {
      return this.$store.state.services
    }
  },
  mounted: function () {
    this.$store.dispatch('list-services')
  },
  methods: {
    showServiceMethod: function (service, method) {
      console.log(method)
    }
  }
}
</script>

<style scoped>
  .scroll {
    padding: 1rem;
  }

  .service-list {
    width: 50%;
    font-size: 1.4rem;
  }

  ul {
    margin: 0;
    padding: 0 0 0 2rem;
    list-style-type: none;
  }

</style>
