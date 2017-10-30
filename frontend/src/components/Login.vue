<template>
  <section class="hero is-light is-fullheight">
    <div class="hero-body">
      <div class="container has-text-centered">
        <div class="column is-4 is-offset-4">
          <div class="box">
            <figure class="avatar">
              <img src="../assets/logo.png">
            </figure>
            <form method="post" action="/api/token">
              <input type="hidden" name="next" value="/" />
              <div class="field">
                <div class="control">
                  <input class="input" name="username" type="text" placeholder="Username" autofocus>
                </div>
              </div>
              <div class="field">
                <div class="control">
                  <input class="input" name="password" type="password" placeholder="Password">
                </div>
              </div>
              <input type="submit" class="button is-fullwidth is-dark" value="Login" />
            </form>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {}
</script>

<style scoped>
.hero .nav, .hero.is-success .nav {
  -webkit-box-shadow: none;
  box-shadow: none;
}
.avatar {
  margin-top: -40px;
  padding-bottom: 10px;
}
.avatar img {
  padding: 5px;
}
</style>
