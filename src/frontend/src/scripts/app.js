/* 
    MAIN SCRIPT

    Author : IvanK Production
*/

//  import styles
import 'normalize.css';
import '../styles/bundle.css';

//  import app modules
import utils, { onScrollPB } from './utils.js';
import spa from './spa.js';

//  import components
import Slider from './components/slider.js'

//  Slider: config
const sliderConfig = {
	speed: 600,

	autoplay: { delay: 6000 },

	pagination: {
		element: '.slider-pagination',
		clickable: true
	}
};

//  Slider: init
window.slider = new Slider('slider-container', sliderConfig);

//  Slider: onSliderTransitionStart
const onSliderTransitionStart = function(slider) {
	const slide = slider.slides[slider.activeIndex].children[0];

	slide.style.opacity = 0;

	slide.children[0].style.opacity   = '0';
	slide.children[0].style.animation = 'none';
	slide.children[2].style.opacity   = '0';
	slide.children[2].style.animation = 'none';
}

//  Slider: onSliderTransitionEnd
const onSliderTransitionEnd = function(slider) {
	const slide = slider.slides[slider.activeIndex].children[0];

	utils.animate({
		duration: 800,
		timing:   utils.makeLinear,
		elem:     slide,
		draw:     function(perc) { utils.drawOpacity(this.elem, perc); },
		callback: function() {
			this.elem.children[2].style.opacity   = '1';
			this.elem.children[2].style.animation = 'slideIn 1000ms cubic-bezier(0.190, 1.000, 0.220, 1.000), fadeIn 400ms linear';
		}
	});

	slide.children[0].style.opacity   = '1';
	slide.children[0].style.animation = 'slideBlockTitle 1600ms cubic-bezier(0.190, 1.000, 0.220, 1.000), fadeIn 600ms linear';
}

//  Slider: attach events
window.slider.on('transitionStart', onSliderTransitionStart);
window.slider.on('transitionEnd',   onSliderTransitionEnd);

//  Elements
let elemSubnavWrappers = document.querySelectorAll('.subnav-container');

//  DOMContentLoaded
document.addEventListener("DOMContentLoaded", function() {
	spa.init();

	//  Submenus
	document.querySelectorAll('a.subnav').forEach(function(elem, i) {
		elem.addEventListener('mouseover', function() {
			elemSubnavWrappers[i].classList.add('showed');
			utils.drawOpacity(elemSubnavWrappers[i], 1);
		});

		elem.addEventListener('mousemove', function() {
			if (!elemSubnavWrappers[i].classList.contains('showed')) {
				elem.dispatchEvent(new Event('mouseover'));
			}
		});

		elem.addEventListener('mouseleave', function(event) {
			if (event.relatedTarget != elemSubnavWrappers[i]) {
				elem.classList.remove('hovered');

				utils.animate({
					"duration": 100,
					"timing": utils.makeLinearEaseInOut,
					"draw": perc => {
						utils.drawOpacity(elemSubnavWrappers[i], 1 - perc);
					},
					"callback": () => {
						elemSubnavWrappers[i].classList.remove('showed');
					}
				});
			}

			elemSubnavWrappers[i].addEventListener('mouseover', function() {
				this.classList.add('showed'); utils.drawOpacity(this, 1);
				elem.classList.add('hovered');
			});

			elemSubnavWrappers[i].addEventListener('mouseleave', function(_event) {
				if (_event.relatedTarget != event.target) {
					elem.classList.remove('hovered');

					utils.animate({
						"duration": 100,
						"elemw": this,
						"timing": utils.makeLinearEaseInOut,
						"draw": function(perc) {
							utils.drawOpacity(this.elemw, 1 - perc);
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

//  Buttons onClick
function fillButtonsOnClick() {
	document.querySelectorAll('button').forEach(elem => {
		if (elem.classList.contains('spa')) {
			elem.onclick = function(e) {
				let dest   = e.target.dataset.href;
				let params = {}

				if (dest.indexOf('?') !== -1) {
					const arr = dest.split('?');

					dest   = arr[0];
					params = utils.queryParse(arr[1]);
				}

				spa.loadPage(dest, params, true);
			};
		} else { elem.onclick = function(e) { window.location.href = this.dataset.href; }; }
	});
}

//  Links onClick
function fillLinksOnClick() {
	document.querySelectorAll('a.spa').forEach(elem => {
		elem.onclick = function(e) {
			e.preventDefault();

			let dest   = elem.getAttribute('href');
			let params = {}

			if (dest.indexOf('?') !== -1) {
				const arr = dest.split('?');

				dest   = arr[0];
				params = utils.queryParse(arr[1]);
			}

			spa.loadPage(dest, params, true);

			return false;
		};
	});
}

//  onResize
window.onresize = function() {
	//
};

//  onScroll
window.onscroll = function() {
	onScrollPB();
};

//  onLoad
window.onload = function() {
	document.body.classList.remove('preload');

	clearTimeout(window.tLoader);
	document.getElementById('loader').style.display = 'none';
	document.getElementById('master-container').style.opacity = '1';

	fillButtonsOnClick();
	fillLinksOnClick();
};

//  onSPAPageLoaded
window.onPageLoaded = dataExtras => {
	fillButtonsOnClick();
	fillLinksOnClick();

	// Slider reinit
	if (document.querySelector('.slider-section')) {
		window.slider = new Slider('slider-container', sliderConfig);

		window.slider.on('transitionStart', onSliderTransitionStart);
		window.slider.on('transitionEnd',   onSliderTransitionEnd);
	}
}

//  onPopState
window.onpopstate = spa.popstate;
