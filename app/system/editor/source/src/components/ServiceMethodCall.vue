<template>
  <div v-if="fn">
    <method-argv :fn="fn"></method-argv>
    <p>
      <button
        @click="callMethod(fn)"
        class="small">Call {{fn.method}}</button>
    </p>
    <method-reply :reply="reply"></method-reply>
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
  computed: {
    reply: function () {
      if (this.callType === 'websocket') {
        return this.socketReply
      } else if (this.callType === 'rest') {
        return this.restReply
      }
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
  watch: {
    fn: function () {
      this.restReply = null
      this.socketReply = null
    }
  },
  methods: {
    callMethod: function (fn) {
      if (this.callType === 'websocket') {
        return this.callSocketMethod(fn)
      } else if (this.callType === 'rest') {
        return this.callRestMethod(fn)
      }
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
