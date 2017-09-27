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
          contenteditable
          @keydown="keydown"
          @keyup.enter="enter"
          class="type">{{field.type}}</span>
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
  methods: {
    keydown: function (e) {
      if (e.key === 'Enter') {
        e.preventDefault()
        e.stopImmediatePropagation()
      }
    },
    enter: function (e) {
      e.preventDefault()
      e.stopImmediatePropagation()
      console.log('enter called')
      return true
    }
  }
}
</script>

<style scoped>
  /*
  h4 {
    margin: 0;
    padding: 0 0 0.5rem 0;
    border-bottom: 1px solid var(--border-color);
    display: inline-block;
    font-size: 1.4rem;
  }
  */
</style>
