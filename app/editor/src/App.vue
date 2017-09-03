<template>
  <main>
    <app-alert></app-alert>
    <app-notify></app-notify>
    <app-header></app-header>
    <app-main></app-main>
    <app-footer></app-footer>
  </main>
</template>

<script>
import AppAlert from '@/components/Alert'
import AppNotify from '@/components/Notify'
import AppHeader from '@/components/Header'
import AppMain from '@/components/Main'
import AppFooter from '@/components/Footer'

export default {
  name: 'app',
  components: {AppAlert, AppNotify, AppHeader, AppMain, AppFooter}
}
</script>

<style>
#app {
  /*
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
  */
}
</style>
