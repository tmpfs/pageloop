<template>
  <p class="small">
    <input id="frontmatter-display" type="checkbox" v-model="frontmatterDisplay" />
    <label for="frontmatter-display">Show frontmatter in code view</label>
  </p>
</template>

<script>
export default {
  name: 'editor-preferences',
  data: function () {
    return {
      frontmatterDisplay: true
    }
  }
}
</script>

<style scoped>
</style>
