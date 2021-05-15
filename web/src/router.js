import Vue from 'vue'
import Router from 'vue-router'

import Base from '@/views/Base'

import BookmarkList from '@/views/BookmarkList'
import FeedList from '@/views/FeedList'
import Login from '@/views/Login'
import ThoughtList from '@/views/ThoughtList'

Vue.use(Router)

export default new Router({
  mode: 'hash',
  base: process.env.BASE_URL,
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
          name: 'bookmarks',
          path: '/',
          component: BookmarkList,
          meta: {
            title: 'Bookmarks',
            subtitle: 'All articles you bookmarked',
            color: 'is-primary'
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
        },
        {
          name: 'thoughts',
          path: '/thoughts/:title?',
          component: ThoughtList,
          meta: {
            title: 'Thoughts',
            subtitle: 'All your thoughts',
            color: 'is-info'
          }
        }
      ]
    }
  ]
})
