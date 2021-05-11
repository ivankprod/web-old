/* 
    MAIN SCRIPT

    Author : IvanK Production
*/

import '../styles/bundle.css';

//  Elements
let   elemSubnavWrappers    = document.querySelectorAll('.subnav-container');
/*const elemFooterBlockFirst  = document.querySelector(".footer-block.first");
const elemFooterBlockSecond = document.querySelector(".footer-block.second");
const elemFooterBlockThird  = document.querySelector(".footer-block.third");
const elemFooterBlockHelper = document.querySelector(".footer-block.helper");
let   elemFooterContainer   = elemFooterBlockThird.parentNode;*/

//  Animations: Animate function
function animate(opts) {
	var start = performance.now();

	requestAnimationFrame(function animate(time) {
		var timeFraction = (time - start) / opts.duration;
		if (timeFraction > 1) timeFraction = 1;

		var progress = opts.timing(timeFraction);

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
var makeLinearEaseInOut = makeEaseInOut(makeLinear);

//  Animations: Opacity 
function drawOpacity(elem, value, grd = 100) {
	elem.style['-webkit-opacity'] = value / 100 * grd;
	elem.style['-khtml-opacity']  = value / 100 * grd;
	elem.style['-moz-opacity']    = value / 100 * grd;
	elem.style['opacity']         = value / 100 * grd;
}

//  swapping footer blocks
/*function fixFooterBlocks() {
	if (window.matchMedia("(min-width: 1200px)").matches) {
		elemFooterContainer.append(elemFooterBlockFirst, elemFooterBlockSecond, elemFooterBlockThird, elemFooterBlockHelper);
	} else if (window.matchMedia("(min-width: 992px)").matches) {
		elemFooterContainer.append(elemFooterBlockFirst, elemFooterBlockSecond, elemFooterBlockThird, elemFooterBlockHelper);
	} else if (window.matchMedia("(min-width: 768px)").matches) {
		elemFooterContainer.append(elemFooterBlockSecond, elemFooterBlockThird, elemFooterBlockFirst, elemFooterBlockHelper);
	} else if (window.matchMedia("(min-width: 10px)").matches) {
		elemFooterContainer.append(elemFooterBlockSecond, elemFooterBlockThird, elemFooterBlockFirst, elemFooterBlockHelper);
	}
}*/

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

	//fixFooterBlocks();
});

//  onResize
window.onresize = function() {
	//fixFooterBlocks();
};

//  onLoad
window.onload = function() {
	document.body.classList.remove('preload');

	/* FOR TEST
	setTimeout(function() {*/
	clearTimeout(window.tLoader);
	document.getElementById('loader').style.display = 'none';
	document.getElementById('master-container').style.opacity = '1';
	/*}, 2000);*/
};

// Buttons onClick
[...document.getElementsByTagName('button')].forEach(elem => {
	elem.addEventListener('click', function(e) { window.location.href = e.target.dataset.href; });
});