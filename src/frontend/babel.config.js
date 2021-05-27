const MODE  = process.argv[process.argv.indexOf('--mode') + 1];
const isDEV = MODE === 'development';

module.exports = {
	presets: [
		['@babel/preset-env',
			{
				useBuiltIns: 'usage',
				debug: true,
				corejs: 3
			}
		]
	]
};