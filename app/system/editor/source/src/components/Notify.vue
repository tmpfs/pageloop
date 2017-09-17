<template>
  <div class="notifications">
    <transition-group appear name="reveal" tag="div" v-on:after-enter="afterEnter">
      <div class="notify"
        :class="{error: item.error}"
        v-bind:key="item.id"
        v-for="item, index in notifications">
          <a class="close" @click="dismiss(item)"></a>
          <h2 v-if="item.title">{{item.title}}</h2>
          <h2 v-if="item.error">Error!</h2>
          <p v-if="item.message">{{item.message}}</p>
          <p v-if="item.error">{{item.error.toString()}}</p>
      </div>
    </transition-group>
  </div>
</template>

<script>
export default {
  name: 'app-notify',
  computed: {
    notifications: function () {
      return this.$store.state.notifier.notifications
    }
  },
  methods: {
    dismiss: function (item) {
      let ind = this.notifications.indexOf(item)
      if (ind === -1) {
        return
      }
      return this.$store.state.notify(this.notifications[ind], true)
    },
    afterEnter: function (el) {
      console.log('enter cmpleted')
      console.log(el)
      // TODO: auto removal of notification
    }
  }
}
</script>

<style scoped>
	.notifications {
		position: absolute;
		z-index: 10000;
		right: 0;
		bottom: 0;
		padding: 1rem;
	}

	.notify {
		position: relative;
		display: block;
		background: var(--base03-color);
		color: var(--color);
		font-size: 1.5rem;
		width: 32rem;
		padding: 0 0 1rem 0;
		border-radius: 3px;
		margin-bottom: 1rem;
	}

  .reveal-enter {
		opacity: 0;
  }

  .reveal-enter-active, .reveal-leave-active {
		transition: all 0.4s ease-in;
    opacity: 1;
    left: 0;
  }

  .reveal-enter, .reveal-leave-to {
    opacity: 0;
		left: 32rem;
  }

	.notify.error {
	  background: var(--red-color);
    color: var(--base2-color);
	}

	.notify:last-child {
		margin-bottom: 0;
	}

	.notify p {
		margin-left: 1rem;
	}

	.notify p:last-child {
		margin-bottom: 0;
	}

	.notify > h2 {
		font-size: 1.6rem;
		padding: 0.5rem 0 0.5rem 1rem;
		border-bottom: 1px solid currentColor;
	}

	.notify > a.close {
		position: absolute;
		top: 0;
		right: 0;
	}

  .notify.error > a.close:hover {
	  background: var(--orange-color);
  }
</style>
