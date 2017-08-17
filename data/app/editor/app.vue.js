/* globals Vue */

Vue.component('app-id', {
  template: `<div class="app-id"><a href="#switch">▾ <span class="name">{{name}}</span></a></div>`,
  data: function () {
    return {
      name: 'container / application'
    }
  },
  created: function () {
    console.log('app id created')
  }
})

Vue.component('app-sidebar', {
  template: `
        <div class="sidebar">
          <h2 class="tab">
            <a class="pages selected" data-target=".pages-list" href="#pages" title="Show pages">Pages</a>
            <a class="files" href="#files"  data-target=".files-list"title="Show files">Files</a>
            <a class="components" href="#components" data-target=".components-list" title="Show components">Components</a>
          </h2>
          <div class="scroll">
              <div class="pages-list"></div>
              <div class="files-list hidden"></div>
              <div class="components-list hidden"></div>
          </div>
        </div>
        `
})

Vue.component('app-preview', {
  template: `
        <div class="preview">
          <h2>Live Preview ~ <a class="preview-url" href="" title="Preview URL"></a></h2>
          <iframe class="live"></iframe>
        </div>
        `
})

Vue.component('app-editor', {
  template: `
          <div class="editor">
            <h2>Editor</h2>
            <div class="scroll">
              <p>Select a page or file to start editing.</p>
            </div>
          </div>
        `
})

Vue.component('editor-main', {
  template: `
          <div class="content-main">
            <div class="content">
              <app-sidebar><slot /></app-sidebar>
              <app-editor><slot /></app-editor>
              <app-preview></app-preview>
            </div>
          </div>
        `,
  created: function () {
    console.log('app created')
  }
})

let vm = new Vue({el: 'main'})
console.log(vm)
