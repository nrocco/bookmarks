<template>
    <div class="columns">
        <div class="column is-one-quarter">
            <div v-if="!new_feed" class="block">
                <a class="button is-primary is-outlined is-rounded is-fullwidth" @click.prevent="onAddNewFeedClicked">Add Feed</a>
            </div>
            <div v-else class="block field has-addons">
                <div class="control is-expanded">
                    <input class="input is-fullwidth" type="text" v-model="new_feed.URL" placeholder="Feed URL">
                </div>
                <div class="control">
                    <a class="button is-primary" @click.prevent="onSaveNewFeedClicked"><i class="fa fa-plus" aria-hidden="true"></i></a>
                </div>
            </div>
            <feed-panel :feeds="feeds" @selected="onFeedSelected" />
        </div>
        <div class="column">
            <div class="block" v-if="selected_feed">
                <h1 class="title">{{ selected_feed.Title }}</h1>
                <h2 class="subtitle">{{ selected_feed.URL }}</h2>
                <p>Last refreshed at: {{ selected_feed.Refreshed }}</p>
                <hr />
            </div>
            <div class="block" v-for="item in items" :key="item.ID">
                <p class="has-text-weight-bold">{{ item.Title }}</p>
                <p class="is-size-7"><a :href="item.URL">{{ item.URL }}</a></p>
                <p class="content">{{ item.Content|excerpt }}</p>
                <p class="block has-text-right">
                    <a @click.prevent="onRemoveClicked(item)" class="button is-small is-danger is-outlined">Remove</a>
                    <a @click.prevent="onReadItLaterClicked(item)" class="button is-small is-primary">Read it later</a>
                </p>
            </div>
        </div>
    </div>
</template>

<script>
import FeedPanel from './FeedPanel'

export default {
  components: {
    FeedPanel
  },
  data () {
    return {
      feeds: [],
      items: [],
      new_feed: null,
      selected_feed: null
    }
  },
  filters: {
    excerpt (value) {
      if (!value) {
        return 'No content'
      }
      return value.toString().substring(0, 500) + '...'
    }
  },
  methods: {
    onAddNewFeedClicked (event) {
      this.new_feed = {
        Category: 'default'
      }
    },
    onSaveNewFeedClicked (event) {
      this.$http.post(`/feeds`, this.new_feed).then(response => {
        this.feeds.push(response)
        this.new_feed = null
      })
    },
    onFeedSelected (feed) {
      this.selected_feed = feed
      this.$http.get(`/items`, {params: {feed: feed.ID}}).then(response => {
        this.items = response.data
      })
    },
    onReadItLaterClicked (item) {
      this.$http.post(`/items/${item.ID}/readitlater`).then(response => {
        this.items.splice(this.items.indexOf(item), 1)
      })
    },
    onRemoveClicked (item) {
      this.$http.delete(`/items/${item.ID}`).then(response => {
        this.items.splice(this.items.indexOf(item), 1)
      })
    },
    fetchFeeds () {
      var payload = {}
      if (this.selected_feed) {
        payload.feed = this.selected_feed.ID
      }

      this.$http.get(`/items`, {params: payload}).then(response => {
        this.items = response.data
      })
      this.$http.get('/feeds').then(response => {
        this.feeds = response.data
      })
    }
  },
  beforeRouteUpdate (to, from, next) {
    next()
    this.fetchFeeds()
  },
  beforeRouteEnter (to, from, next) {
    next(vm => {
      vm.fetchFeeds()
    })
  }
}
</script>

<style>
</style>
