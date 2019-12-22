<template>
  <div>
    <section class="hero is-bold" :class="color">
      <div class="hero-header">
        <div class="container has-text-right logout">
          <p>
            <a @click="isBookmarkletModalActive = true" class="is-size-7">Bookmarklet</a>
            <a @click.prevent="onLogoutClicked" class="is-size-7"><span class="icon"><i class="fa fa-sign-out" aria-hidden="true"></i></span><span>Logout</span></a>
          </p>
        </div>
      </div>
      <div class="hero-foot">
        <nav class="tabs is-centered">
          <ul>
            <li><router-link exact active-class="is-active" tag="li" to="/"><a>Bookmarks</a></router-link></li>
            <li><router-link exact active-class="is-active" tag="li" to="/feeds"><a>Feeds</a></router-link></li>
            <li><router-link exact active-class="is-active" tag="li" to="/thoughts"><a>Thoughts</a></router-link></li>
          </ul>
        </nav>
      </div>
    </section>

    <section class="section">
      <div class="container">
        <router-view></router-view>
      </div>
    </section>

    <div class="modal" scroll="keep" :class="{'is-active': isBookmarkletModalActive}">
      <div class="modal-background" @click="isBookmarkletModalActive = false"></div>
      <div class="modal-card">
        <section class="modal-card-body">
          <pre class="bookmarklet">javascript:(function(){window.location='{{ baseurl }}/bookmarks/save?url='+encodeURIComponent(location.href)+'&title='+encodeURIComponent(document.title);})();</pre>
          <p>Bookmark this page, then replace the url of the bookmark you just created with the above javascript snippet.</p>
        </section>
      </div>
    </div>
  </div>
</template>

<script>
export default {
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
  data: () => ({
    isBookmarkletModalActive: false
  }),
  methods: {
    onLogoutClicked () {
      this.$http.delete('/token').then(() => {
        this.$router.push({ name: 'login' })
      })
    }
  }
}
</script>

<style scoped>
  .bookmarklet {
    word-wrap: break-word;
    word-break: break-all;
  }
  .logout {
    margin-top: 1rem;
    margin-right: 2rem;
  }
</style>
