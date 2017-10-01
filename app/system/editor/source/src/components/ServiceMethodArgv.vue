<template>
  <div v-if="fn.fields">
    <argv-fields label="Argument Fields" :title="fn.arg" :fn="fn" :params="params" :fields="fn.fields"></argv-fields>
  </div>
</template>

<script>

import ArgvFields from '@/components/ServiceArgvFields'

export default {
  name: 'method-argv',
  components: {ArgvFields},
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
