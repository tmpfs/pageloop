<template>
  <div class="notifications" :class="{hidden: !notifications.length}">
    <div class="notify" :class="{reveal: item.reveal, rendered: item.rendered, error: item.error}" v-for="item, index in notifications">
      <a class="close" @click="dismiss(item)"></a>
      <h2 v-if="item.title">{{item.title}}</h2>
      <h2 v-if="item.error">Error!</h2>
      <p v-if="item.message">{{item.message}}</p>
      <p v-if="item.error">{{item.error.toString()}}</p>
    </div>
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
      let el = this.$el.childNodes[ind]
      if (!el) {
        return
      }
      let cb = () => {
        el.removeEventListener('transitionend', cb)
        this.$store.state.notify(item, true)
      }
      el.addEventListener('transitionend', cb)
      item.rendered = false
      item.reveal = false
    }
  },
  updated: function () {
    setTimeout(() => {
      this.notifications.forEach((n) => {
        n.rendered = true
        if (!n.persist) {
          setTimeout(() => {
            this.dismiss(n)
          }, n.timeout || 5000)
        }
      })
    }, 50)
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
		left: 32rem;
		display: block;
		background: var(--base03-color);
		color: var(--color);
		font-size: 1.5rem;
		width: 32rem;
		padding: 0 0 1rem 0;
		border-radius: 3px;
		margin-bottom: 1rem;
		opacity: 0;
		transition: all 0.4s ease-in;
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

	.notify.reveal {
		opacity: 1;
		left: 0;
	}

	.notify.rendered {
		transition: none;
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
