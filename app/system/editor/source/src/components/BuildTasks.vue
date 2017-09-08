<template>
  <div class="build-tasks">
    <div class="toggle">
      <a class="select" @click="toggle" title="Select a task">⏷ Tasks</a>
    </div>
    <div class="tasks" :class="{hidden: !show}">
      <div v-for="task, key in tasks">
        <a
          @click="select($event, key, task)"
          title="Run task">{{key}}
        </a>
        <span class="small">{{task}}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'build-tasks',
  data: function () {
    return {
      show: false
    }
  },
  computed: {
    app: function () {
      return this.$store.state.app
    },
    tasks: function () {
      return this.$store.state.app.build.tasks
    }
  },
  watch: {
    app: function (val) {
      if (val && val.build && val.build.tasks) {
        this.tasks = val.build.tasks
      }
    }
  },
  methods: {
    toggle: function (e) {
      e.stopImmediatePropagation()
      this.show = !this.show
      const hide = () => {
        document.removeEventListener('click', hide)
        this.show = false
      }

      if (this.show) {
        document.addEventListener('click', hide)
      }
    },
    select: function (e, name, task) {
      console.log(name)
      console.log(task)
    }
  }
}
</script>

<style scoped>

  .build-tasks {
    position: relative;
    flex: 1 0;
    min-width: 16rem;
  }

  .build-tasks .toggle {
    display: inline-block;
    display: flex;
  }

  .build-tasks .toggle > a:not(.select) {
    flex: 1 0;
  }

  .select {
    display: inline-block;
    background: var(--base03-color);
    padding: 0 2rem;
    flex: none;
  }

  .tasks {
    position: absolute;
    top: 3.2rem;
    left: 0;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    background: var(--base02-color);
    border: 1px solid var(--border-color);
    min-width: 10rem;
  }

  a {
    display: block;
  }

  a:hover {
    text-decoration: none;
  }

  .small {
    display: block;
    text-transform: none;
  }
</style>
