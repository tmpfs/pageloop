<template>
  <div v-if="fn">
    <method-argv :fn="fn" :params="params"></method-argv>
    <p>
      <button
        @click="invoke(fn)"
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
      replies: {
        websocket: null,
        rest: null
      },
      params: {}
    }
  },
  computed: {
    reply: function () {
      return this.replies[this.callType]
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
      // Reset replies when the service method changes
      this.replies = {websocket: null, rest: null}
      this.params = {}
    }
  },
  methods: {
    invoke: function (fn) {
      const params = this.params
      const client = this.$store.state.client
      const req = Request.rpc(fn.method, params)
      console.log(this.callType === 'rest')
      client.rpc(req, {http: this.callType === 'rest'})
        .then((res) => {
          this.replies[this.callType] = res
        })
    }
  }
}
</script>
