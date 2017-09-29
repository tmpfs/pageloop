<template>
  <div v-if="fn">
    <method-argv :fn="fn" v-on:submit="submit"></method-argv>
    <p>
      <button
        @click="invoke(fn)"
        class="small">Call {{fn.method}}</button>
    </p>
    <method-reply :fn="fn" :reply="reply"></method-reply>
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
      }
    }
  },
  computed: {
    reply: function () {
      return this.replies[this.callType]
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
    }
  },
  methods: {
    submit: function () {
      this.invoke(this.fn)
    },
    invoke: function (fn) {
      const params = this.params
      const client = this.$store.state.client
      const req = Request.rpc(fn.method, params)
      client.rpc(req, {http: this.callType === 'rest'})
        .then((res) => {
          this.replies[this.callType] = res

          // Update the num calls
          const req = Request.rpc('Service.ReadMethodCalls', {service: fn.service.toLowerCase(), method: fn.name.toLowerCase()})
          client.rpc(req, {http: this.callType === 'rest'})
            .then((res) => {
              if (res.response.status === 200) {
                fn.calls = res.document
              }
            })
        })
    }
  }
}
</script>
