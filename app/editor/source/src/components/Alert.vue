<template>
  <div class="alert" :class="{hidden: !alert.visible}">
    <div class="background"></div>
    <div class="dialog">
      <a class="close" @click="dismiss"></a>
      <h2>{{alert.title}}</h2>
      <div class="dialog-panel">
        <p v-if="alert.message">{{alert.message}}</p>
        <small v-if="alert.note">{{alert.note}}</small>
        <div class="form-actions">
          <button class="sml" @click="dismiss">Cancel</button>
          <button class="sml primary" @click="ok">OK</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app-alert',
  computed: {
    alert: function () {
      return this.$store.state.alert
    }
  },
  methods: {
    dismiss: function () {
      this.$store.commit('alert-hide')
    },
    ok: function () {
      this.alert.ok()
      this.dismiss()
    }
  }
}
</script>

<style scoped>
	.alert {
		position: absolute;
		left: 0;
		top: 0;
		right: 0;
		bottom: 0;
		z-index: 10001;
		text-align: center;
	}

	.alert > .background {
		position: absolute;
		z-index: 1;
		left: 0;
		top: 0;
		right: 0;
		bottom: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.6);
		pointer-events: none;
	}

	.alert > .dialog {
		position: relative;
		z-index: 2;
		display: inline-block;
		text-align: left;
		margin: 0 auto;
		background: var(--background);
		color: var(--color);
		font-size: 1.5rem;
		min-width: 32rem;
		padding: 0 0 1rem 0;
		border-bottom-right-radius: 3px;
		border-bottom-left-radius: 3px;
	}

	.alert .dialog-panel {
		padding: 0 1rem;
	}

	.alert > .dialog > h2, .notify > h2 {
		font-size: 1.6rem;
		padding: 0.5rem 0 0.5rem 1rem;
		border-bottom: 1px solid currentColor;
	}

	.alert > .dialog > a.close, .notify > a.close {
		position: absolute;
		top: 0;
		right: 0;
	}
</style>
