import Vue from 'vue'
import App from './App'
import Moment from 'vue-moment'

import router from './router'
import client from './client'

import 'bulma/css/bulma.css'

Vue.use(Moment)

Vue.config.productionTip = false
Vue.prototype.$http = client

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
