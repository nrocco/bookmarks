<template>
  <div>
    <div class="field has-addons">
      <p class="control is-expanded">
        <input class="input" type="search" placeholder="Search" v-model="filters.q" @search="onFilterChange">
      </p>
    </div>
    <hr/>
    <thought v-for="thought in thoughts" :thought="thought" :key="thought.Title" @saved="onThoughtSaved(thought, $event)" @removed="onThoughtRemoved" />
  </div>
</template>

<script>
import LoaderMixin from '../mixins/loader.js'
import Thought from './Thought.vue'

export default {
  mixins: [
    LoaderMixin
  ],

  components: {
    Thought
  },

  data: () => ({
    filters: {},
    thoughts: []
  }),

  computed: {
  },

  methods: {
    onLoad (filters) {
      this.thoughts = []
      this.filters = filters

      if (this.$route.params.title) {
        this.$http.get(`/thoughts/${this.$route.params.title}`).then(response => {
          this.thoughts.push({
            Title: this.$route.params.title,
            Created: response.headers['x-created'],
            Updated: response.headers['x-updated'],
            Tags: response.headers['x-tags'].split(','),
            Content: response.data
          })
        }).catch(error => {
          console.log(error)
          this.thoughts.push({
            Title: this.$route.params.title,
            Created: null,
            Updated: null,
            Tags: [],
            Content: ''
          })
        })
      } else {
        let params = {}
        if (this.filters.q) {
          params.q = this.filters.q
        }

        this.$http.get(`/thoughts`, { params: params }).then(response => {
          this.thoughts = response.data
        })
      }
    },

    onFilterChange (event) {
      this.changeRouteOnFilterChange(this.filters, '/thoughts')
    },

    onThoughtSaved (oldThought, newThought) {
      this.thoughts.splice(this.thoughts.indexOf(oldThought), 1, newThought)
    },

    onThoughtRemoved (thought) {
      this.thoughts.splice(this.thoughts.indexOf(thought), 1)
    }
  }
}
</script>

<style>
</style>
