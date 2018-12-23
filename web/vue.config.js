module.exports = {
  productionSourceMap: false,
  devServer: {
    proxy: 'http://backend:3000'
  }
}
