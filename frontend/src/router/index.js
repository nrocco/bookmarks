import Vue from 'vue'
import Router from 'vue-router'
import Base from '@/components/Base'
import BookmarkList from '@/components/BookmarkList'
import FeedList from '@/components/FeedList'
import Login from '@/components/Login'

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
          path: '/',
          component: BookmarkList,
          meta: {
            title: 'Read it later',
            subtitle: 'All articles you recently saved',
            archived: false,
            color: 'is-primary'
          }
        },
        {
          path: '/archive',
          component: BookmarkList,
          meta: {
            title: 'Archive',
            subtitle: 'All articles you archived',
            archived: true,
            color: 'is-dark'
          }
        },
        {
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
