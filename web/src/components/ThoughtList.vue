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
        <div class="control" v-if="!filters.q">
          <button class="button is-info" @click="onThoughtNewClicked()">New Thought</button>
        </div>
      </b-field>
    </div>

    <hr/>

    <div class="block thought" v-for="thought in thoughts" :key="thought.ID">
      <p class="is-size-5 has-text-weight-semibold">{{ thought.Tags.join(', ') }}</p>
      <p class="is-size-7 mb-2">
        <time :title="thought.Created">{{ thought.Created|moment("from", "now") }}</time>
        <span> - </span>
        <a @click.prevent="onThoughtModifyClicked(thought)" class="has-text-primary">Modify</a>
        <span> - </span>
        <a @click.prevent="onThoughtRemoveClicked(thought)" class="has-text-danger">Remove</a>
      </p>
      <VueShowdown class="content" :markdown="thought.Content" />
    </div>

    <infinite-loading :identifier="filters" @infinite="infiniteHandler"></infinite-loading>

    <div v-if="thought" class="modal" :class="{'is-active': thought}">
      <div class="modal-background"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <p v-if="thought.ID">Modify thought {{ thought.ID }}</p>
          <p v-else>New Thought</p>
        </header>
        <section class="modal-card-body">
          <b-field :type="{'is-success': thought.Tags.length > 0, 'is-danger': thought.Tags.length === 0}">
            <b-taginput v-model="thought.Tags" allow-new placeholder="Add tags" open-on-focus clear-on-select autocomplete :data="tags" @typing="onTagsTyping"></b-taginput>
          </b-field>
          <b-field :type="{'is-success': thought.Content, 'is-danger': !thought.Content}">
            <b-input v-model="thought.Content" type="textarea" rows="6"></b-input>
          </b-field>
        </section>
        <footer class="modal-card-foot" style="justify-content:space-between;">
          <button class="button" @click="onThoughtModifyClicked(null)">Cancel</button>
          <button class="button is-success" @click="onThoughtSaveClicked" :disabled="thought.Tags.length === 0 || !thought.Content">Save</button>
        </footer>
      </div>
    </div>
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
    limit: 10,
    filters: {},
    thoughts: [],
    thought: null,
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
    }
  },

  methods: {
    onLoad (filters) {
      this.thoughts = []
      this.filters = filters
    },

    infiniteHandler ($state) {
      let payload = Object.assign({ _limit: this.limit, _offset: this.thoughts.length }, this.filters)
      this.$http.get(`/thoughts`, { params: payload }).then(response => {
        this.thoughts.push(...response.data)
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

    onThoughtNewClicked () {
      this.thought = {
        Content: '',
        Tags: []
      }
    },

    onThoughtModifyClicked (thought) {
      this.thought = thought
    },

    onThoughtRemoveClicked (thought) {
      if (!confirm('Are you sure?')) {
        return false
      }
      this.$http.delete(`/thoughts/${thought.ID}`).then(() => {
        this.thoughts.splice(this.thoughts.indexOf(thought), 1)
      })
    },

    onThoughtSaveClicked () {
      this.$http.request({
        method: this.thought.ID ? 'put' : 'post',
        url: this.thought.ID ? `/thoughts/${this.thought.ID}` : '/thoughts',
        data: this.thought.Content,
        headers: {
          'X-Tags': this.thought.Tags.join(',')
        }
      }).then(response => {
        this.thought.Created = response.headers['x-created']
        this.thought.Updated = response.headers['x-updated']
        if (!this.thought.ID) {
          this.thought.ID = response.headers['x-id']
          this.thoughts.unshift(this.thought)
        }
        this.thought = null
      })
    },

    onTagsTyping (value) {
      this.tags = []

      if (!value) {
        return
      }

      // TODO use upstream filtering
      this.$http.get(`/thoughts/_tags`).then(response => {
        this.tags = response.data.filter((tag) => {
          return tag.toString().toLowerCase().indexOf(value.toLowerCase()) >= 0
        });
      })
    }
  }
}
</script>

<style>
.thought {
  background-color: hsl(0, 0%, 98%);
  border-radius: 4px;
  border: 1px solid hsl(0, 0%, 94%);
  padding: 1rem;
}
.thought:hover {
  border: 1px solid hsl(0, 0%, 90%);
}
.thought .content h1,
.thought .content h2,
.thought .content h3,
.thought .content h4 {
  font-size: 1rem;
}
</style>
