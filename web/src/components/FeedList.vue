<template>
  <div>
    <div class="block">
      <div class="field has-addons">
        <div class="control is-expanded">
          <div v-if="newFeed==null" class="select">
            <select v-model="filters.feed" @change="onFilterChange">
              <option :value="undefined">All</option>
              <option v-for="feed in feeds" :key="feed.ID" :value="feed.ID">{{ feed.Title }} ({{ feed.Items }})</option>
            </select>
          </div>
          <div v-else>
            <input class="input" type="text" v-model="newFeed" placeholder="New atom/rss feed url">
          </div>
        </div>

        <div class="control" v-if="!filters.feed">
          <a class="button" :class="{'is-info':newFeed!=='', 'is-warning': newFeed===''}" @click="onAddFeedClicked">{{ newFeed!=='' ? 'New Feed' : 'Cancel' }}</a>
        </div>

        <div v-if="selectedFeed" class="dropdown is-right is-pulled-right is-hoverable">
          <div class="dropdown-trigger">
            <button class="button is-info" aria-haspopup="true" aria-controls="feed-menu">
              <span>Manage</span>
              <span class="icon is-small">
                <i class="fas fa-angle-down" aria-hidden="true"></i>
              </span>
            </button>
          </div>
          <div class="dropdown-menu" id="feed-menu" role="menu">
            <div class="dropdown-content">
              <div class="dropdown-item">
                <p><a class="url" :href="selectedFeed.URL" :target="isIphone ? '_blank' : ''">{{ selectedFeed.URL }}</a></p>
              </div>
              <div class="dropdown-item">
                <p>Last item created: <i :title="selectedFeed.LastAuthored|moment('dddd, MMMM Do YYYY, HH:mm')">{{ selectedFeed.LastAuthored|moment("from") }}</i></p>
              </div>
              <div class="dropdown-item">
                <p>Last refreshed: <i :title="selectedFeed.Refreshed|moment('dddd, MMMM Do YYYY, HH:mm')">{{ selectedFeed.Refreshed|moment("from") }}</i></p>
              </div>
              <hr class="dropdown-divider">
              <a class="dropdown-item is-primary is-outlined" @click.prevent="onRefreshFeedClicked(selectedFeed)">Refresh Feed</a>
              <a class="dropdown-item is-danger is-outlined" @click.prevent="onDeleteFeedClicked(selectedFeed)">Delete Feed</a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <hr/>

    <div class="block feed-item" v-for="item in items" :key="item.ID">
      <p class="has-text-weight-bold">{{ item.Title }}</p>
      <p class="is-size-7"><a class="url" :href="item.URL" :target="isIphone ? '_blank' : ''">{{ item.URL }}</a></p>
      <p class="content"><i>{{ item.Date|moment("from", "now") }}</i> - {{ item.Content }}&#8230;</p>
      <p class="buttons">
      <div class="field is-grouped">
        <p class="control"><a @click.prevent="onRemoveClicked(item)" class="button is-small is-danger is-outlined">Remove</a></p>
        <span class="is-button-pusher" />
        <p class="control"><a @click.prevent="onReadItLaterClicked(item)" class="button is-small is-primary">Read it later</a></p>
      </div>
    </div>
  </div>
</template>

<script>
import LoaderMixin from '../mixins/loader.js'

export default {
  mixins: [
    LoaderMixin
  ],

  data: () => ({
    newFeed: null,
    filters: {},
    feeds: [],
    items: []
  }),

  computed: {
    isIphone () {
      return window.navigator.userAgent.includes('iPhone')
    },
    selectedFeed () {
      if (!this.filters.feed) {
        return null
      }
      return this.feeds.filter(feed => feed.ID === parseInt(this.filters.feed)).shift()
    }
  },

  methods: {
    onLoad (filters) {
      this.feeds = []
      this.items = []
      this.filters = filters

      this.$http.get(`/feeds`).then(response => {
        this.feeds = response.data
      })

      this.$http.get(`/items`, { params: this.filters }).then(response => {
        this.items = response.data
      })
    },

    onAddFeedClicked () {
      if (this.newFeed === null) {
        this.newFeed = ''
      } else if (this.newFeed !== '') {
        this.$http.post(`/feeds`, { url: this.newFeed }).then(() => {
          this.newFeed = null
          setTimeout(() => window.location.reload(true), 2000)
        })
      } else {
        this.newFeed = null
      }
    },

    onFilterChange () {
      this.changeRouteOnFilterChange(this.filters)
    },

    onReadItLaterClicked (item) {
      this.$http.post(`/items/${item.ID}/readitlater`).then(() => {
        this.items.splice(this.items.indexOf(item), 1)
      })
    },

    onRemoveClicked (item) {
      this.$http.delete(`/items/${item.ID}`).then(() => {
        this.items.splice(this.items.indexOf(item), 1)
      })
    },

    onRefreshFeedClicked (feed) {
      this.$http.post(`/feeds/${feed.ID}/refresh`).then(() => {
        setTimeout(() => window.location.reload(true), 2000)
      })
    },

    onDeleteFeedClicked (feed) {
      this.$http.delete(`/feeds/${feed.ID}`).then(() => {
        this.changeRouteOnFilterChange({})
      })
    }
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
.is-button-pusher {
  width: 100%;
}
</style>
