<template>
  <div>
    <section class="hero is-bold" :class="color">
      <div class="hero-header">
        <div class="container has-text-right logout">
          <p>
            <a @click.prevent="open=true" class="is-size-7">Bookmarklet</a>
            <a @click.prevent="onLogoutClicked" class="is-size-7"><span class="icon"><i class="fa fa-sign-out" aria-hidden="true"></i></span><span>Logout</span></a>
          </p>
        </div>
      </div>
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
              <li><router-link exact active-class="is-active" tag="li" to="/"><a>Bookmarks</a></router-link></li>
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

    <div class="modal" :class="{'is-active':open}">
      <div class="modal-background"></div>
      <div class="modal-content">
        <textarea class="textarea" v-model="snippet"></textarea>
      </div>
      <button class="modal-close is-large" aria-label="close" @click.prevent="open=false"></button>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    snippet () {
      return "javascript:(function(){window.location='"+location.protocol+'//'+location.host+"/bookmarks/save?url='+encodeURIComponent(location.href)+'&title='+encodeURIComponent(document.title);})();"
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
  data: () => ({
    open: false
  }),
  methods: {
    onLogoutClicked (event) {
      this.$http.delete('/token').then(response => {
        this.$router.push({ name: 'login' })
      })
    }
  }
}
</script>

<style scoped>
  .logout {
    margin-top: 1rem;
    margin-right: 2rem;
  }
</style>
