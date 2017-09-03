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
  .data-editor {
    font-size: 1.4rem;
    user-select: none;
    padding-right: 1rem;
  }

  .data-editor ul, .data-editor ol {
    list-style-type: none;
    margin-left: 1rem;
    padding: 0;
  }

  .data-key::after {
    content: ':';
    margin-right: 1rem;
  }

  .data-key {
    display: inline-block;
    min-width: 6rem;
  }

  .data-key, .data-value {
    padding: 0.1rem 0.3rem;
  }

  .data-value {
    background: var(--base03-color);
    color: var(--base2-color);
    border-radius: 2px;
  }
</style>
