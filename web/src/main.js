import Vue from 'vue'
import App from './App'
import Moment from 'vue-moment'

import router from './router'
import store from './store'

import 'bulma/css/bulma.css'

Vue.use(Moment)

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
