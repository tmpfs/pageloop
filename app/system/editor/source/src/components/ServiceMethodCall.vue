<template>
  <div v-if="fn">
    <!--
    <p v-if="callType === 'websocket'" class="small center">JSON-RPC API over the websocket transport</p>
    <p v-else class="small center">REST API over the HTTP transport</p>
    -->
    <method-argv :fn="fn"></method-argv>
    <div v-if="callType === 'websocket'">
      <p>
        <a
          @click="callSocketMethod(fn)"
          class="small">Call </a>
      </p>
      <method-reply :reply="socketReply"></method-reply>
    </div>

    <div v-if="callType === 'rest'">
      <p>
        <a
          @click="callRestMethod(fn)"
          class="small">Call</a>
      </p>
      <method-reply :reply="restReply"></method-reply>
    </div>
  </div>
</template>

<script>

import {Request} from '../lib/net/client'
import MethodReply from '@/components/ServiceMethodReply'
import MethodArgv from '@/components/ServiceMethodArgv'

export default {
  name: 'services',
  components: {MethodReply, MethodArgv},
  data: function () {
    return {
      socketReply: null,
      restReply: null
    }
  },
  props: {
    fn: {
      type: Object
    },
    callType: {
      type: String
    }
  },
  methods: {
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
