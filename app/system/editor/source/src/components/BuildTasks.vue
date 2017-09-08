<template>
  <div class="build-tasks">
    <div class="toggle">
      <a class="select" @click="toggle" title="Select a task">⏷ Tasks</a>
    </div>
    <div class="tasks" :class="{hidden: !show}">
      <div
        title="Run task"
        class="task"
        v-for="task, key in tasks"
        @click="select($event, key, task)">
        <span class="name">{{key}}</span>
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

  .select {
    display: inline-block;
    background: var(--base03-color);
  }

  .tasks {
    position: absolute;
    top: 2.8rem;
    left: 0;
    background: var(--base03-color);
    border: 1px solid var(--border-color);
    width: 24rem;
  }

  .task {
    padding-left: 2rem;
    cursor: pointer;
  }

  .task:hover {
    color: var(--base2-color);
  }

  .tasks a {
    padding: 0;
  }

  a:hover {
    text-decoration: none;
  }

  .name {
    font-size: 1.4rem;
    text-transform: uppercase;
  }

  .small {
    display: block;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }
</style>
