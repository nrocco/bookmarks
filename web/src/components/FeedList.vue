<template>
  <div>
    <div class="block">
      <div class="field">
        <div class="control">
          <div class="select">
            <select v-model="selectedFeedId">
              <option value="">All ({{ totalUnread }})</option>
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
      <div class="tags" v-if="getFeedForItem(item).Tags">
        <span class="tag" v-for="tag in getFeedForItem(item).Tags" :key="tag">{{ tag }}</span>
      </div>
      <p class="content"><i>{{ item.Date|moment("from", "now") }}</i> - {{ item.Content|excerpt }}</p>
      <p class="buttons is-right">
        <a @click.prevent="onRemoveClicked(item)" class="button is-small is-danger is-outlined">Remove</a>
        <a @click.prevent="onReadItLaterClicked(item)" class="button is-small is-primary">Read it later</a>
      </p>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import FeedDetails from './FeedDetails'

export default {
  components: {
    FeedDetails
  },
  data () {
    return {
      selectedFeedId: ''
    }
  },
  computed: {
    feeds () {
      return this.$store.getters.feeds
    },
    selectedFeed () {
      return this.$store.getters.feeds.filter(feed => {
        return feed.ID === this.selectedFeedId
      }).shift()
    },
    items () {
      return this.$store.getters.items.filter(item => {
        return this.selectedFeedId === '' || (item.FeedID === this.selectedFeedId)
      })
    },
    totalUnread () {
      return this.$store.getters.items.length
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
    getFeedForItem (item) {
      return this.$store.getters.feeds.filter(feed => {
        return item.FeedID === feed.ID
      })[0]
    },
    ...mapActions({
      onRefreshFeedClicked: 'refreshFeed',
      onDeleteFeedClicked: 'deleteFeed',
      onReadItLaterClicked: 'readLaterFeedItem',
      onRemoveClicked: 'removeFeedItem'
    })
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
.feed-item .tags:not(:last-child) {
  margin-top: 0.5rem;
  margin-bottom: 0;
}
</style>
