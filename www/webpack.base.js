const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack')

const pages = [
  "activate",
  "activation",
  "app",
  "appcenter",
  "error",
  "index",
  "login",
  "network",
  "settings",
  "storage",
  "support",
  "updates",
  "backup",
];

module.exports = {
  entry: pages.reduce((jsobject, page) => 
    ( jsobject[page] = `./js/${page}.js`, jsobject ),
    {} // jsobject
  ),
  plugins: pages.map(page =>
    new HtmlWebpackPlugin({ 
      	template: `./${page}.html`,
      	inject: 'body', 
      	chunks: [page], 
      	filename: `${page}.html`
    }))
}
