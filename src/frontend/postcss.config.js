const MODE  = process.argv[process.argv.indexOf('--mode') + 1];
const isDEV = MODE === 'development';

module.exports = {
	plugins: [
		require('autoprefixer')({ remove: false }),
		require('postcss-preset-env')({ browsers: isDEV ? 'last 1 chrome version, last 1 firefox version' : 'since 2015, not ie <= 11' }),
		require('postcss-import')
	]
};