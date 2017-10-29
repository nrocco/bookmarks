<template>
  <div>
    <div class="block control">
      <input class="input" type="search" placeholder="Search" v-model="q" autofocus @search="onSearch">
    </div>
    <div class="block" v-for="bookmark in bookmarks" :key="bookmark.ID">
      <p class="has-text-weight-bold">{{ bookmark.Title }}</p>
      <p class="is-size-7"><a :href="bookmark.URL">{{ bookmark.URL }}</a></p>
      <p class="content">{{ bookmark.Content|excerpt }}</p>
      <p class="block has-text-right">
        <a @click.prevent="onRemoveClicked(bookmark)" class="button is-small is-danger is-outlined">Remove</a>
        <a @click.prevent="onReadItLaterClicked(bookmark)" class="button is-small is-primary" v-if="bookmark.Archived">Read it later</a>
        <a @click.prevent="onArchiveClicked(bookmark)" class="button is-small is-dark" v-else="bookmark.Archived">Archive</a>
      </p>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      bookmarks: [],
      q: ''
    }
  },
  filters: {
    excerpt: function (value) {
      if (!value) {
        return 'No content'
      }
      return value.toString().substring(0, 500) + '...'
    }
  },
  methods: {
    onSearch (event) {
      this.$router.push({query: {q: this.q}})
    },
    onReadItLaterClicked (bookmark) {
      axios.post(`/api/bookmarks/${bookmark.ID}/readitlater`).then(response => {
        this.bookmarks.splice(self.bookmarks.indexOf(bookmark), 1)
      })
    },
    onArchiveClicked (bookmark) {
      axios.post(`/api/bookmarks/${bookmark.ID}/archive`).then(response => {
        this.bookmarks.splice(self.bookmarks.indexOf(bookmark), 1)
      })
    },
    onRemoveClicked (bookmark) {
      axios.delete(`/api/bookmarks/${bookmark.ID}`).then(response => {
        this.bookmarks.splice(self.bookmarks.indexOf(bookmark), 1)
      })
    },
    fetchBookmarks () {
      var payload = {}

      if (this.$route.meta.archived) {
        payload.archived = 'true'
      }
      if (this.q !== '') {
        payload.q = this.q
      }

      axios.get(`/api/bookmarks`, {params: payload}).then(response => {
        this.bookmarks = response.data
      })
    }
  },
  beforeRouteUpdate (to, from, next) {
    next()
    this.fetchBookmarks()
  },
  beforeRouteEnter (to, from, next) {
    next(vm => {
      vm.fetchBookmarks()
    })
  }
}
</script>

<style>
</style>
