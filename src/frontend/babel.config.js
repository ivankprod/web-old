const MODE  = process.argv[process.argv.indexOf('--mode') + 1];
const isDEV = MODE === 'development';

module.exports = {
	presets: [
		['@babel/preset-env',
			{
				useBuiltIns: 'usage',
				corejs: 3,
				targets: isDEV ? 'last 1 chrome version, last 1 firefox version' : 'since 2015, not ie <= 11'
			}
		]
	]
};