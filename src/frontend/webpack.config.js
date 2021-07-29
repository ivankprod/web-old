const fs = require('fs');
const path = require('path');
const HTMLWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyPlugin = require('copy-webpack-plugin');
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');
const SitemapPlugin = require('sitemap-webpack-plugin').default;
const SriPlugin = require('webpack-subresource-integrity');
const DotEnv = require('dotenv-webpack');
const babelConfig = require('./babel.config');
const postcssConfig = require('./postcss.config');

const MODE  = process.argv[process.argv.indexOf('--mode') + 1];
const isDEV = MODE === 'development', isPROD = !isDEV;

const sitemapPaths = JSON.parse(fs.readFileSync(path.resolve(__dirname, '../server/misc/sitemap.json')));

module.exports = {
	context: path.resolve(__dirname, './src'),
	entry: {
		app:   './scripts/app.js',
		admin: './scripts/admin.js'
	},
	output: {
		path: path.resolve(__dirname, '../../build_' + (isDEV ? 'dev' : 'prod')),
		filename: 'static/js/[name].[contenthash].js',
		crossOriginLoading: 'anonymous'
	},
	optimization: isPROD ? { minimizer: [new OptimizeCssAssetsPlugin(), new TerserPlugin()] } : { minimize: false },
	cache: false,
	plugins: [
		new SriPlugin({ hashFuncNames: ['sha256'] }),
		new SitemapPlugin({
			base: 'https://ivankprod.ru',
			paths: sitemapPaths,
			options: { filename: './sitemap.xml', skipgzip: true, lastmod: (new Date()).toDateString() }
		}),
		new HTMLWebpackPlugin({
			chunks: ['app'],
			template: path.resolve(__dirname, './src/views/partials/footer.hbs'),
			filename: path.resolve(__dirname, '../server/views/partials/footer.hbs'),
			publicPath: '/',
			inject: false,
			scriptLoading: 'defer',
			minify: false
		}),
		new HTMLWebpackPlugin({
			chunks: ['app'],
			template: path.resolve(__dirname, './src/views/partials/header.hbs'),
			filename: path.resolve(__dirname, '../server/views/partials/header.hbs'),
			publicPath: '/',
			inject: false,
			scriptLoading: 'defer',
			minify: false
		}),
		new CleanWebpackPlugin({ cleanOnceBeforeBuildPatterns: ['./static/js/*', './static/css/*'] }),
		new MiniCssExtractPlugin({ filename: 'static/css/[name].[contenthash].css', ignoreOrder: false }),
		new CopyPlugin({
			patterns: [
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
				},
				{
					from: path.resolve(__dirname, '../server/certs/ivankprod.ru/'),
					to:   './certs/ivankprod.ru/'
				},
				// {
				// 	from: path.resolve(__dirname, '../server/views'),
				// 	to:   './views',
				// 	globOptions: {
				// 		ignore: [
				// 			'**/footer.hbs',
				// 			'**/header.hbs'
				// 		]
				// 	}
				// }
			]
		}),
		new DotEnv({ path: path.resolve(__dirname, '../../.env') })
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