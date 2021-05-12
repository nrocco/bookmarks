import Vue from 'vue'
import Buefy from 'buefy'
import Moment from 'vue-moment'
import VueShowdown from 'vue-showdown'
import App from '@/App.vue'

import client from '@/client.js'
import router from '@/router.js'

import 'buefy/dist/buefy.css'

Vue.config.productionTip = false
Vue.prototype.$http = client

Vue.use(Buefy)
Vue.use(VueShowdown, {
  options: {
    simplifiedAutoLink: true
  }
})
Vue.use(Moment)

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
