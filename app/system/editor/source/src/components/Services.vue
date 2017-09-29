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
                  @click="showServiceMethod(srv, fn)"
                  :title="fn.method">
                  <span class="fn">{{fn.method}}</span>
                </a>
                <span v-bind:calls="fn.calls" title="Number of method calls" class="calls">{{fn.calls}}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <div class="content-column method-info">
        <div class="column-header">
          <h2>Info</h2>
        </div>
        <div class="scroll" v-if="fn">
          <ul class="details small">
            <li>
              <span>Service</span>
              <span>{{fn.service}}</span>
            </li>
            <li>
              <span>Method</span>
              <span>{{fn.method}}</span>
            </li>
            <li>
              <span>Description</span>
              <span>{{fn.meta.description}}</span>
            </li>
            <li>
              <span>Calls</span>
              <span>{{fn.calls}}</span>
            </li>
            <li>
              <span>Argument Type</span>
              <span class="type">{{fn.arg}}</span>
            </li>
            <li>
              <span>Reply Type</span>
              <span class="type">{{fn.reply}}</span>
            </li>
          </ul>

          <div v-if="fn.fields">
            <hr />
            <ul class="details small">
              <li>
                <span>Argument Fields</span>
                <span class="type">{{fn.arg}}</span>
              </li>
              <li v-for="field in fn.fields">
                <span>{{field.alias}}</span>
                <span class="type">{{field.type}}</span>
              </li>
            </ul>
          </div>

          <div v-if="fn.info" v-for="route in fn.info">
            <hr />
            <ul class="details small">
              <li>
                <span>Verb</span>
                <span>{{route.method}}</span>
              </li>
              <li>
                <span>Pattern</span>
                <span>{{route.path}}</span>
              </li>
              <li>
                <span>Status</span>
                <span>{{route.status}}</span>
              </li>
              <li>
                <span>Response Type</span>
                <span>{{getResponseType(route)}}</span>
              </li>
            </ul>
          </div>

        </div>
      </div>

      <div class="content-column method-call">
        <div class="column-header">
          <h2 v-if="fn" class="method-name">{{fn.method}}</h2>
          <h2 v-else>Call</h2>
          <nav class="tabs" v-if="fn">
            <a
              :class="{selected: callType === 'websocket'}"
              @click="callType = 'websocket'">
              Websocket</a>
            <a
              :class="{selected: callType === 'rest'}"
              @click="callType = 'rest'">
              Rest</a>
          </nav>
        </div>
        <div class="scroll" v-if="fn">
          <method-call :fn="fn" :callType="callType"></method-call>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

import MethodCall from '@/components/ServiceMethodCall'

export default {
  name: 'services',
  components: {MethodCall},
  data: function () {
    return {
      callType: 'websocket'
    }
  },
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
    },
    params: {
      get: function () {
        return this.$store.state.services.params
      },
      set: function (val) {
        this.$store.state.services.params = val
      }
    }
  },
  created: function () {
    this.$store.dispatch('list-services')
  },
  methods: {
    getResponseType: function (route) {
      const t = route['response-type']
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

      // Set up default parameters for binding
      let params = {}
      if (method.fields) {
        method.fields.forEach((field) => {
          params[field.alias] = ''
        })
      }
      this.params = params
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

  .method-name {
    text-transform: none;
  }

  .tabs > :first-child {
    border-left: 1px solid var(--border-color);
  }

  h3, h4 {
    margin: 0;
    padding: 0 0 0.5rem 0;
    border-bottom: 1px solid var(--border-color);
    font-size: 1.5rem;
  }

  h4 {
    display: inline-block;
    font-size: 1.4rem;
  }

  .service-list {
    font-size: 1.4rem;
    user-select: none;
  }

  .methods {
    margin: 0;
    padding: 0 0 0 2rem;
    list-style-type: none;
  }

  .methods li {
    display: flex;
  }

  .methods li > * {
    flex: 1 0;
  }

  .methods li > span.calls {
    text-align: right;
  }

</style>

<style>

  .type {
    color: var(--yellow-color);
  }

</style>
