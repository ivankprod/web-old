const path = require('path');
const SitemapGenerator = require('sitemap-generator');

const generator = SitemapGenerator('https://ivankprod.ru', {
	maxDepth: 0,
	lastMod: true,
	priorityMap: [1.0, 0.8, 0.6, 0.4, 0.2],
	filepath: path.resolve(__dirname, '/sitemap.xml'),
	maxEntriesPerFile: 50000,
	stripQuerystring: true,
	customDomain: 'https://ivankprod.ru',
	ignore: url => { return /(auth|admin|api|legal|join)/g.test(url) }
});

generator.on('error', (error) => {
	console.log('Error during sitemap.xml generation: ', error);
	if (error.code != 404) process.exit(1);
});

generator.on('done', () => {
	console.log(`sitemap.xml has been generated`);
});

generator.start();
