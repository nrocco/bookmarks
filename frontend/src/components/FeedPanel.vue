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
                        <select v-model="selected_category">
                            <option :value="null">All</option>
                            <option :value="category" v-for="category in categories">{{ category }}</option>
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <a class="panel-block" v-for="feed in feeds_filtered" @click.prevent="onFeedClicked(feed)"><span style="width: 100%;">{{ feed.Title }}</span> <span class="tag is-rounded">{{ feed.Items }}</span></a>

        <div class="panel-block">
            <a v-if="!new_feed" class="button is-primary is-outlined is-rounded is-small is-fullwidth" @click.prevent="onAddNewFeedClicked"><i class="fa fa-plus" aria-hidden="true"></i>&nbsp;Add Feed</a>
            <div v-else class="field has-addons" style="width:100%;">
                <div class="control is-expanded">
                    <input class="input is-small is-fullwidth" type="text" v-model="new_feed.URL" placeholder="Feed URL" @keyup.esc="new_feed=null" @keyup.enter="onSaveNewFeedClicked">
                </div>
                <div class="control">
                    <a class="button is-small is-primary" @click.prevent="onSaveNewFeedClicked"><i class="fa fa-plus" aria-hidden="true"></i></a>
                </div>
            </div>
        </div>
    </nav>
</template>

<script>
export default {
  props: {
    feeds: Array
  },
  data () {
    return {
      selected_category: null,
      new_feed: null,
      filter: ''
    }
  },
  methods: {
    onFeedClicked (feed) {
      this.$emit('selected', feed)
    },
    onAddNewFeedClicked (event) {
      this.new_feed = {
        Category: 'default'
      }
    },
    onSaveNewFeedClicked (event) {
      this.$http.post(`/feeds`, this.new_feed).then(response => {
        this.$emit('added', response.data)
        this.$emit('selected', response.data)
        this.new_feed = null
      })
    }
  },
  computed: {
    feeds_filtered () {
      return this.feeds.filter(feed => {
        return (!this.selected_category || feed.Category === this.selected_category) && feed.Title.toLowerCase().indexOf(this.filter.toLowerCase()) > -1
      })
    },
    categories () {
      let categories = []
      for (let feed of this.feeds) {
        categories.push(feed.Category)
      }
      return [...new Set(categories)]
    }
  }
}
</script>

<style scoped>
</style>
