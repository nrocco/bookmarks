import Vue from 'vue'
import Router from 'vue-router'
import BookmarkList from '@/components/BookmarkList'
import FeedList from '@/components/FeedList'

Vue.use(Router)

export default new Router({
  routes: [
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
})
