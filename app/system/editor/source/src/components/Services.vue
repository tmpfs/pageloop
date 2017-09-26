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

      <div class="content-column method-info">
        <div class="column-header">
          <h2>Method Info</h2>
        </div>
        <div class="scroll" v-if="fn">
          <ul class="details small">
            <li>
              <span>Method</span>
              <span>{{fn.method}}</span>
            </li>
            <li>
              <span>Name</span>
              <span>{{fn.name}}</span>
            </li>
            <li>
              <span>Calls</span>
              <span>{{fn.calls}}</span>
            </li>
            <li>
              <span>Argument Type</span>
              <span>{{fn.arg}}</span>
            </li>
            <li>
              <span>Reply Type</span>
              <span>{{fn.reply}}</span>
            </li>

            <li v-if="fn.info">
              <span>Verb</span>
              <span>{{fn.info.method}}</span>
            </li>
            <li v-if="fn.info">
              <span>Pattern</span>
              <span>{{fn.info.path}}</span>
            </li>
            <li v-if="fn.info">
              <span>Status</span>
              <span>{{fn.info.status}}</span>
            </li>
            <li v-if="fn.info">
              <span>Response Type</span>
              <span>{{getResponseType(fn)}}</span>
            </li>
          </ul>
        </div>
      </div>

      <div class="content-column method-call">
        <div class="column-header">
          <h2>Method Call</h2>
        </div>
        <div class="scroll" v-if="fn">
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'services',
  computed: {
    services: function () {
      return this.$store.state.services.list
    },
    fn: {
      get: function () {
        return this.$store.state.services.method
      },
      set: function (val) {
        this.$store.state.services.method = val
      }
    }
  },
  created: function () {
    this.$store.dispatch('list-services')
  },
  methods: {
    getResponseType: function (fn) {
      const t = fn.info['response-type']
      if (t === 0) {
        return 'json'
      }
      // NOTE: that we treat response type 2 when the service
      // NOTE: function sends the response directly as binary
      return 'binary'
    },
    showServiceMethod: function (service, method) {
      // console.log(method)
      this.fn = method
    }
  }
}
</script>

<style scoped>

  .content-main {
    border-top: 1px solid var(--border-color);
  }

  .scroll {
    padding: 1rem;
    width: 100%;
    height: calc(100% - 2.3rem);
  }

  .service-list {
    width: 20%;
    font-size: 1.4rem;
  }

  .method-info, .method-call {
    width: 40%;
  }

  .methods {
    margin: 0;
    padding: 0 0 0 2rem;
    list-style-type: none;
  }

</style>
