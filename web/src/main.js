import Vue from 'vue'
import Buefy from 'buefy'
import App from './App.vue'
import Moment from 'vue-moment'
import VueShowdown from 'vue-showdown'

import router from './router'
import client from './client'

import 'buefy/dist/buefy.css'

Vue.use(Buefy)
Vue.use(Moment)
Vue.use(VueShowdown, {
  options: {
    simplifiedAutoLink: true
  }
})

Vue.config.productionTip = false
Vue.prototype.$http = client

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
