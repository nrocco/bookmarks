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
      <p>
        <span v-for="tag in bookmark.Tags" :key="tag" class="tag is-link">{{ tag }}</span>
        <span v-if="bookmark.Tags.length > 0"> - </span>
        <span>{{ bookmark.Excerpt }}&#8230;</span>
      </p>
    </div>

    <infinite-loading :identifier="filters" @infinite="infiniteHandler"></infinite-loading>
  </div>
</template>

<script>
import LoaderMixin from '../mixins/loader.js'
import InfiniteLoading from 'vue-infinite-loading';

export default {
  mixins: [
    LoaderMixin
  ],

  components: {
    InfiniteLoading,
  },

  data: () => ({
    limit: 20,
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
    },

    infiniteHandler ($state) {
      let payload = Object.assign({ _limit: this.limit, _offset: this.bookmarks.length }, this.filters)
      this.$http.get(`/bookmarks`, { params: payload }).then(response => {
        this.bookmarks.push(...response.data)
        if (response.data.length === 0 || this.limit > response.data.length) {
          $state.complete()
        } else {
          $state.loaded()
        }
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
