module.exports = {
  productionSourceMap: false,
  devServer: {
    disableHostCheck: true,
    proxy: 'http://backend:3000'
  }
}
