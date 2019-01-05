import client from '@/client'

const state = {
  bookmarks: [],
  filter: '',
  archived: false
}

const getters = {
  bookmarks: state => {
    return state.bookmarks
  },
  archived: state => {
    return state.archived
  },
  filter: state => {
    return state.filter
  }
}

const actions = {
  getBookmarks (context) {
    let payload = {}

    if (context.getters.filter !== '') {
      payload.q = context.getters.filter
    }

    if (context.getters.archived) {
      payload.archived = 'true'
    }

    client.get(`/bookmarks`, { params: payload }).then(response => {
      context.commit('bookmarks', response.data)
    })
  },
  readLaterBookmark (context, bookmark) {
    client.patch(`/bookmarks/${bookmark.ID}`, { Archived: false }).then(response => {
      context.commit('removeBookmark', bookmark)
    })
  },
  archiveBookmark (context, bookmark) {
    client.patch(`/bookmarks/${bookmark.ID}`, { Archived: true }).then(response => {
      context.commit('removeBookmark', bookmark)
    })
  },
  removeBookmark (context, bookmark) {
    client.delete(`/bookmarks/${bookmark.ID}`).then(response => {
      context.commit('removeBookmark', bookmark)
    })
  }
}

const mutations = {
  bookmarks (state, bookmarks) {
    state.bookmarks = bookmarks
  },
  filter (state, filter) {
    state.filter = filter
  },
  archived (state, archived) {
    state.archived = archived
  },
  removeBookmark (state, bookmark) {
    state.bookmarks.splice(state.bookmarks.indexOf(bookmark), 1)
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
