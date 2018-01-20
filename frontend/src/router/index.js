import Vue from 'vue'
import Router from 'vue-router'

import Base from '@/components/Base'
import Login from '@/components/Login'

import BookmarkList from '@/components/BookmarkList'
import FeedList from '@/components/FeedList'

Vue.use(Router)

export default new Router({
  routes: [
    {
      name: 'login',
      path: '/login',
      component: Login
    },
    {
      path: '/',
      component: Base,
      children: [
        {
          name: 'readitlater',
          path: '/',
          component: BookmarkList,
          meta: {
            title: 'Read it later',
            subtitle: 'All articles you recently saved',
            color: 'is-primary'
          }
        },
        {
          name: 'archive',
          path: '/archive',
          component: BookmarkList,
          meta: {
            title: 'Archive',
            subtitle: 'All articles you archived',
            color: 'is-dark'
          }
        },
        {
          name: 'feeds',
          path: '/feeds',
          component: FeedList,
          meta: {
            title: 'Feeds',
            subtitle: 'All your rss and atom feeds',
            color: 'is-warning'
          }
        }
      ]
    }
  ]
})
