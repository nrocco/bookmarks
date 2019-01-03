'use strict'

import Vue from 'vue'
import Vuex from 'vuex'

import bookmarks from './modules/bookmarks'
import feeds from './modules/feeds'

Vue.use(Vuex)

export default new Vuex.Store({
  strict: process.env.NODE_ENV !== 'production',
  modules: {
    bookmarks,
    feeds
  },
  actions: {
    readitlater (context, route) {
      context.commit('filter', route.query.q)
      context.commit('archived', false)
      context.dispatch('getBookmarks')
    },
    archive (context, route) {
      context.commit('filter', route.query.q)
      context.commit('archived', true)
      context.dispatch('getBookmarks')
    },
    feeds (context, route) {
      context.dispatch('getFeeds')
      context.dispatch('getItems')
    }
  }
})
