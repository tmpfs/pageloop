export default {
  computed: {
    currentFile: function () {
      return this.$store.state.app.current
    }
  },
  methods: {
    click: function (e, item) {
      let target = e.currentTarget
      let i, n, j, start
      let urls = []

      const lookAhead = () => {
        // Look ahead for last selection
        for (j = this.$el.childNodes.length - 1; j >= i; j--) {
          n = this.$el.childNodes[j]
          if (n.nodeType === 1) {
            if (!start && ~n.className.indexOf('selected')) {
              start = n
            }
            if (start) {
              urls.push(n.getAttribute('data-url'))
            }
          }
        }
      }

      const lookBehind = () => {
        for (i = 0; i < this.$el.childNodes.length; i++) {
          n = this.$el.childNodes[i]
          if (n.nodeType === 1) {
            if (!start && ~n.className.indexOf('selected')) {
              start = n
            }
            if (start) {
              urls.push(n.getAttribute('data-url'))
            }
            if (target === n) {
              if (start) {
                break
              } else {
                lookAhead()
                break
              }
            }
          }
        }
      }

      if (e.ctrlKey) {
        this.selection.push(item)
      } else if (e.shiftKey && this.selection.length) {
        lookBehind()
        this.selection = []
        urls.forEach((u) => {
          this.selection.push(this.getSelectionByUrl(u))
        })
      } else {
        this.selection = [item]
        return this.go(item)
      }
    }
  }
}
