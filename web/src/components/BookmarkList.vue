<template>
  <div>
    <div class="field has-addons">
      <p class="control is-expanded">
        <input class="input" type="search" placeholder="Search" v-model="filters.q" @search="onFilterChange">
      </p>
    </div>
    <hr/>
    <div class="block bookmark" v-for="bookmark in bookmarks" :key="bookmark.ID">
      <p class="has-text-weight-bold">{{ bookmark.Title }}</p>
      <p class="is-size-7">
        <a class="url" :href="bookmark.URL">{{ bookmark.URL }}</a>
        <span> - </span>
        <a @click.prevent="onToggleArchivedClicked(bookmark)" :class="{'has-text-primary': !bookmark.Archived, 'has-text-info': bookmark.Archived}">{{ bookmark.Archived ? 'Read it Later' : 'Archive' }}</a>
        <span> - </span>
        <a @click.prevent="onRemoveClicked(bookmark)" class="has-text-danger">Remove</a>
      </p>
      <p class="content">{{ bookmark.Content }}&#8230;</p>
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

  methods: {
    onLoad (filters) {
      this.bookmarks = []
      this.filters = filters

      let params = {}
      if (!this.filters.q) {
        params.readitlater = 'true'
      } else {
        params.q = this.filters.q
      }

      this.$http.get(`/bookmarks`, { params: params }).then(response => {
        this.bookmarks = response.data
      })
    },

    onFilterChange (event) {
      this.changeRouteOnFilterChange(this.filters)
    },

    onToggleArchivedClicked (bookmark) {
      this.$http.patch(`/bookmarks/${bookmark.ID}`, { Archived: !bookmark.Archived }).then(response => {
        this.bookmarks.splice(this.bookmarks.indexOf(bookmark), 1)
      })
    },

    onRemoveClicked (bookmark) {
      this.$http.delete(`/bookmarks/${bookmark.ID}`).then(response => {
        this.bookmarks.splice(this.bookmarks.indexOf(bookmark), 1)
      })
    }
  }
}
</script>

<style>
.bookmark {
  border-bottom: 1px solid hsl(0, 0%, 96%);
  padding-bottom: 1.5rem;
}
.bookmark .url {
  word-break: break-all;
}
.content:not(:last-child) {
  margin-bottom: 0.5rem;
}
</style>
