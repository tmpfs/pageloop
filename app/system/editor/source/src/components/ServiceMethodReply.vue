<template>
  <div v-if="reply">
    <h4>Reply</h4>
    <ul class="details small">
      <li>
        <span>Status</span>
        <span>{{reply.response.status}}</span>
      </li>
      <li>
        <span>Duration</span>
        <span>{{reply.response.duration}}ms</span>
      </li>
      <li>
        <span>Request Method</span>
        <span>{{reply.response.method}}</span>
      </li>
      <li>
        <span>URL</span>
        <span>{{reply.response.url}}</span>
      </li>
      <li>
        <span>Transport</span>
        <span>{{reply.response.transport}}</span>
      </li>
    </ul>
    <h4>Document</h4>
    <pre class="small">{{JSON.stringify(reply.document, undefined, 2)}}</pre>
  </div>
</template>

<script>
export default {
  name: 'method-reply',
  props: {
    reply: {
      type: Object
    }
  }
}
</script>

<style scoped>

  h4 {
    margin: 0;
    padding: 0 0 0.5rem 0;
    border-bottom: 1px solid var(--border-color);
    display: inline-block;
    font-size: 1.4rem;
  }

</style>
