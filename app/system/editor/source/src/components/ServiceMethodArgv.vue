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
          @focus="focus"
          @blur="blur"
          @keydown="keydown"
          @keyup.enter="enter"
          class="field input">{{params[field.alias]}}</span>
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
  /*
  mounted: function () {
    console.log(arguments)
    console.log('argv mounted')
  },
  */
  methods: {
    focus: function (e) {
      const sel = getSelection()
      sel.selectAllChildren(e.target)
    },
    blur: function (e) {
      const sel = getSelection()
      sel.removeAllRanges()
      this.setParam(e.currentTarget)
    },
    /*
    keyup: function (e) {
      const el = e.currentTarget
      const alias = el.getAttribute('data-alias')
      let value = el.innerText

      // TODO: type coercion
      const type = el.getAttribute('data-type')

      // Strings are passed through verbatim
      if (type !== 'string') {
        // Quick type conversion for numbers, booleans and null
        const doc = `{"value": ${value}}`
        let result
        try {
          result = JSON.parse(doc)
        } catch (e) {
          // Can and will fail
          console.error(e)
        }

        // Coercion succeeded
        if (result) {
          value = result.value
          console.log('got json parse value: ' + value)
        }
      }

      console.log('alias: ' + alias)
      console.log('field type: ' + type)
      console.log('field value: ' + value)
      console.log('value type: ' + typeof (value))

      this.params[alias] = value
    },
    */
    keydown: function (e) {
      if (e.key === 'Enter') {
        e.preventDefault()
        e.stopImmediatePropagation()
      }
    },
    setParam: function (el) {
      const alias = el.getAttribute('data-alias')
      let value = el.innerText
      const type = el.getAttribute('data-type')

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
          // console.log('got json parse value: ' + value)
        }
      }

      console.log('alias: ' + alias)
      console.log('field type: ' + type)
      console.log('field value: ' + value)
      console.log('value type: ' + typeof (value))

      this.params[alias] = value
    },
    update: function () {
      const fields = this.$el.querySelectorAll('.field')
      fields.forEach((el) => {
        this.setParam(el)
      })
    },
    enter: function (e) {
      e.preventDefault()
      e.stopImmediatePropagation()
      this.update()
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
