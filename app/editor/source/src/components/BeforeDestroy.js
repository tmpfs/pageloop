export default {
  beforeDestroy: function () {
    this.$el.parentNode.removeChild(this.$el)
  }
}
