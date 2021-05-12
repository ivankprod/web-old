/* 
    MAIN SCRIPT

    Author : IvanK Production
*/

import '../styles/bundle.css';

import utils from './utils.js';
import spa from './spa.js';

//  Elements
let elemSubnavWrappers = document.querySelectorAll('.subnav-container');

//
//  MAIN PAGE EVENTS:
//

//  DOMContentLoaded
document.addEventListener("DOMContentLoaded", function() {
	spa.init();

	//  Submenus
	document.querySelectorAll('a.subnav').forEach(function(elem, i) {
		elem.addEventListener('mouseover', function() {
			elemSubnavWrappers[i].classList.add('showed');
			utils.drawOpacity(elemSubnavWrappers[i], 1);
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
};

//  onPopState
window.onpopstate = spa.popstate;

//  Buttons onClick
[...document.getElementsByTagName('button')].forEach(elem => {
	elem.addEventListener('click', function(e) {
		let dest   = e.target.dataset.href;
		let params = {}

		if (dest.indexOf('?') !== -1) {
			const arr = dest.split('?');

			dest   = arr[0];
			params = utils.queryParse(arr[1]);
		}

		spa.loadPage(dest, params, true);
	});
});

//  Links onClick
document.querySelectorAll('a.spa').forEach(elem => {
	elem.addEventListener('click', function(e) {
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
	});
});
