/* 
    MAIN SCRIPT

    Author : IvanK Production
*/

import '../styles/bundle.css';

import utils from './utils.js';

//  Elements
let elemSubnavWrappers = document.querySelectorAll('.subnav-container');

//  Animations: Animate function
function animate(opts) {
	let start = performance.now();

	requestAnimationFrame(function animate(time) {
		let timeFraction = (time - start) / opts.duration;
		if (timeFraction > 1) timeFraction = 1;

		let progress = opts.timing(timeFraction);

		if (opts.draw) opts.draw(progress);
		if (opts.move) opts.move(progress);
		
		if (timeFraction < 1) { requestAnimationFrame(animate); }
		else { opts.callback(); }
	});
}

//  Animations: linear
function makeLinear(timeFraction) {
	return timeFraction;
}

//  Animations: EaseInOut
function makeEaseInOut(timing) {
	return function(timeFraction) {
		if (timeFraction < 0.5) return timing(2 * timeFraction) / 2;
		else return (2 - timing(2 * (1 - timeFraction))) / 2;
	};
}

//  Animations: Complete timing function
const makeLinearEaseInOut = makeEaseInOut(makeLinear);

//  Animations: Opacity 
function drawOpacity(elem, value, grd = 100) {
	elem.style['-webkit-opacity'] = value / 100 * grd;
	elem.style['-khtml-opacity']  = value / 100 * grd;
	elem.style['-moz-opacity']    = value / 100 * grd;
	elem.style['opacity']         = value / 100 * grd;
}

//
//  MAIN PAGE EVENTS:
//

//  DOMContentLoaded
document.addEventListener("DOMContentLoaded", function() {
	//  Submenus
	document.querySelectorAll('a.subnav').forEach(function(elem, i) {
		elem.addEventListener('mouseover', function() {
			elemSubnavWrappers[i].classList.add('showed');
			drawOpacity(elemSubnavWrappers[i], 1);
		});

		elem.addEventListener('mouseleave', function(event) {
			if (event.relatedTarget != elemSubnavWrappers[i]) {
				elem.classList.remove('hovered');

				animate({
					"duration": 100,
					"timing": makeLinearEaseInOut,
					"draw": function(perc) {
						drawOpacity(elemSubnavWrappers[i], 1 - perc);
					},
					"callback": function() {
						elemSubnavWrappers[i].classList.remove('showed');
					}
				});
			}

			elemSubnavWrappers[i].addEventListener('mouseover', function() {
				this.classList.add('showed'); drawOpacity(this, 1);
				elem.classList.add('hovered');
			});

			elemSubnavWrappers[i].addEventListener('mouseleave', function(_event) {
				if (_event.relatedTarget != event.target) {
					elem.classList.remove('hovered');

					animate({
						"duration": 100,
						"elemw": this,
						"timing": makeLinearEaseInOut,
						"draw": function(perc) {
							drawOpacity(this.elemw, 1 - perc);
						},
						"callback": function() {
							this.elemw.classList.remove('showed');
						}
					});
				}
			});
		});
	});
});

//  onResize
window.onresize = function() {
	//
};

//  onLoad
window.onload = function() {
	document.body.classList.remove('preload');

	clearTimeout(window.tLoader);
	document.getElementById('loader').style.display = 'none';
	document.getElementById('master-container').style.opacity = '1';

	/*utils.rewriteMetas({
		docSource: document,
		docDest:   document,
		metas: [
			'description',
			'og:type'
		]
	});*/
};

// Buttons onClick
[...document.getElementsByTagName('button')].forEach(elem => {
	elem.addEventListener('click', function(e) { window.location.href = e.target.dataset.href; });
});