<script>
export default {
  name: 'data-editor',
  computed: {
    pageData: function () {
      const state = this.$store.state
      if (state.current && state.current.data) {
        return this.$store.state.current.data
      }

      if (state.current.page) {
        return state.current.page.data
      }

      return {}
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
          let v = h('span', {class: 'data-value'}, [])
          v.children.push(h('span', {class: 'data-input'}, ['' + value]))
          li.children.push(v)
        }
        el.children.push(li)
      })
      return [el]
    }

    let children = list(this.pageData)
    let el = h('div', {class: 'data-editor scroll'}, children)
    return el
  }
}
</script>

<style scoped>
  .data-editor {
    font-size: 1.5rem;
    user-select: none;
  }

  .data-editor ul, .data-editor ol {
    display: table;
    list-style-type: none;
    padding-left: 1rem;
    padding: 0;
    margin: 0;
    width: 100%;
  }

  .data-editor ul li, .data-editor ol li {
    display: table-row;
  }

  .data-editor li:hover {
    background: var(--base03-color);
  }

  .data-key, .data-value {
    display: table-cell;
    padding: 0.2rem 0 0.2rem 0.2rem;
  }

  .data-key {
    width: 25%;
    text-shadow: 1px 1px 1px var(--base03-color);
    clear: both;
    padding-left: 1rem;
  }

  .data-key::after {
    content: ':';
    float: right;
    margin-right: 0.5rem;
  }

  .data-value {
    width: 75%;
  }

  .data-value > .data-input {
    display: inline-block;
    background: var(--base03-color);
    border-radius: 0.3rem;
    color: var(--base1-color);
    vertical-align: top;
  }

  .data-input::before, .data-input::after {
    content: '';
    margin-right: 0.5rem;
  }
</style>
