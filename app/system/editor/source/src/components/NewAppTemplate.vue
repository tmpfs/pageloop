<template>
  <div class="new-app-template">
    <ul class="new-app-step small">
      <li>
        <span>Choose a template</span>
        <span>Step 2/3</span>
      </li>
    </ul>
    <form @submit="nextStep">
      <div class="templates-list">
        <div
          @click="selectTemplate"
          class="app-template"
          :class="{selected: applicationTemplate === tpl.url}"
          v-for="tpl, index in templates">
          <div>
            <input
              :id="tpl.name"
              type="radio"
              :value="tpl.url"
              v-model="applicationTemplate"
              name="template" />
          </div>
          <div>
            <label :for="tpl.name">{{tpl.name}}</label>
            <p class="small">{{tpl.description}}</p>
            <!--
            <a
              class="small"
              @click="preview($event, tpl)">
                <span v-if="previewUrl !== tpl.url">Show Preview</span>
                <span v-else>Hide Preview</span>
            </a>
            -->

            <transition-group appear name="reveal">
                <iframe key="preview" v-if="previewUrl === tpl.url" @load="loaded" :src="tpl.url"></iframe>
            </transition-group>
          </div>
        </div>
      </div>

      <div class="form-actions">
        <input type="submit" value="Next: Create" class="primary" v-focus />
      </div>
    </form>
  </div>
</template>

<script>
export default {
  name: 'new-app-template',
  computed: {
    templates: function () {
      return this.$store.state.templates
    },
    previewUrl: {
      get: function () {
        return this.$store.state.newApp.previewUrl
      },
      set: function (val) {
        this.$store.state.newApp.previewUrl = val
      }
    },
    applicationTemplate: {
      get: function () {
        return this.$store.state.newApp.templateUrl
      },
      set: function (val) {
        this.$store.state.newApp.templateUrl = val
      }
    },
    template: {
      get: function () {
        return this.$store.state.newApp.template
      },
      set: function (val) {
        this.$store.state.newApp.template = val
      }
    }
  },
  mounted: function () {
    this.$store.dispatch('list-templates')
  },
  methods: {
    preview: function (tpl) {
      if (this.previewUrl !== tpl.url) {
        this.previewUrl = tpl.url
      } else {
        this.previewUrl = ''
      }
    },
    loaded: function (e) {
      // Hide scrollbars for preview iframes
      e.currentTarget.contentDocument.querySelector('body').style = 'overflow: hidden'
    },
    selectTemplate: function (e) {
      e.preventDefault()
      const radio = e.currentTarget.querySelector('input[type="radio"]')
      radio.checked = true
      this.applicationTemplate = radio.value
      const templates = this.$store.state.templates
      for (let i = 0; i < templates.length; i++) {
        if (radio.value === templates[i].url) {
          this.template = templates[i]
          break
        }
      }
      this.preview(this.template)
    },
    nextStep: function (e) {
      e.preventDefault()
      this.$store.commit('new-app-view', 'new-app-create')
    }
  }
}
</script>

<style scoped>
  input[type="radio"] {
    display: inline-block;
    line-height: 3rem;
    padding-top: 0.4rem;
    pointer-events: none;
    vertical-align: middle;
  }

  .new-app-templates {
    margin-left: 1rem;
    padding: 1rem 0;
  }

  iframe {
    width: 100%;
    height: 240px;
    pointer-events: none;
    user-select: none;
  }

  .app-template {
    display: flex;
    background: var(--base03-color);
    transition: all 0.3s ease-out;
    cursor: pointer;
    margin-bottom: 0.8rem;
    border-bottom: 2px solid transparent;
    padding: 1rem 0;
  }

  .app-template > div:first-child {
    min-width: 4rem;
    border-right: 2px solid var(--base00-color);
    text-align: center;
  }

  .app-template > div:last-child {
    flex: 1 0;
    padding: 0 1rem;
  }

  .app-template label {
    vertical-align: middle;
  }

  .app-template p {
    margin: 0;
  }

  .app-template.selected {
    border-bottom: 2px solid var(--base3-color);
  }

  .reveal-enter {
		opacity: 0;
    height: 0px;
  }

  .reveal-enter-active, .reveal-leave-active {
		transition: all 0.4s ease-in;
    opacity: 1;
    height: 240px;
  }

  .reveal-enter, .reveal-leave-to {
		opacity: 0;
    height: 0px;
  }

</style>
