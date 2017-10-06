<template>
  <span class="typed-input" :class="{error: error}">
    <input
      v-if="isBoolean()"
      type="checkbox"
      @blur="blur"
      @keyup.enter="enter"
      :name="name"
      v-model="target[name]"
      :checked="booleanValue" />
    <input
      v-else-if="isNumber()"
      type="number"
      class="input"
      @blur="blur"
      @keyup.enter="enter"
      v-model.number="target[name]"
      :name="name" />
    <input
      v-else-if="isString()"
      type="text"
      :name="name"
      @blur="blur"
      @keyup.enter="enter"
      v-model="target[name]"
      class="input" />
    <span v-else class="input unsupported">{{type}}</span>
  </span>
</template>

<script>
export default {
  name: 'typed-input',
  data: function () {
    return {
      error: false,
      booleanValue: false,
      numberValue: 0,
      stringValue: ''
    }
  },
  props: {
    target: {
      type: Object
    },
    type: {
      type: String
    },
    name: {
      type: String
    },
    value: {
      type: String
    },
    placeholder: {
      type: String
    }
  },
  methods: {
    isBoolean: function () {
      return /^bool/i.test(this.type)
    },
    isString: function () {
      return this.type === 'string'
    },
    isNumber: function () {
      return /(int|float)/i.test(this.type)
    },
    toggleBoolean: function () {
      this.booleanValue = !this.booleanValue
    },
    blur: function (e) {
      this.error = this.hasError()
    },
    hasError: function () {
      if (this.isString() && this.target[this.name] === '') {
        return true
      }
      return false
    },
    enter: function (e) {
      e.preventDefault()
      e.stopImmediatePropagation()
      console.log('enter listener called')
      this.error = this.hasError()
      if (this.error) {
        return true
      }
      this.$emit('submit', e, this)
    }
  }
}
</script>

<style scoped>

  .typed-input {
    display: inline-block;
    padding-left: 0 !important;
  }

  .typed-input.error > input[type="text"] {
    border-color: var(--red-color);
  }

  .input {
    font-size: 1.4rem;
    padding: 0;
    min-height: 1.8rem;
    border-radius: 0;
    display: inline-block;
    color: var(--cyan-color);
    user-select: none;
    cursor: default;
    padding-bottom: 0.5rem;
    margin: 0 1rem 0.5rem 0;
  }

  .input[data-type="bool"] {
    cursor: pointer;
  }

  input[type="checkbox"] {
    margin: 0;
  }

  input[type="number"] {
    margin: 0 0 0 1rem;
  }

  input[type="text"] {
    user-select: auto;
    cursor: auto;
    border-bottom: 1px solid var(--base00-color);
  }

  .unsupported {
    color: var(--red-color);
  }
</style>
