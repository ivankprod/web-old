const utils = require('./utils.js');

test('utils.rewriteMetas to throw', () => {
	expect(utils.rewriteMetas).toThrow("Options not specified!");
});
