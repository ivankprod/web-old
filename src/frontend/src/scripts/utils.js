/* 
    UTILS SCRIPT

    Author : IvanK Production
*/

////////////////////
//  MAIN SECTION  //
////////////////////

//  sleep
export function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

//  ProgressBar onScroll
export function onScrollPB() {
	let elemProgressBar = document.getElementById('progress-bar');
	if (!elemProgressBar) return;

	const scrolled = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop;

	if (scrolled > 24) {
		elemProgressBar.style.position = 'fixed';
		elemProgressBar.style.top      = '0';
	} else {
		elemProgressBar.style.position = 'absolute';
		elemProgressBar.style.top      = '24px';
	}
}

//////////////////////////
//  ANIMATIONS SECTION  //
//////////////////////////

//  Animations: main function
export function animate(opts) {
	let start = performance.now();

	if (opts.stoppable) { window.lastRAF = null; }

	requestAnimationFrame(function animate(time) {
		let timeFraction = (time - start) / opts.duration;
		if (timeFraction > 1) timeFraction = 1;
		if (timeFraction < 0) timeFraction = 0;

		let progress = opts.timing(timeFraction);

		if (opts.draw) { opts.draw(progress); }
		if (opts.move) { opts.move(progress); }

		if (timeFraction < 1) {
			if (opts.stoppable) { window.lastRAF = requestAnimationFrame(animate); }
			else { requestAnimationFrame(animate); }
		} else {
			if (opts.callback) { opts.callback(); }
		}
	});
}

//  Animations: linear
export function makeLinear(timeFraction) {
	return timeFraction;
}

//  Animations: pow
export function makePow(timeFraction) {
	return Math.pow(timeFraction, 5);
}

//  Animations: circ
export function makeCirc(timeFraction) {
	return 1 - Math.sin(Math.acos(timeFraction));
}

//  Animations: EaseOut
export function makeEaseOut(timing) {
	return function(timeFraction) {
		return 1 - timing(1 - timeFraction);
	};
}

//  Animations: EaseInOut
export function makeEaseInOut(timing) {
	return function(timeFraction) {
		if (timeFraction < 0.5) return timing(2 * timeFraction) / 2;
		else return (2 - timing(2 * (1 - timeFraction))) / 2;
	};
}

//  Animations: complete timing functions
export var makeLinearEaseInOut = makeEaseInOut(makeLinear);
export var makePowEaseOut      = makeEaseOut(makePow);
export var makeCircEaseInOut   = makeEaseInOut(makeCirc);

//  Animations: opacity 
export function drawOpacity(elem, value) {
	elem.style.opacity = value;
}

//  Animations: async opacity (fadeout)
export async function fadeOut(elem) {
	elem.style.opacity = '0';
};

//  Animations: height
export function drawHeight(elem, value) {
	elem.style.height = value + 'px';
}

//  Animations: move by pixel parameter
export function movePX(elem, style, value) {
	elem.style[style] = value + 'px';
}

//  Animations: complete progress bar animation
export function completeProgress(elem, start, value) {
	elem.style.width = (value > start ? value : start) + '%';
}

/////////////////////
//  QUERY SECTION  //
/////////////////////

//  Query object to string
export function queryStringify(obj) {
	let params = [];

	Object.keys(obj).forEach(key => {
		if (obj[key] !== '') params.push(String(key + '=' + obj[key]).replace(/\s/g, '_'));
	});

	return (params.length ? '?' + params.join('&') : '');
}

//  Query string to object
export function queryParse(str) {
	let result = {};
	if (str == '') return result;

	const obj = new URLSearchParams(str);
	for (const [key, value] of obj.entries()) { result[key] = value; }

	return result;
};

////////////////////
//  META SECTION  //
////////////////////

//  Get canonical link
function getCanonical(doc) {
	let result = doc.querySelector('link[rel="canonical"]');
	if (!result) throw new Error("Canonical link not found!");

	return result.href;
}

//  Set canonical link
function setCanonical(doc, value) {
	let result = doc.querySelector('link[rel="canonical"]');
	if (!result) throw new Error("Canonical link not found!");

	result.href = value;
}

//  Get meta tag
export function getMeta(doc, metaName) {
	let result = doc.querySelector('meta[name="' + metaName + '"]');
	if (!result) result = doc.querySelector('meta[property="' + metaName + '"]');
	if (!result) throw new Error("Meta not found!");

	return result.content;
}

//  Set meta tag
export function setMeta(doc, metaName, metaContent) {
	let result = doc.querySelector('meta[name="' + metaName + '"]');
	if (!result) result = doc.querySelector('meta[property="' + metaName + '"]');
	if (!result) throw new Error("Meta not found!");

	result.content = metaContent;
}

//  Rewrite meta tags from ajax-loaded page to current page
export function rewriteMetas(opts) {
	if (opts && opts.metas && opts.docSource && opts.docDest) {
		opts.metas.forEach(elem => {
			setMeta(opts.docDest, elem, getMeta(opts.docSource, elem));
		});

		if (opts.withCanonical) setCanonical(opts.docDest, getCanonical(opts.docSource));
	} else {
		throw new Error("Options not specified!");
	}
}

//////////////////////
//  EXPORT SECTION  //
//////////////////////

export default {
	animate, makeLinear, makePow, makeCirc, makeEaseOut, makeEaseInOut, makeLinearEaseInOut, makePowEaseOut, makeCircEaseInOut,
	drawOpacity, drawHeight, movePX, fadeOut,

	queryStringify, queryParse, rewriteMetas,

	sleep, onScrollPB
}
