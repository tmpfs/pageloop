<template>
  <div class="notifications" :class="{hidden: !notifications.length}">
    <div class="notify" :class="{reveal: item.reveal, rendered: item.rendered}" v-for="item, index in notifications">
      <a class="close" @click="dismiss(item)"></a>
      <h2 v-if="item.title">{{item.title}}</h2>
      <p>{{item.message}}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app-notify',
  computed: {
    notifications: function () {
      return this.$store.state.notifications
    }
  },
  methods: {
    dismiss: function (item) {
      let ind = this.notifications.indexOf(item)
      let el = this.$el.childNodes[ind]
      if (!el) {
        return
      }
      let cb = () => {
        el.removeEventListener('transitionend', cb)
        this.$store.state.notify(item, true)
      }
      el.addEventListener('transitionend', cb)
      item.rendered = false
      item.reveal = false
    }
  },
  updated: function () {
    setTimeout(() => {
      this.notifications.forEach((n) => {
        n.rendered = true
        if (!n.persist) {
          setTimeout(() => {
            this.dismiss(n)
          }, n.timeout || 5000)
        }
      })
    }, 50)
  }
}
</script>

<style scoped>

</style>
