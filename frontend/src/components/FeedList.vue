<template>
  <div>
    <div class="block control is-expanded">
      <div class="select is-fullwidth">
        <select v-model="selected" @change="onChange">
          <option value="">All feeds</option>
          <option v-for="feed in feeds" :key="feed.ID" :value="feed.ID">{{ feed.Title }}</option>
        </select>
      </div>
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
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      feeds: [],
      items: [],
      selected: ''
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
    onChange (event) {
      this.$router.push({query: {feed: this.selected}})
    },
    onReadItLaterClicked (item) {
      axios.post(`/api/items/{item.ID}/readitlater`).then(response => {
        this.items.splice(self.items.indexOf(item), 1)
      })
    },
    onRemoveClicked (item) {
      axios.delete(`/api/items/${item.ID}`).then(response => {
        this.items.splice(self.items.indexOf(item), 1)
      })
    },
    fetchFeeds () {
      var payload = {}
      if (this.selected !== '') {
        payload.feed = this.selected
      }

      axios.get(`/api/items`, {params: payload}).then(response => {
        this.items = response.data
      })
      axios.get('/api/feeds').then(response => {
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
