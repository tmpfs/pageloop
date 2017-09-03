<template>
  <div class="content-main">
    <div class="content">
      <sidebar></sidebar>
      <editor></editor>
      <preview></preview>
    </div>
  </div>
</template>

<script>

import Sidebar from '@/components/Sidebar'
import Editor from '@/components/Editor'
import Preview from '@/components/Preview'

export default {
  name: 'edit',
  components: {Sidebar, Editor, Preview}
}
</script>

<style scoped>
</style>
