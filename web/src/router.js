import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

export default new VueRouter({
  mode: 'hash',
  base: process.env.BASE_URL,
  routes: [
    {
      name: 'login',
      path: '/login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/',
      component: () => import('@/views/Base.vue'),
      children: [
        {
          name: 'bookmarks',
          path: '/',
          component: () => import('@/views/BookmarkList.vue'),
          meta: {
            title: 'Bookmarks',
            subtitle: 'All articles you bookmarked',
            color: 'is-primary'
          }
        },
        {
          name: 'feeds',
          path: '/feeds',
          component: () => import('@/views/FeedList.vue'),
          meta: {
            title: 'Feeds',
            subtitle: 'All your rss and atom feeds',
            color: 'is-warning'
          }
        },
        {
          name: 'thoughts',
          path: '/thoughts/:title?',
          component: () => import('@/views/ThoughtList.vue'),
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
