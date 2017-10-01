<template>
  <div class="argv">
    <p class="small title" v-if="label && title">
      <span>{{label}}</span>
      <span class="type">{{title}}</span>
    </p>
    <ul class="details small" v-if="fields">
      <li v-for="field in fields">
        <span v-if="!field.fields"><span class="type">{{field.type}}</span> {{field.alias}}</span>
        <typed-input
          v-if="!field.fields"
          @blur="blur"
          v-on:submit="enter"
          :type="field.type"
          :name="field.alias"
          v-bind:value="params[field.alias]"
          :value="params[field.alias]"></typed-input>
          <argv-fields :label="field.alias" :title="field.type" v-if="field.fields" :fn="fn" :params="params" :fields="field.fields"></argv-fields>
      </li>
    </ul>
  </div>
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

<style scoped>
  .argv {
    flex: 1 0;
  }

  p.title {
    display: flex;
    border-bottom: 1px solid var(--border-color);
  }

  p.title > span {
    flex: 1 0;
  }

  p.title > span.type {
    text-align: right;
  }
</style>
