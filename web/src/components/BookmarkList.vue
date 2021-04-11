<template>
  <div>
    <div class="field has-addons">
      <p class="control is-expanded">
        <input class="input" type="search" placeholder="Search" v-model="filters.q" @search="onFilterChange">
      </p>
    </div>

    <hr/>

    <div class="block bookmark" v-for="bookmark in bookmarks" :key="bookmark.ID">
      <p class="is-size-5 has-text-weight-semibold">{{ bookmark.Title }}</p>
      <p class="is-size-7 mb-2">
        <time :title="bookmark.Created">{{ bookmark.Created|moment("from", "now") }}</time>
        <span> - </span>
        <a class="url" :href="bookmark.URL" :target="isIphone ? '_blank' : ''">{{ bookmark.URL }}</a>
        <span> - </span>
        <a @click.prevent="onRemoveClicked(bookmark)" class="has-text-danger">Remove</a>
      </p>
      <p>{{ bookmark.Excerpt }}&#8230;</p>
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
    filters: {},
    bookmarks: []
  }),

  computed: {
    isIphone () {
      return window.navigator.userAgent.includes('iPhone')
    }
  },

  methods: {
    onLoad (filters) {
      this.bookmarks = []
      this.filters = filters

      let params = {}
      if (this.filters.q) {
        params.q = this.filters.q
      }

      this.$http.get(`/bookmarks`, { params: params }).then(response => {
        this.bookmarks = response.data
      })
    },

    onFilterChange () {
      this.changeRouteOnFilterChange(this.filters)
    },

    onRemoveClicked (bookmark) {
      this.$http.delete(`/bookmarks/${bookmark.ID}`).then(() => {
        this.bookmarks.splice(this.bookmarks.indexOf(bookmark), 1)
      })
    }
  }
}
</script>

<style>
.bookmark {
  padding: 1rem;
  background-color: hsl(0, 0%, 100%);
  border: 1px solid hsl(0, 0%, 97%);
  border-radius: 4px;
}
.bookmark:hover {
  background-color: hsl(0, 0%, 98%);
}
.bookmark .url {
  word-break: break-all;
}
</style>
