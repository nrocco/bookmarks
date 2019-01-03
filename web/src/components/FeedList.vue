<template>
  <div>
    <div class="block">
      <div class="field has-addons has-addons-right">
        <p class="control">
          <span class="select">
            <select>
              <option value="">All ({{ totalUnread }})</option>
              <option v-for="feed in feeds" :key="feed.ID" :value="feed.ID">{{ feed.Title }} ({{ feed.Items }})</option>
            </select>
          </span>
        </p>
        <p class="control is-expanded">
          <input class="input" type="search" placeholder="Search" autofocus>
        </p>
      </div>
    </div>

    <hr/>

    <div class="block feed-item" v-for="item in items" :key="item.ID">
      <p class="has-text-weight-bold">{{ item.Title }}</p>
      <p class="is-size-7"><a class="url" :href="item.URL">{{ item.URL }}</a></p>
      <p class="content">{{ item.Content|excerpt }}</p>
      <p class="buttons is-right">
        <a @click.prevent="onRemoveClicked(item)" class="button is-small is-danger is-outlined">Remove</a>
        <a @click.prevent="onReadItLaterClicked(item)" class="button is-small is-primary">Read it later</a>
      </p>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  data () {
    return {}
  },
  computed: {
    feeds () {
      return this.$store.getters.feeds
    },
    items () {
      return this.$store.getters.items
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
    ...mapActions({
      onRefreshFeedClicked: 'refreshFeed',
      onDeleteFeedClicked: 'deleteFeed',
      onReadItLaterClicked: 'readLaterFeedItem',
      onRemoveClicked: 'removeFeedItem'
    })
  },
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
</style>
