<template>
  <div class="block">
    <h1 class="title">{{ feed.Title }}</h1>
    <h2 class="subtitle">{{ feed.URL }}</h2>
    <div class="block">
      <span :title="feed.LastAuthored|moment('dddd, MMMM Do YYYY, HH:mm')"><i>Last item created at:</i> {{ feed.LastAuthored|moment("from") }}</span>
    </div>
    <div class="block">
      <span :title="feed.Refreshed|moment('dddd, MMMM Do YYYY, HH:mm')"><i>Last refreshed at:</i> {{ feed.Refreshed|moment("from") }}</span>
      <a class="button is-small is-primary is-outlined" @click.prevent="onRefreshFeedClicked(feed)">Refresh</a>
      <a class="button is-small is-danger is-outlined" @click.prevent="onDeleteFeedClicked(feed)">Delete</a>
    </div>
    <hr />
    <div v-if="items.length > 0">
      <div v-for="item in items" :key="item.ID">
        <p class="has-text-weight-bold">{{ item.Title }}</p>
        <p class="is-size-7"><a :href="item.URL">{{ item.URL }}</a></p>
        <p class="content">{{ item.Content|excerpt }}</p>
        <p class="block has-text-right">
          <a @click.prevent="onRemoveClicked(item)" class="button is-small is-danger is-outlined">Remove</a>
          <a @click.prevent="onReadItLaterClicked(item)" class="button is-small is-primary">Read it later</a>
        </p>
      </div>
    </div>
    <div v-else class="block has-text-centered">
      <i>No items in this feed</i>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  computed: {
    feed () {
      return this.$store.getters.selectedFeed
    },
    items () {
      return this.$store.getters.items
    }
  },
  methods: {
    ...mapActions({
      onRefreshFeedClicked: 'refreshFeed',
      onDeleteFeedClicked: 'deleteFeed',
      onReadItLaterClicked: 'readLaterFeedItem',
      onRemoveClicked: 'deleteFeedItem'
    })
  },
  filters: {
    excerpt (value) {
      if (!value) {
        return 'No content'
      }
      return value.toString().substring(0, 500) + '...'
    }
  }
}
</script>

<style>
</style>
