<template>
  <div>
    <div class="field has-addons">
      <p class="control">
        <span class="select">
          <select v-model="tag" @change="onSearch">
            <option :value="undefined">all tags</option>
            <option>read-it-later</option>
          </select>
        </span>
      </p>
      <p class="control is-expanded">
        <input class="input" type="search" placeholder="Search" v-model="filter" @search="onSearch">
      </p>
    </div>

    <hr/>

    <div class="block bookmark" v-for="bookmark in filteredBookmarks" :key="bookmark.ID">
      <p class="has-text-weight-bold">{{ bookmark.Title }}</p>
      <p class="is-size-7"><a class="url" :href="bookmark.URL">{{ bookmark.URL }}</a></p>
      <div class="tags has-addons" v-if="bookmark.Tags" v-for="tag in bookmark.Tags" :key="tag">
        <span class="tag">{{ tag }}</span>
        <a class="tag is-delete"></a>
      </div>
      <p class="content">{{ bookmark.Content }}</p>
      <p class="buttons is-right">
        <a @click.prevent="onRemoveClicked(bookmark)" class="button is-small is-danger is-outlined">Remove</a>
        <a @click.prevent="onReadItLaterClicked(bookmark)" class="button is-small is-primary" v-if="bookmark.Tags.indexOf('read-it-later') == -1">Read it later</a>
        <a @click.prevent="onArchiveClicked(bookmark)" class="button is-small is-dark" v-else>Archive</a>
      </p>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      filter: null,
      tag: 'read-it-later',
      bookmarks: []
    }
  },
  computed: {
    filteredBookmarks () {
      return this.bookmarks.filter(el => !this.tag || el.Tags.indexOf(this.tag) !== -1)
    }
  },
  methods: {
    onSearch (event) {
      this.$router.push({ query: { q: this.filter, tag: this.tag } })
    },
    loadBookmarks () {
      let payload = {}
      if (this.filter !== '') {
        payload.q = this.filter
      }
      if (this.tag) {
        payload.tag = this.tag
      }
      this.$http.get(`/bookmarks`, { params: payload }).then(response => {
        this.bookmarks = response.data
      })
    },
    onReadItLaterClicked (bookmark) {
      bookmark.Tags.push('read-it-later')
      this.$http.patch(`/bookmarks/${bookmark.ID}`, { Tags: bookmark.Tags }).then(response => {
        bookmark = response.data
      })
    },
    onArchiveClicked (bookmark) {
      bookmark.Tags.splice(bookmark.Tags.indexOf('read-it-later'), 1)
      this.$http.patch(`/bookmarks/${bookmark.ID}`, { Tags: bookmark.Tags }).then(response => {
        bookmark = response.data
      })
    },
    onRemoveClicked (bookmark) {
      this.$http.delete(`/bookmarks/${bookmark.ID}`).then(response => {
        this.bookmarks.splice(this.bookmarks.indexOf(bookmark), 1)
      })
    }
  },
  watch: {
    '$route' (to, from) {
      this.filter = to.query.q
      this.tag = to.query.tag
      this.loadBookmarks()
    }
  },
  mounted () {
    this.filter = this.$route.query.q
    this.tag = this.$route.query.tag
    this.loadBookmarks()
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
.bookmark .tags:not(:last-child) {
  margin-top: 0.5rem;
  margin-bottom: 0;
}
</style>
