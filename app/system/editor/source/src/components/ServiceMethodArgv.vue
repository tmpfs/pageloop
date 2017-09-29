<template>
  <div v-if="fn.fields">
    <ul class="details small">
      <li>
        <span>Argument Fields</span>
        <span class="type">{{fn.arg}}</span>
      </li>
      <li v-for="field in fn.fields">
        <span><span class="type">{{field.type}}</span> {{field.alias}}</span>
        <typed-input
          @blur="blur"
          v-on:submit="enter"
          :type="field.type"
          :name="field.alias"
          v-bind:value="params[field.alias]"
          :value="params[field.alias]"></typed-input>
        <!--
        <span v-if="field.fields">
          {{field.fields}}
        </span>
        -->
      </li>
    </ul>
  </div>
</template>

<script>

import TypedInput from '@/components/TypedInput'

export default {
  name: 'method-argv',
  components: {TypedInput},
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
    focus: function (e) {
      const sel = getSelection()
      sel.selectAllChildren(e.target)
    },
    blur: function (e, input) {
      this.setParam(e.currentTarget, input)
    },
    keydown: function (e) {
      if (e.key === 'Enter') {
        e.preventDefault()
        e.stopImmediatePropagation()
      }
    },
    setParam: function (el, input) {
      const alias = input.name
      this.params[alias] = input.getValue()
    },
    update: function (input) {
      const fields = this.$el.querySelectorAll('.field')
      fields.forEach((el) => {
        this.setParam(el, input)
      })
    },
    enter: function (e, input) {
      this.update(input)
      this.$emit('submit')
    }
  }
}
</script>

<style scoped>
</style>
