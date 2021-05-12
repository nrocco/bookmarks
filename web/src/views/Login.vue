<template lang="html">
  <div class="container mt-6">
    <div class="columns is-centered is-mobile">
      <div class="column is-three-quarters-mobile is-half-tablet is-one-third-widescreen">
        <div class="box has-text-centered">
          <figure class="avatar p-5">
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
</template>

<script>
export default {}
</script>

<style lang="css" scoped>
</style>
