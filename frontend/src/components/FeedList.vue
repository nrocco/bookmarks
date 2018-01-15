<template>
  <div class="columns">
    <div class="column is-one-quarter">
      <feed-panel />
    </div>
    <div class="column">
      <feed-details v-if="selectedFeed" />
      <i v-else class="has-text-centered">Select a feed...</i>
    </div>
  </div>
</template>

<script>
import FeedPanel from './FeedPanel'
import FeedDetails from './FeedDetails'

export default {
  components: {
    FeedPanel,
    FeedDetails
  },
  computed: {
    selectedFeed () {
      return this.$store.getters.selectedFeed
    }
  }
}
</script>

<style>
</style>
