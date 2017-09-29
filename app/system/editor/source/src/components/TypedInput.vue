<template>
  <span
    :data-type="type"
    :data-name="name"
    contenteditable
    @focus="focus"
    @blur="blur"
    @keyup="keyup"
    @keydown="keydown"
    @keyup.enter="enter"
    v-model="value"
    class="field input">{{value}}</span>
</template>

<script>
export default {
  name: 'typed-input',
  data: function () {
    return {
      text: ''
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
    }
  },
  methods: {
    getText: function () {
      return this.text
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
    keyup: function (e) {
      this.text = e.currentTarget.innerText
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
      this.$emit('submit', e, this)
    }
  }
}
</script>

<style scoped>
  .input {
    color: var(--cyan-color);
  }
</style>
