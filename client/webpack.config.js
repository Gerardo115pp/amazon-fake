const path = require('path');
const htmlWebpackPlugin = require('html-webpack-plugin');



module.exports = {
	entry: './src/index.js',
	output: {
		path: path.resolve(__dirname, 'build'),
		filename: 'boundle.js'
	},
	devServer: {
		host: "127.0.0.1",
		port: 5001,
		hot: true,
		contentBase: "./public",
		inline: true,
		disableHostCheck: true,
		historyApiFallback: true
	},
	resolve: {
		extensions: ['*', '.mjs', '.js', '.svelte']
	},
	module: {
		rules: [
			{
				test:  /\.js?$/,
				exclude: /node_modules/,
				use: {
					loader: 'babel-loader'
				}
			},
			{
				test: /\.svelte$/,
				use: {
					loader: 'svelte-loader'
				}
			},
			{
				test: /\.svg$/,
				exclude: /node_modules/,
				use: {
					loader: 'svg-inline-loader',
					options: {
					  removeSVGTagAttrs: true
					}
				}
			}
		]
	},
	plugins: [
		new htmlWebpackPlugin({
			inject: true,
			template: './public/index.html',
			filename: './index.html'
		})
	]
}
