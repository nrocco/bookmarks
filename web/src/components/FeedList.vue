<template>
  <div>
    <div class="block">
      <div class="field">
        <div class="control">
          <div class="select">
            <select v-model="selectedFeedId">
              <option :value="undefined">All</option>
              <option v-for="feed in feeds" :key="feed.ID" :value="feed.ID">{{ feed.Title }} ({{ feed.Items }})</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <feed-details v-if="selectedFeedId" :feed="selectedFeed" />

    <hr/>

    <div class="block feed-item" v-for="item in items" :key="item.ID">
      <p class="has-text-weight-bold">{{ item.Title }}</p>
      <p class="is-size-7"><a class="url" :href="item.URL">{{ item.URL }}</a></p>
      <p class="content"><i>{{ item.Date|moment("from", "now") }}</i> - {{ item.Content }}...</p>
      <p class="buttons is-right">
        <a @click.prevent="onRemoveClicked(item)" class="button is-small is-danger is-outlined">Remove</a>
        <a @click.prevent="onReadItLaterClicked(item)" class="button is-small is-primary">Read it later</a>
      </p>
    </div>
  </div>
</template>

<script>
import FeedDetails from './FeedDetails'

export default {
  components: {
    FeedDetails
  },
  data () {
    return {
      selectedFeedId: '',
      feeds: [],
      items: []
    }
  },
  computed: {
    selectedFeed () {
      return this.feeds.filter(feed => feed.ID === this.selectedFeedId).shift()
    }
  },
  methods: {
    loadFeeds () {
      this.$http.get(`/feeds`).then(response => {
        this.feeds = response.data
      })
      this.$http.get(`/items`).then(response => {
        this.items = response.data
      })
    },
    getFeedForItem (item) {
      return this.feeds.filter(feed => {
        return item.FeedID === feed.ID
      }).shift()
    },
    onReadItLaterClicked (item) {
      this.$http.post(`/items/${item.ID}/readitlater`).then(response => {
        item = response.data
      })
    },
    onRemoveClicked (item) {
      this.$http.delete(`/items/${item.ID}`).then(response => {
        this.items.splice(this.items.indexOf(item), 1)
      })
    }
  },
  watch: {
    '$route' (to, from) {
      this.selectedFeedId = to.query.selectedFeedId
      this.loadFeeds()
    }
  },
  mounted () {
    this.selectedFeedId = this.$route.query.selectedFeedId
    this.loadFeeds()
  }
}
</script>

<style>
.feed-item {
  border-bottom: 1px solid hsl(0, 0%, 96%);
  padding-bottom: 1.5rem;
}
.feed-item .url {
  word-break: break-all;
}
.content:not(:last-child) {
  margin-bottom: 0.5rem;
}
.select {
  width: 100%;
}
.select select {
  width: 100%;
}
</style>
