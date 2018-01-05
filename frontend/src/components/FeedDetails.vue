<template>
    <div class="block">
        <h1 class="title">{{ feed.Title }}</h1>
        <h2 class="subtitle">{{ feed.URL }}</h2>
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
export default {
  props: {
    feed: Object
  },
  data () {
    return {
      items: []
    }
  },
  methods: {
    onRefreshFeedClicked (feed) {
      this.$http.post(`/feeds/${feed.ID}/refresh`).then(response => {
        this.$emit('refreshed', feed)
        // TODO sleep a second
        // TODO load items again
      })
    },
    onDeleteFeedClicked (feed) {
      this.$http.delete(`/feeds/${feed.ID}`).then(response => {
        this.$emit('deleted', feed)
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
    fetchItems () {
      this.$http.get(`/items`, {params: {feed: this.feed.ID}}).then(response => {
        this.items = response.data
      })
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
  watch: {
    feed: {
      immediate: true,
      handler (feed) {
        if (feed) {
          this.fetchItems()
        }
      }
    }
  }
}
</script>

<style>
</style>
