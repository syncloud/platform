const api = require('./src/stub/api')

module.exports = {
  devServer: {
    before: api.mock
    // proxy: 'http://localhost:8081'
  }
}
