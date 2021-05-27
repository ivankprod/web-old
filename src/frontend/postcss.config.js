const MODE  = process.argv[process.argv.indexOf('--mode') + 1];
const isDEV = MODE === 'development';

module.exports = {
	plugins: [
		require('autoprefixer')({ remove: false }),
		require('postcss-preset-env'),
		require('postcss-import')
	]
};