<template>
  <div class="small editor-preferences">
    <div>
      <input id="show-data-editor" type="checkbox" v-model="showDataEditor" />
      <label for="show-data-editor">Show data editor</label>
    </div>
    <div>
      <input id="show-code-editor" type="checkbox" v-model="showCodeEditor" />
      <label for="show-code-editor">Show code editor</label>
    </div>
    <div>
      <input id="frontmatter-display" type="checkbox" v-model="frontmatterDisplay" />
      <label for="frontmatter-display">Show frontmatter in code view</label>
    </div>
  </div>
</template>

<script>
export default {
  name: 'editor-preferences',
  data: function () {
    return {
      showDataEditor: true,
      showCodeEditor: true,
      frontmatterDisplay: true
    }
  }
}
</script>

<style scoped>
  .editor-preferences {
    margin-top: 1rem;
  }
</style>
