<template>
  <div>
    <div class="block">
      <div class="field has-addons">
        <div class="control is-expanded">
          <input class="input" type="search" placeholder="Search" v-model="filters.q" @search="onFilterChange">
        </div>
        <div class="control" v-if="!filters.q">
          <button class="button is-info" @click="onNewClicked()">New Thought</button>
        </div>
      </div>
    </div>

    <hr/>

    <div class="block thought" v-for="thought in thoughts" :key="thought.ID">
      <p class="is-size-5 has-text-weight-semibold">{{ thought.Tags.join(', ') }}</p>
      <p class="is-size-7 mb-2">
        <time :title="thought.Created">{{ thought.Created|moment("from", "now") }}</time>
        <span> - </span>
        <a @click.prevent="onModifyClicked(thought)" class="has-text-primary">Modify</a>
        <span> - </span>
        <a @click.prevent="onRemoveClicked(thought)" class="has-text-danger">Remove</a>
      </p>
      <VueShowdown class="content" :markdown="thought.Content" />
    </div>

    <div v-if="thought" class="modal" :class="{'is-active': thought}">
      <div class="modal-background"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <p v-if="thought.ID">Modify thought {{ thought.ID }}</p>
          <p v-else>New Thought</p>
        </header>
        <section class="modal-card-body">
          <div class="field is-grouped is-grouped-multiline">
            <div class="control" v-for="tag in thought.Tags" :key="tag">
              <div class="tags has-addons">
                <span class="tag is-light">{{ tag }}</span>
                <a @click="onRemoveTag(tag)" class="tag is-delete"></a>
              </div>
            </div>
            <div class="control">
              <input @keyup.enter="onAddTag" @keydown.tab="onAddTag" @blur="onAddTag" class="input is-small" type="text" placeholder="Tag" autofocus>
            </div>
          </div>
          <textarea v-model="thought.Content" class="textarea" rows="6"></textarea>
        </section>
        <footer class="modal-card-foot" style="justify-content:space-between;">
          <button class="button" @click="onModifyClicked(null)">Cancel</button>
          <button class="button is-success" @click="onSaveClicked" :disabled="thought.Tags.length === 0 || !thought.Content">Save</button>
        </footer>
      </div>
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
    thoughts: [],
    thought: null
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
        }).catch(() => {
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

    onFilterChange () {
      this.changeRouteOnFilterChange(this.filters, '/thoughts')
    },

    onNewClicked () {
      this.thought = {
        Content: '',
        Tags: []
      }
    },

    onModifyClicked (thought) {
      this.thought = thought
    },

    onRemoveClicked (thought) {
      if (!confirm('Are you sure?')) {
        return false
      }
      this.$http.delete(`/thoughts/${thought.ID}`).then(() => {
        this.thoughts.splice(this.thoughts.indexOf(thought), 1)
      })
    },

    onSaveClicked () {
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

    onRemoveTag (tag) {
      this.thought.Tags.splice(this.thought.Tags.indexOf(tag), 1)
    },

    onAddTag (event) {
      if (!event.target.value) {
        return
      }
      this.thought.Tags.push(event.target.value)
      event.target.value = ''
    }
  }
}
</script>

<style>
.thought {
  padding: 1rem;
  background-color: hsl(0, 0%, 99%);
  border: 1px solid hsl(0, 0%, 97%);
  border-radius: 4px;
}
.thought:hover {
  background-color: hsl(0, 0%, 98%);
}
.thought .content h1,
.thought .content h2,
.thought .content h3,
.thought .content h4 {
  font-size: 1rem;
}
</style>
