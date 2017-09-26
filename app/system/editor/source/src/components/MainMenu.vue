<template>
  <transition name="reveal">
    <div class="main-menu scroll" v-if="show">
      <nav>
        <a
          tabindex="1"
          @click="$store.dispatch('navigate', {href: 'home'})"
          class="home"
          :class="{selected: selectedView === 'home'}"
          title="Home page">
          <i class="fa fa-home"></i>
          Home
        </a>
        <div>
          <a
            tabindex="2"
            @click="$store.dispatch('navigate', {href: 'apps'})"
            :class="{selected: selectedView === 'apps'}"
            title="View and edit applications">
            <i class="fa fa-cube"></i>
            Applications
          </a>
          <ul class="apps" v-if="apps.length">
            <li v-for="app in apps">
              <a
                :class="{selected: isSelected(app)}"
                :title="app.description"
                @click="editApplication(app)">
              {{app.display}}
              </a>
            </li>
          </ul>
        </div>
        <a
          tabindex="3"
          @click="$store.dispatch('navigate', {href: 'docs'})"
          :class="{selected: selectedView === 'docs'}"
          title="Documentation">
          <i class="fa fa-book"></i>
          Documentation
        </a>
        <a
          tabindex="4"
          @click="$store.dispatch('navigate', {href: 'settings'})"
          :class="{selected: selectedView === 'settings'}"
          title="Settings">
          <i class="fa fa-cog"></i>
          Settings
        </a>
      </nav>
    </div>
  </transition>
</template>

<script>
export default {
  name: 'main-menu',
  created: function () {
    return this.$store.dispatch('containers')
  },
  computed: {
    selectedView: function () {
      return this.$store.state.main.view
    },
    show: function () {
      return this.$store.state.showMainMenu
    },
    apps: function () {
      return this.$store.state.apps
    }
  },
  methods: {
    isSelected: function (app) {
      const state = this.$store.state
      return app.container === state.container && app.name === state.application && this.selectedView === 'edit'
    },
    getContainer: function (app) {
      return this.$store.state.getContainerByName(app.container)
    },
    editApplication: function (app) {
      const container = this.getContainer(app)
      return this.$store.dispatch('edit-app', {container: container, application: app})
    }
  }
}
</script>

<style scoped>

  .scroll {
    height: 100%;
  }

  ul {
    padding: 0;
    margin: 0;
    list-style-type: none;
  }

  .apps {
    margin-left: 1rem;
  }

  .main-menu {
    width: 16%;
    border-top: 1px solid var(--border-color);
    border-right: 1px solid var(--border-color);
    font-size: 1.5rem;
    background: var(--base03-color);
    user-select: none;
  }

  .main-menu > nav {
    padding: 1rem;
  }

  .main-menu > nav a {
    display: block;
    transition: all 0.4s ease-out;
    color: var(--base00-color);
    margin-bottom: 0.5rem;
    border-bottom: 1px solid transparent;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .main-menu > nav a:hover {
    color: var(--base2-color);
  }

  .main-menu > nav a.selected {
    color: var(--base2-color);
    border-bottom: 1px solid var(--base2-color);
  }

  .reveal-enter-active, .reveal-leave-active {
		transition: all 0.4s ease-out;
    opacity: 1;
    width: 16%;
  }

  .reveal-enter, .reveal-leave-to {
		opacity: 0;
    width: 0%;
  }

</style>
