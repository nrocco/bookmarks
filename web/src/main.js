import Vue from 'vue'
import Buefy from 'buefy'
import Moment from 'vue-moment'
import VueShowdown from 'vue-showdown'
import App from '@/App.vue'

import client from '@/client.js'
import router from '@/router.js'

import 'buefy/dist/buefy.css'

import { library } from '@fortawesome/fontawesome-svg-core'
import { faAngleLeft, faAngleRight, faBook, faCalendarAlt, faChartLine, faEdit, faEuroSign, faExclamationCircle, faFile, faFilePdf, faFilter, faIndustry, faMinus, faMoneyBill, faPlus, faRulerVertical, faSave, faSearch, faSignOutAlt, faSort, faTimes, faTimesCircle, faTrashAlt, faUpload } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
library.add(faAngleLeft, faAngleRight, faBook, faCalendarAlt, faChartLine, faEdit, faEuroSign, faExclamationCircle, faFile, faFilePdf, faFilter, faIndustry, faMinus, faMoneyBill, faPlus, faRulerVertical, faSave, faSearch, faSignOutAlt, faSort, faTimes, faTimesCircle, faTrashAlt, faUpload)

Vue.config.productionTip = false
Vue.prototype.$http = client
Vue.component('vue-fontawesome', FontAwesomeIcon);

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
