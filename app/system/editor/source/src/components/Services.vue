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
                <span title="Number of method calls" class="calls">{{fn.calls}}</span>
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
              <span>Description</span>
              <span>{{fn.meta.description}}</span>
            </li>
            <li>
              <span>Calls</span>
              <span>{{fn.calls}}</span>
            </li>
            <li>
              <span>Argument Type</span>
              <span>{{fn.arg}}</span>
            </li>
            <li v-for="field in fn.fields">
              <span>Field</span>
              <span>{{field.type}} {{field.alias}}</span>
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
          <h3>Websocket</h3>
          <p class="small">Use the JSON-RPC API over the websocket transport.</p>
          <!-- TODO: arguments -->
          <p>
            <a
              @click="callSocketMethod(fn)"
              class="small">Call </a>
          </p>
          <method-reply :reply="socketReply"></method-reply>
          <h3>REST</h3>
          <p class="small">Use the REST API over the HTTP transport.</p>
          <!-- TODO: arguments -->
          <p>
            <a
              @click="callRestMethod(fn)"
              class="small">Call</a>
          </p>
          <method-reply :reply="restReply"></method-reply>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

import {Request} from '../lib/net/client'
import MethodReply from '@/components/ServiceMethodReply'

export default {
  name: 'services',
  components: {MethodReply},
  data: function () {
    return {
      socketReply: null,
      restReply: null
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
    },
    callSocketMethod: function (fn) {
      const params = undefined
      const client = this.$store.state.client
      const req = Request.rpc(fn.method, params)
      client.rpc(req)
        .then((res) => {
          this.socketReply = res
        })
    },
    callRestMethod: function (fn) {
      const params = undefined
      const client = this.$store.state.client
      const req = Request.rpc(fn.method, params)
      client.rpc(req, {http: true})
        .then((res) => {
          this.restReply = res
        })
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
