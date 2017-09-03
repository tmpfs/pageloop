<script>
export default {
  name: 'data-editor',
  computed: {
    pageDataJson: function () {
      return JSON.stringify(this.pageData, undefined, 2)
    },
    pageData: {
      get: function () {
        if (!this.$store.state.current || !this.$store.state.current.page) {
          return {}
        }
        return this.$store.state.current.page.data
      },
      set: function (val) {
        //
      }
    }
  },
  render: function (h) {
    // We need recursion to render meta page data
    function list (target) {
      let isArr = Array.isArray(target)
      function it (o, fn) {
        if (Array.isArray(o)) {
          o.forEach(fn)
        } else {
          for (let k in o) {
            fn(o[k], k)
          }
        }
      }

      let el = h(isArr ? 'ol' : 'ul', null, [])
      it(target, (value, key) => {
        let li = h('li', {'data-type': typeof (value), 'data-key': key}, [])
        if (!isArr) {
          let k = h('span', {class: 'data-key'}, ['' + key])
          li.children.push(k)
        }
        if (Array.isArray(value) || value && typeof (value) === 'object') {
          let nodes = list(value)
          li.children.push(...nodes)
        } else {
          let v = h('span', {class: 'data-value'}, ['' + value])
          li.children.push(v)
        }
        el.children.push(li)
      })
      return [el]
    }

    let children = list(this.pageData)
    let el = h('div', {class: 'data-editor'}, children)
    return el
  }
}
</script>

<style scoped>
</style>
