<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'app',
  watch: {
    '$route' (to, from) {
      this.$store.dispatch(to.name, to)
    }
  },
  mounted () {
    this.$store.dispatch(this.$route.name, this.$route)
  }
}
</script>

<style lang="css">
</style>
