<template>
  <div class="build-tasks">
    <div class="toggle">
      <a class="select" @click="toggle" title="Select a task">⏷ Tasks</a>
    </div>
    <div class="tasks" :class="{hidden: !show}">
      <a v-for="task, key in tasks"
        @click="select($event, key, task)"
        title="Run task">{{key}}</a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'build-tasks',
  data: function () {
    return {
      show: false,
      filter: 'media',
      label: 'media'
    }
  },
  computed: {
    tasks: function () {
      return this.$store.state.app.build.tasks
    }
  },
  mounted: function () {
    console.log(this.tasks)
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
  }

  .build-tasks .toggle {
    display: inline-block;
    display: flex;
  }

  .build-tasks .toggle > a:not(.select) {
    flex: 1 0;
  }

  .tasks {
    position: absolute;
    top: 2.1rem;
    right: -2.2rem;
    display: flex;
    flex-direction: column;
    background: var(--base02-color);
    border: 1px solid var(--border-color);
    min-width: 10rem;
  }

  .select {
    display: inline-block;
    flex: none;
    background: var(--base02-color);
  }
</style>
