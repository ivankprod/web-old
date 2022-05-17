const path = require('path');
const HTMLWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");
const TerserPlugin = require('terser-webpack-plugin');
const { SubresourceIntegrityPlugin } = require('webpack-subresource-integrity');
const babelConfig = require('./babel.config');
const postcssConfig = require('./postcss.config');

const MODE  = process.argv[process.argv.indexOf('--mode') + 1];
const isDEV = MODE === 'development', isPROD = !isDEV;

module.exports = {
	context: path.resolve(__dirname, './src'),

	entry: {
		app:   './scripts/app.js',
		admin: './scripts/admin.js'
	},

	output: {
		path: path.resolve(__dirname, '../../build'),
		filename: 'static/js/[name].[contenthash].js',
		crossOriginLoading: 'anonymous'
	},

	optimization: isPROD ? {
		realContentHash: true,
		minimize: true,
		minimizer: [new CssMinimizerPlugin(), new TerserPlugin()]
	} : {
		realContentHash: true,
		minimize: false
	},

	cache: false,

	plugins: [
		new CopyPlugin({
			patterns: [
				{
					from: path.resolve(__dirname, './sitemap.xml'),
					to:   './sitemap.xml'
				},
				{
					from: path.resolve(__dirname, './src/favicon.ico'),
					to:   './favicon.ico'
				},
				{
					from: path.resolve(__dirname, './src/images'),
					to:   './static/images'
				},
				{
					from: path.resolve(__dirname, './src/fonts'),
					to:   './static/fonts'
				}
			]
		}),
		new SubresourceIntegrityPlugin({ enabled: true, hashFuncNames: ["sha256"] }),
		new HTMLWebpackPlugin({
			chunks: ['app'],
			template: path.resolve(__dirname, '../server/views/partials/footer.hbs'),
			filename: './views/partials/footer.hbs',
			publicPath: '/',
			inject: false,
			scriptLoading: 'defer',
			minify: false
		}),
		new HTMLWebpackPlugin({
			chunks: ['app'],
			template: path.resolve(__dirname, '../server/views/partials/header.hbs'),
			filename: './views/partials/header.hbs',
			publicPath: '/',
			inject: false,
			scriptLoading: 'defer',
			minify: false
		}),
		new CleanWebpackPlugin({ cleanOnceBeforeBuildPatterns: ['./static/js/*', './static/css/*'] }),
		new MiniCssExtractPlugin({ filename: 'static/css/[name].[contenthash].css', ignoreOrder: false })
	],

	module: {
		rules: [
			{
				test: /\.css$/,
				use: [
					MiniCssExtractPlugin.loader,
					{
						loader: 'css-loader',
						options: { url: false, importLoaders: 1 }
					},
					{
						loader: 'postcss-loader',
						options: { postcssOptions: postcssConfig }
					}
				]
			},
			{
				test: /\.js$/,
				exclude: /node_modules/,
				use: {
					loader: 'babel-loader',
					options: babelConfig
				}
			}
		]
	}
}