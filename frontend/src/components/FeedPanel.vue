<template>
  <nav class="panel">
    <p class="panel-heading">Feeds</p>
    <div class="panel-block">
      <div class="field has-addons">
        <div class="control has-icons-left is-expanded">
          <input class="input is-small" type="text" placeholder="search" v-model="filter">
          <span class="icon is-small is-left">
            <i class="fa fa-search"></i>
          </span>
        </div>
        <div class="control">
          <div class="select is-small">
            <select v-model="selectedCategory">
              <option :value="null">All</option>
              <option :value="category" v-for="category in categories">{{ category }}</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <a class="panel-block" v-for="feed in filteredFeeds" @click.prevent="onFeedClicked(feed)"><span style="width: 100%;">{{ feed.Title }}</span> <span class="tag is-rounded">{{ feed.Items }}</span></a>

    <div class="panel-block">
      <a v-if="!newFeed" class="button is-primary is-outlined is-rounded is-small is-fullwidth" @click.prevent="onAddNewFeedClicked"><i class="fa fa-plus" aria-hidden="true"></i>&nbsp;Add Feed</a>
      <div v-else class="field has-addons" style="width:100%;">
        <div class="control is-expanded">
          <input class="input is-small is-fullwidth" type="text" v-model="newFeed.URL" placeholder="Feed URL" @keyup.esc="newFeed=null" @keyup.enter="onSaveNewFeedClicked">
        </div>
        <div class="control">
          <a class="button is-small is-primary" @click.prevent="onSaveNewFeedClicked"><i class="fa fa-plus" aria-hidden="true"></i></a>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  data () {
    return {
      selectedCategory: null,
      newFeed: null,
      filter: ''
    }
  },
  methods: {
    ...mapActions({
      onFeedClicked: 'selectFeed'
    }),
    onSaveNewFeedClicked (event) {
      this.$store.dispatch('addFeed', this.newFeed)
    },
    onAddNewFeedClicked (event) {
      this.newFeed = {
        Category: 'default'
      }
    }
  },
  computed: {
    filteredFeeds () {
      return this.$store.getters.feeds.filter(feed => {
        return (!this.selectedCategory || feed.Category === this.selectedCategory) && feed.Title.toLowerCase().indexOf(this.filter.toLowerCase()) > -1
      })
    },
    categories () {
      return this.$store.getters.categories
    }
  }
}
</script>

<style scoped>
</style>
