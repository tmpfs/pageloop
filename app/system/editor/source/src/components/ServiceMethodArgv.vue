<template>
  <div v-if="fn.fields">
    <ul class="details small">
      <li>
        <span>Argument Fields</span>
        <span class="type">{{fn.arg}}</span>
      </li>
      <li v-for="field in fn.fields">
        <span>{{field.alias}}</span>
        <span
          :data-type="field.type"
          :data-alias="field.alias"
          contenteditable
          @keyup="keyup"
          @keydown="keydown"
          @keyup.enter="enter"
          class="input">{{params[field.alias]}}</span>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'method-argv',
  props: {
    fn: {
      type: Object
    }
  },
  computed: {
    params: {
      get: function () {
        return this.$store.state.services.params
      },
      set: function (val) {
        this.$store.state.services.params = val
      }
    }
  },
  methods: {
    keyup: function (e) {
      const el = e.currentTarget
      const alias = el.getAttribute('data-alias')
      const value = el.innerText

      // TODO: type coercion
      // const type = el.getAttribute('data-type')

      this.params[alias] = value
    },
    keydown: function (e) {
      if (e.key === 'Enter') {
        e.preventDefault()
        e.stopImmediatePropagation()
      }
    },
    enter: function (e) {
      e.preventDefault()
      e.stopImmediatePropagation()
      this.$emit('submit')
    }
  }
}
</script>

<style scoped>
  .input {
    color: var(--cyan-color);
  }
</style>
