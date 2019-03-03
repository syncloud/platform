const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack')

const entries = [
    "activation",
//  "app",
//  "appcenter",
  "error",
  "index",
  "login",
//  "network",
  "settings",
//  "storage",
//  "support",
//  "updates",
];

module.exports = {
  entry: entries.reduce((map, key) => 
    ( map[key] = `./js/${key}.js`, map ),
    {}
  ),
  plugins: entries.map(entry =>
    new HtmlWebpackPlugin({ 
      	template: `./${entry}.html`,
      	inject: 'body', 
      	chunks: [entry], 
      	filename: `${entry}.html`
    })).concat([
    new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery"
    })
  ])
}
