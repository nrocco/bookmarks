<template>
    <div class="columns">
        <div class="column is-one-quarter">
            <feed-panel :feeds="feeds" @selected="onFeedSelected" @added="onFeedAdded" />
        </div>
        <div class="column" v-if="selected_feed">
            <feed-details :feed="selected_feed" @deleted="onFeedDeleted" />
        </div>
        <div class="column" v-else>
            <i class="has-text-centered">Select a feed from the right</i>
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
  data () {
    return {
      feeds: [],
      selected_feed: null
    }
  },
  methods: {
    onFeedAdded (feed) {
      this.feeds.push(feed)
    },
    onFeedDeleted (feed) {
      this.selected_feed = null
      this.feeds.splice(this.feeds.indexOf(feed), 1)
    },
    onFeedSelected (feed) {
      this.selected_feed = feed
      this.$http.get(`/items`, {params: {feed: feed.ID}}).then(response => {
        this.items = response.data
      })
    },
    load (filters = null) {
      // TODO handle filters

      var payload = {}

      if (this.selected_feed) {
        payload.feed = this.selected_feed.ID
      }

      this.$http.get('/feeds').then(response => {
        this.feeds = response.data
      })
    }
  },
  beforeRouteUpdate (to, from, next) {
    next()
    this.load(to.query)
  },
  beforeRouteEnter (to, from, next) {
    next(vm => {
      vm.load(to.query)
    })
  }
}
</script>

<style>
</style>
