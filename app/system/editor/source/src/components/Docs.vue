<template>
  <div class="content-main">
    <iframe class="docs" src="http://localhost:3577/docs/"></iframe>
  </div>
</template>

<script>
export default {
  name: 'docs'
}
</script>

<style scoped>
  .content-main {
    border-top: 1px solid var(--border-color);
  }

  iframe.docs {
    height: 100%;
  }
</style>
