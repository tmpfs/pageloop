<template>
  <span class="typed-input">
    <input
      v-if="isBoolean()"
      type="checkbox"
      @keyup.enter="enter"
      :name="name"
      :checked="booleanValue" />
    <input
      v-else-if="isNumber()"
      type="number"
      class="input"
      @keyup.enter="enter"
      v-model="numberValue"
      :name="name" />
    <input
      v-else-if="isString()"
      type="text"
      :data-type="type"
      :data-name="name"
      :contenteditable="!isBoolean()"
      @keyup.enter="enter"
      :value="getDefaultValue(value)"
      class="input" />
    <span v-else class="input unsupported">{{type}}</span>
  </span>
</template>

<script>
export default {
  name: 'typed-input',
  data: function () {
    return {
      text: '',
      booleanValue: false,
      numberValue: 0
    }
  },
  props: {
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
    getDefaultValue: function (val) {
      if (val === '') {
        if (this.isBoolean()) {
          return false
        }
      }
      return val
    },
    isBoolean: function () {
      return /^bool/i.test(this.type)
    },
    isString: function () {
      return this.type === 'string'
    },
    isNumber: function () {
      return /(int|float)/i.test(this.type)
    },
    getText: function () {
      return this.text || this.value
    },
    toggleBoolean: function () {
      this.booleanValue = !this.booleanValue
    },
    getValue: function () {
      // const alias = this.name
      const type = this.type
      let value = this.getText()

      // Strings are passed through verbatim
      if (type !== 'string') {
        // Quick type conversion for numbers, booleans, json arrays/objects and null
        const doc = `{"value": ${value}}`
        let result
        try {
          result = JSON.parse(doc)
        } catch (e) {
          // Can and will fail
        }

        // Coercion succeeded
        if (result) {
          value = result.value
        }
      }
      return value
    },
    focus: function (e) {
      const sel = getSelection()
      sel.selectAllChildren(e.target)
    },
    blur: function (e) {
      const sel = getSelection()
      sel.removeAllRanges()
      this.$emit('blur', e, this)
    },
    enter: function (e) {
      e.preventDefault()
      e.stopImmediatePropagation()
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

  .input {
    font-size: 1.4rem;
    padding: 0;
    min-height: 1.8rem;
    border-radius: 0;
    display: inline-block;
    width: 100%;
    color: var(--cyan-color);
    user-select: none;
    cursor: default;
    padding-bottom: 0.5rem;
    margin: 0 0 0.5rem 1rem;
  }

  .input[data-type="bool"] {
    cursor: pointer;
  }

  input[type="checkbox"] {
    margin: 0 0 0 1rem;
  }

  input[type="number"] {
    margin: 0 0 0 1rem;
  }

  .input[contenteditable="true"] {
    user-select: auto;
    cursor: auto;
    border-bottom: 1px solid var(--base00-color);
  }

  .unsupported {
    color: var(--red-color);
  }
</style>
