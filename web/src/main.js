import Vue from 'vue'
import Buefy from 'buefy'
import App from './App'
import Moment from 'vue-moment'

import router from './router'
import client from './client'

import 'buefy/dist/buefy.css'

Vue.use(Moment)
Vue.use(Buefy)

Vue.config.productionTip = false
Vue.prototype.$http = client

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
