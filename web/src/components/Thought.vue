<template>
  <div class="box thought">
    <div class="media-content">
      <div class="content">
        <div class="columns">
          <div class="column is-four-fifths">
            <h4 class="title is-4"><router-link :to="{name: 'thoughts',params:{title:thought.Title}}">{{ thought.Title }}</router-link></h4>
            <p class="subtitle is-6">
              <i v-if="thought.Updated">Modified {{ thought.Updated|moment("from", "now") }}</i>
              <i v-else>New thought</i>
            </p>
          </div>

          <div v-if="!isEditing" class="column has-text-right">
            <a :href="'/api/thoughts/'+thought.Title" class="button is-white">
              <i class="fas fa-download" aria-hidden="true"></i>
            </a>
            <a @click="onEditClicked" class="button is-white">
              <i class="fas fa-edit" aria-hidden="true"></i>
            </a>
          </div>
          <div v-else class="column has-text-right">
            <a @click="onRemoveClicked(thought)" class="button is-white">
              <i class="fas fa-trash" aria-hidden="true"></i>
            </a>
            <a @click="onCancelClicked" class="button is-white">
              <i class="fas fa-times" aria-hidden="true"></i>
            </a>
          </div>
        </div>

        <div v-if="!isEditing">
          <div class="tags">
            <span v-for="tag in thought.Tags" :key="tag" class="tag is-light">{{ tag }}</span>
          </div>
          <p class="content"><VueShowdown :markdown="thought.Content" /></p>
        </div>
        <div v-else>
          <div class="field is-grouped is-grouped-multiline">
            <div class="control" v-for="tag in modifiedThought.Tags" :key="tag">
              <div class="tags has-addons">
                <span class="tag is-light">{{ tag }}</span>
                <a @click="onRemoveTag(tag)" class="tag is-delete"></a>
              </div>
            </div>
            <div class="control">
              <input @keyup.enter="onTagAdded" class="input is-small" type="text" placeholder="Text input">
            </div>
          </div>
          <textarea v-model="modifiedThought.Content" @keyup="onContentChanged" class="textarea" placeholder="10 lines of textarea" rows="10"></textarea>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import _ from 'lodash'

export default {
  props: {
    thought: {
      type: [Object],
      default: () => ({})
    }
  },

  data: () => ({
    modifiedThought: null
  }),

  computed: {
    isEditing () {
      return this.modifiedThought !== null
    }
  },

  methods: {
    onEditClicked (event) {
      this.modifiedThought = JSON.parse(JSON.stringify(this.thought))
    },

    onCancelClicked (event) {
      this.modifiedThought = null
    },

    onContentChanged: _.debounce(function () {
      this.$http.put(`/thoughts/${this.thought.Title}`, this.modifiedThought.Content).then(response => {
        this.$emit('saved', this.modifiedThought)
      })
    }, 1000),

    onSaveClicked (event) {
      this.modifiedThought = null
    },

    onTagAdded (event) {
      this.modifiedThought.Tags.push(event.target.value)
      event.target.value = ''

      this.$http.put(`/thoughts/${this.thought.Title}`, null, {
        headers: {
          'X-Tags': this.modifiedThought.Tags.join(',')
        }
      }).then(response => {
        this.modifiedThought.Created = response.headers['x-created']
        this.modifiedThought.Updated = response.headers['x-updated']
        this.$emit('saved', this.modifiedThought)
      })
    },

    onRemoveTag (tag) {
      this.modifiedThought.Tags.splice(this.modifiedThought.Tags.indexOf(tag), 1)
      this.$http.put(`/thoughts/${this.thought.Title}`, null, {
        headers: {
          'X-Tags': this.modifiedThought.Tags.join(',')
        }
      }).then(response => {
        this.modifiedThought.Created = response.headers['x-created']
        this.modifiedThought.Updated = response.headers['x-updated']
        this.$emit('saved', this.modifiedThought)
      })
    },

    onRemoveClicked (thought) {
      if (!confirm('Are you sure?')) {
        return false
      }
      this.$http.delete(`/thoughts/${this.thought.Title}`).then(response => {
        this.$emit('removed', this.thought)
      })
    }
  }
}
</script>

<style>
</style>
