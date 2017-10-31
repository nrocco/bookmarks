<template>
  <div>
    <section class="hero is-bold" :class="color">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">{{ title }}</h1>
          <h2 class="subtitle">{{ subtitle }}</h2>
        </div>
      </div>
      <div class="hero-foot">
        <nav class="tabs is-boxed is-right">
          <div class="container">
            <ul>
              <li><router-link exact active-class="is-active" tag="li" to="/"><a>Read It Later</a></router-link></li>
              <li><router-link exact active-class="is-active" tag="li" to="/archive"><a>Archive</a></router-link></li>
              <li><router-link exact active-class="is-active" tag="li" to="/feeds"><a>Feeds</a></router-link></li>
            </ul>
          </div>
        </nav>
      </div>
    </section>

    <section class="section">
      <div class="container">
        <router-view></router-view>
      </div>
    </section>

    <footer class="footer">
      <div class="container">
        <div class="content has-text-centered is-size-7">
          <p>Made with love by CasaDiRocco</p>
          <p><a @click.prevent="open=true">bookmarklet</a></p>
          <p><a @click.prevent="logout">logout</a></p>
        </div>
      </div>
    </footer>

    <div class="modal" :class="{'is-active':open}">
      <div class="modal-background"></div>
      <div class="modal-content">
        <textarea class="textarea">javascript:(function(){window.location = '{{ baseurl }}/bookmarks/save?url='+encodeURIComponent(location.href)+'&title='+encodeURIComponent(document.title);})();</textarea>
      </div>
      <button class="modal-close is-large" aria-label="close" @click.prevent="open=false"></button>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      open: false
    }
  },
  computed: {
    baseurl () {
      return location.protocol + '//' + location.host
    },
    title () {
      return this.$route.meta.title
    },
    subtitle () {
      return this.$route.meta.subtitle
    },
    color () {
      return this.$route.meta.color
    }
  },
  methods: {
    logout (event) {
      this.$http.delete('/token').then(response => {
        console.log(response)
      })
    }
  }
}
</script>

<style>
</style>
