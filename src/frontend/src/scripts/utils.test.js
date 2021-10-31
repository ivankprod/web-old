//const utils = require('./utils.js');
import utils from "./utils.js"

test('utils.rewriteMetas to throw', () => {
	expect(utils.rewriteMetas).toThrow("Options not specified!");
});
