import Vue from 'vue'
import Router from 'vue-router'

import Base from '@/components/Base'
import Login from '@/components/Login'

import BookmarkList from '@/components/BookmarkList'
import FeedList from '@/components/FeedList'
import ThoughtList from '@/components/ThoughtList'

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
