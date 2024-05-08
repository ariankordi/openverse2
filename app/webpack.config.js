const path = require('path');
const webpack = require('webpack');
const VueLoaderPlugin = require('vue-loader/lib/plugin');

module.exports = {
	entry: path.join(__dirname, './main.js'),
	output: {
		path: path.join(__dirname, '../static'),
		filename: 'openverse.js'
	},
	mode: 'production',
	watch: true,
	module: {
		rules: [
			{
				test: /\.vue$/,
				loader: 'vue-loader'
			},
			{
				test: /\.js$/,
				loader: 'babel-loader',
				exclude: /node_modules/,
				query: {
					presets: ['es2015']
				}
			}
		]
	},
	plugins: [
		new VueLoaderPlugin(),
		new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/)
	]/*,
	resolve: {
		alias: {
			'vue$': 'vue/dist/vue.esm.js' // 'vue/dist/vue.common.js' for webpack 1
		}
	}*/
}
