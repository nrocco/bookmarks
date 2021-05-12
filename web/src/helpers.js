export default {
  beforeRouteUpdate (to, from, next) {
    next()
    this.onLoad(sanitiseFilters(to.query))
  },

  beforeRouteEnter (to, from, next) {
    next(vm => vm.onLoad(sanitiseFilters(to.query)))
  },

  methods: {
    changeRouteOnFilterChange (filters, path) {
      this.$router.push({ path: path, query: sanitiseFilters(filters) })
    }
  }
}

function sanitiseFilters (parameters) {
  parameters = JSON.parse(JSON.stringify(parameters))
  Object.keys(parameters).forEach((key) => (!parameters[key]) && delete parameters[key])

  return parameters
}
