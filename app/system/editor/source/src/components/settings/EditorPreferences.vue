<template>
  <div class="small editor-preferences">
    <div>
      <input id="show-data-editor" type="checkbox" v-model="showDataEditor" />
      <label for="show-data-editor">Show data editor</label>
    </div>
    <div>
      <input id="show-text-editor" type="checkbox" v-model="showTextEditor" />
      <label for="show-text-editor">Show text editor</label>
    </div>
    <div>
      <input id="use-vim-mode" type="checkbox" v-model="useVimCodeEditor" />
      <label for="use-vim-mode">Use VIM text editor mode</label>
    </div>
    <div>
      <input id="frontmatter-display" type="checkbox" v-model="frontmatterDisplay" />
      <label for="frontmatter-display">Show frontmatter in text editor</label>
    </div>
  </div>
</template>

<script>
export default {
  name: 'editor-preferences',
  data: function () {
    return {
      showDataEditor: true,
      showTextEditor: true,
      useVimCodeEditor: true,
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
