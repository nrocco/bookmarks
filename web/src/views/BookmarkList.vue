<template>
  <div>
    <div class="block">
      <b-field grouped>
        <div class="control">
          <b-taginput placeholder="Filter tags" v-model="filterTags" autocomplete :data="tags" @typing="onTagsTyping" @input="onFilterChange"></b-taginput>
        </div>
        <div class="control is-expanded">
          <input class="input" type="search" placeholder="Search" v-model="filters.q" @search="onFilterChange">
        </div>
      </b-field>
    </div>

    <hr/>

    <div class="block bookmark" v-for="bookmark in bookmarks" :key="bookmark.ID">
      <p class="is-size-5 has-text-weight-semibold">{{ bookmark.Title }}</p>
      <p class="is-size-7 mb-2">
        <time :title="bookmark.Created">{{ bookmark.Created|moment("from", "now") }}</time>
        <span> - </span>
        <a class="url" :href="bookmark.URL" :target="isIphone ? '_blank' : ''">{{ bookmark.URL }}</a>
        <span> - </span>
        <a @click.prevent="onBookmarkRemoveClicked(bookmark)" class="has-text-danger">Remove</a>
      </p>
      <p>
        <b-tag v-for="tag in bookmark.Tags" :key="tag" type="is-link" closable @close="onBookmarkTagRemoved(bookmark, tag)">{{ tag }}</b-tag>
        <span v-if="bookmark.Tags.length > 0"> - </span>
        <span>{{ bookmark.Excerpt }}&#8230;</span>
      </p>
    </div>

    <infinite-loading :identifier="filters" @infinite="onInfiniteScroll">
      <span slot="no-more"></span>
      <div slot="no-results">
        <i>No bookmarks found!</i>
      </div>
    </infinite-loading>
  </div>
</template>

<script>
import LoaderMixin from '@/helpers.js'
import InfiniteLoading from 'vue-infinite-loading';

export default {
  mixins: [
    LoaderMixin
  ],

  components: {
    InfiniteLoading,
  },

  data: () => ({
    filters: {},
    bookmarks: [],
    tags: []
  }),

  computed: {
    filterTags: {
      get () {
        if (!this.filters.tags) {
          return []
        }
        return this.filters.tags.split(',')
      },
      set (value) {
        this.filters.tags = value.join(',')
      }
    },
    isIphone () {
      return window.navigator.userAgent.includes('iPhone')
    }
  },

  methods: {
    onLoad (filters) {
      this.bookmarks = []
      this.filters = filters
    },

    onInfiniteScroll ($state) {
      let payload = Object.assign({ _limit: 20, _offset: this.bookmarks.length }, this.filters)
      this.$http.get(`/bookmarks`, { params: payload }).then(response => {
        this.bookmarks.push(...response.data)
        if (response.data.length > 0) {
          $state.loaded()
        }
        if (response.data.length < 20) {
          $state.complete()
        }
      })
    },

    onFilterChange () {
      this.changeRouteOnFilterChange(this.filters)
    },

    onBookmarkRemoveClicked (bookmark) {
      this.$http.delete(`/bookmarks/${bookmark.ID}`).then(() => {
        this.bookmarks.splice(this.bookmarks.indexOf(bookmark), 1)
      })
    },

    onBookmarkTagRemoved (bookmark, tag) {
      bookmark.Tags.splice(bookmark.Tags.indexOf(tag), 1)
      this.$http.patch(`/bookmarks/${bookmark.ID}`, {Tags: bookmark.Tags})
    },

    onTagsTyping (value) {
      this.tags = []

      if (!value) {
        return
      }

      // TODO add rest api to get bookmark tags
      this.tags = ['read-it-later'].filter((tag) => {
        return tag.toString().toLowerCase().indexOf(value.toLowerCase()) >= 0
      })
    }
  }
}
</script>

<style scoped>
.bookmark {
  background-color: hsl(0, 0%, 98%);
  border-radius: 4px;
  border: 1px solid hsl(0, 0%, 94%);
  padding: 1rem;
}
.bookmark:hover {
  border: 1px solid hsl(0, 0%, 90%);
}
.bookmark .url {
  word-break: break-all;
}
</style>
