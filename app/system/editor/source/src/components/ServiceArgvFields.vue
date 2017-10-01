<template>
  <ul class="details small" v-if="fields">
    <li v-if="label && title">
      <span>{{label}}</span>
      <span class="type">{{title}}</span>
    </li>
    <li v-for="field in fields">
      <span><span class="type">{{field.type}}</span> {{field.alias}}</span>
      <typed-input
        @blur="blur"
        v-on:submit="enter"
        :type="field.type"
        :name="field.alias"
        v-bind:value="params[field.alias]"
        :value="params[field.alias]"></typed-input>
      <argv-fields v-if="field.fields" :fn="fn" :params="params" :fields="field.fields"></argv-fields>
    </li>
  </ul>
</template>

<script>

import TypedInput from '@/components/TypedInput'

export default {
  name: 'argv-fields',
  components: {TypedInput},
  props: {
    fn: {
      type: Object
    },
    params: {
      type: Object
    },
    fields: {
      type: Array
    },
    label: {
      type: String
    },
    title: {
      type: String
    }
  },
  methods: {
    blur: function (e, input) {
      this.$emit('blur', e, input)
    },
    enter: function (e, input) {
      this.$emit('submit', e, input)
    }
  }
}
</script>
