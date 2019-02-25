const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack')

module.exports = {
  entry: {
    index: './js/index.js',
    error: './js/error.js'
  },
  plugins: [
    new HtmlWebpackPlugin({ 
      	template: './index.html',
      	inject: 'body', 
      	chunks: ['index'], 
      	filename: 'index.html' 
    }),
    new HtmlWebpackPlugin({ 
      	template: './error.html', 
      	inject: 'body', 
      	chunks: ['error'], 
      	filename: 'error.html' 
    }),
    new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery"
    })
  ]
}