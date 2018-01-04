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
      filter: ''
    }
  },
  methods: {
    onFeedClicked (feed) {
      this.$emit('selected', feed)
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
