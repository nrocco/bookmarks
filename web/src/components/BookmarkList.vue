<template>
  <div>
    <div class="field has-addons">
      <p class="control">
        <span class="select">
          <select v-model="filters.tag" @change="onFilterChange">
            <option :value="undefined">all tags</option>
            <option v-for="tag in tags" :value="tag.Name" :key="tag.Name">{{ tag.Name }}</option>
          </select>
        </span>
      </p>
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
        <a @click.prevent="onRemoveClicked(bookmark)" class="has-text-danger">Remove</a>
      </p>
      <div class="tags">
        <span class="tag" v-for="tag in bookmark.Tags" :key="bookmark.ID+tag">{{ tag }} <button class="delete is-small" @click.prevent="onRemoveTagClicked(bookmark, tag)"></button></span>
        <span class="tag" contenteditable @keyup.enter="onTagEntered(bookmark, $event)"></span>
      </div>
      <p class="content">{{ bookmark.Content }}</p>
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
    tags: [],
    bookmarks: []
  }),

  methods: {
    onLoad (filters) {
      this.bookmarks = []
      this.tags = []
      this.filters = filters

      this.$http.get(`/bookmarks`, { params: this.filters }).then(response => {
        this.bookmarks = response.data
      })

      this.$http.get(`/tags`).then((response) => {
        this.tags = response.data
      })
    },

    onFilterChange (event) {
      this.changeRouteOnFilterChange(this.filters)
    },

    onRemoveTagClicked (bookmark, tag) {
      bookmark.Tags.splice(bookmark.Tags.indexOf(tag), 1)
      this.$http.patch(`/bookmarks/${bookmark.ID}`, { Tags: bookmark.Tags }).then(response => {
        bookmark = response.data
      })
    },

    onTagEntered (bookmark, event) {
      let tag = event.target.innerText.toString().trim()
      if (!tag) {
        return
      }
      bookmark.Tags.push(tag)
      this.$http.patch(`/bookmarks/${bookmark.ID}`, { Tags: bookmark.Tags }).then(response => {
        bookmark = response.data
      })
      event.target.innerText = ''
      return false
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
.bookmark-tags {
  margin-top: 0.75rem;
}
.bookmark .url {
  word-break: break-all;
}
.content:not(:last-child) {
  margin-bottom: 0.5rem;
}
.bookmark .tags:not(:last-child) {
  margin-top: 0.5rem;
  margin-bottom: 0;
}
</style>
