<template>
  <div>
    <div class="field has-addons">
      <p class="control">
        <span class="select">
          <select v-model="tag" @change="onSearch">
            <option :value="undefined">all tags</option>
            <option v-for="tag in tags" :value="tag.Name" :key="tag.Name">{{ tag.Name }}</option>
          </select>
        </span>
      </p>
      <p class="control is-expanded">
        <input class="input" type="search" placeholder="Search" v-model="filter" @search="onSearch">
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
      <b-taglist>
        <b-tag closable v-for="tag in bookmark.Tags" :key="tag" @close="onRemoveTagClicked(bookmark, tag)">{{ tag }}</b-tag>
        <b-tag contenteditable @keyup.native.enter="onTagEntered(bookmark, $event)"></b-tag>
      </b-taglist>
      <p class="content">{{ bookmark.Content }}</p>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      filter: null,
      tag: null,
      tags: [],
      bookmarks: []
    }
  },
  methods: {
    onSearch (event) {
      this.$router.push({query: {q: this.filter, tag: this.tag}})
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
    loadTags () {
      this.$http.get(`/tags`).then((response) => {
        this.tags = response.data
      })
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
    this.loadTags()
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
