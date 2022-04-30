import utils from "./utils.js";

afterEach(() => {
	jest.useRealTimers();
});

test('utils.sleep to be called', () => {
	jest.useFakeTimers();
	jest.spyOn(global, 'setTimeout');

	const callback = jest.fn();

	utils.sleep(10).then(() => { callback(); });
	return expect(setTimeout).toHaveBeenLastCalledWith(expect.any(Function), 10);
});

describe('utils.onScrollPB', () => {
	test('to be undefined', () => {
		utils.onScrollPB();

		expect(document.getElementById('progress-bar')).toBeNull();
	});

	test('to position absolute', () => {
		document.body.innerHTML = '<div id="progress-bar"></div>';

		utils.onScrollPB();

		expect(document.getElementById('progress-bar').style.position).toEqual('absolute');
	});

	test('to position fixed', () => {
		document.body.innerHTML = '<div id="progress-bar"></div>';
		document.documentElement.scrollTop = 25;

		utils.onScrollPB();

		expect(document.getElementById('progress-bar').style.position).toEqual('fixed');
	});
});

test('utils.rewriteMetas to throw', () => {
	expect(utils.rewriteMetas).toThrow("Options not specified!");
});
