import Vue from 'vue'
import App from './App'
import Moment from 'vue-moment'

import router from './router'
import store from './store'

Vue.use(Moment)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  store,
  router,
  template: '<App/>',
  components: { App }
})
