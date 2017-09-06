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
export default {
  name: 'hint',
  props: ['id'],
  computed: {
    message: function () {
      return this.$store.state.hints.messages[this.id]
    }
  },
  methods: {
    dismiss: function () {
      this.$store.commit('dismiss-hint', this.id)
    }
  }
}
</script>

<style scoped>
  .hint {
    position: relative;
    background: var(--base03-color);
    padding: 1rem;
    font-size: 1.4rem;
  }

  .hint.compact {
    margin: 0;
  }

  .hint::before {
    display: inline-block;
    content: 'ℹ';
    color: var(--cyan-color);
    margin-right: 0.5rem;
    text-align: center;
    line-height: 2rem;
    background: var(--base02-color);
    width: 1.8rem;
    height: 1.8rem;
    border-radius: 0.9rem;
  }

  .hint > a.close {
		position: absolute;
		top: 0.2rem;
		right: 0.2rem;
  }
</style>
