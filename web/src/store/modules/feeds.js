import client from '@/client'

const state = {
  feeds: [],
  items: [],
  selectedFeed: null
}

const getters = {
  feeds: state => {
    return state.feeds
  },
  items: state => {
    return state.items
  },
  categories: state => {
    let categories = []
    for (let feed of state.feeds) {
      categories.push(feed.Category)
    }
    return [...new Set(categories)]
  },
  selectedFeed: state => {
    return state.selectedFeed
  }
}

const actions = {
  getFeeds (context) {
    client.get(`/feeds`).then(response => {
      context.commit('feeds', response.data)
    })
  },
  addFeed (context, feed) {
    console.log('POST /feeds', feed)
    client.post(`/feeds`, feed).then(response => {
      context.commit('addFeed', response.data)
      context.commit('selectedFeed', response.data)
    })
  },
  selectFeed (context, feed) {
    context.commit('selectedFeed', feed)
    client.get(`/items`, {params: {feed: feed.ID}}).then(response => {
      context.commit('items', response.data)
    })
  },
  refreshFeed (context, feed) {
    client.post(`/feeds/${feed.ID}/refresh`).then(response => {
      // TODO sleep a second
      // TODO load items again
    })
  },
  deleteFeed (context, feed) {
    client.delete(`/feeds/${feed.ID}`).then(response => {
      context.commit('removeFeed', feed)
    })
  },
  readLaterFeedItem (context, item) {
    client.post(`/items/${item.ID}/readitlater`).then(response => {
      context.commit('removeItem', item)
    })
  },
  removeFeedItem (context, item) {
    client.delete(`/items/${item.ID}`).then(response => {
      context.commit('removeItem', item)
    })
  }
}

const mutations = {
  feeds (state, feeds) {
    state.feeds = feeds
  },
  items (state, items) {
    state.items = items
  },
  addFeed (state, feed) {
    state.feeds.push(feed)
  },
  selectedFeed (state, feed) {
    state.selectedFeed = feed
  },
  removeFeed (state, feed) {
    if (state.selectedFeed === feed) {
      state.selectedFeed = null
    }
    state.feeds.splice(state.feeds.indexOf(feed), 1)
  },
  removeItem (state, item) {
    state.items.splice(state.items.indexOf(item), 1)

    let feed = state.feeds.find(feed => {
      return feed.ID === item.FeedID
    })

    if (feed) {
      feed.Items -= 1
    }
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
