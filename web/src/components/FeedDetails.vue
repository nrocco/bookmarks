<template>
  <div class="card">
    <div class="card-content">
      <div class="dropdown is-right is-pulled-right is-hoverable">
        <div class="dropdown-trigger">
          <button class="button" aria-haspopup="true" aria-controls="feed-menu">
            <span>Info</span>
            <span class="icon is-small">
              <i class="fas fa-angle-down" aria-hidden="true"></i>
            </span>
          </button>
        </div>
        <div class="dropdown-menu" id="feed-menu" role="menu">
          <div class="dropdown-content">
            <div class="dropdown-item">
              <p>Manage this feed.</p>
            </div>
            <hr class="dropdown-divider">
            <a class="dropdown-item is-primary is-outlined" @click.prevent="onRefreshFeedClicked(feed)">Refresh Feed</a>
            <a class="dropdown-item is-danger is-outlined" @click.prevent="onDeleteFeedClicked(feed)">Delete Feed</a>
          </div>
        </div>
      </div>

      <p class="title">{{ feed.Title }}</p>
      <p class="subtitle">{{ feed.URL }}</p>
      <p><i :title="feed.LastAuthored|moment('dddd, MMMM Do YYYY, HH:mm')">Last item created {{ feed.LastAuthored|moment("from") }}</i></p>
      <p><i :title="feed.Refreshed|moment('dddd, MMMM Do YYYY, HH:mm')">Last refreshed {{ feed.Refreshed|moment("from") }}</i></p>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  props: {
    feed: Object
  },
  data () {
    return {}
  },
  methods: {
    ...mapActions({
      onRefreshFeedClicked: 'refreshFeed',
      onDeleteFeedClicked: 'deleteFeed'
    })
  }
}
</script>

<style>
</style>
