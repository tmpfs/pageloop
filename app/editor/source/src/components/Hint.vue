<template>
  <p class="hint" v-if="message">
    <span v-html="message"></span>
    <a
      @click="dismiss"
      class="close minimal"
      title="Dismiss this hint"></a>
  </p>
</template>

<script>

import BeforeDestroy from '@/components/BeforeDestroy'

export default {
  name: 'hint',
  mixins: [BeforeDestroy],
  props: ['id'],
  computed: {
    message: function () {
      return this.$store.state.hints.messages[this.id]
    }
  },
  methods: {
    dismiss: function () {
      this.$store.commit('dismiss-hint', this.id)
      this.$destroy()
    }
  }
}
</script>
