/* 
    AJAX SCRIPT

    Author : IvanK Production
*/

import { sleep, onScrollPB, animate, makeLinear, drawOpacity, completeProgress, queryStringify } from './utils.js';

let elemMasterContainer = document.getElementById('master-container');
let ajaxController;

//  Progress Bar class
export class ProgressBar {
	constructor() {
		this.loadFinished = false;

		this.elem    = document.createElement('div');
		this.elem.id = 'progress-bar';

		let elemOld = document.getElementById('progress-bar');
		if (elemOld) {
			elemOld.remove(); this.elem.style.opacity = '1';
			if (window.lastRAF) { cancelAnimationFrame(window.lastRAF); }
		}

		document.body.insertBefore(this.elem, elemMasterContainer);
		onScrollPB();
	}

	start() {
		animate({
			"stoppable": true,
			"duration": 10000,
			"timing": makeLinear,
			"draw": perc => {
				completeProgress(this.elem, 0, perc * 100);
			}
		});

		sleep(1000).then(() => {
			if (!this.loadFinished) { this.elem.style.opacity = '1'; }
		});
	}

	finish() {
		if (window.lastRAF) { cancelAnimationFrame(window.lastRAF); }

		animate({
			"elem": this.elem,
			"duration": 400,
			"timing": makeLinear,
			"draw": function(perc) {
				completeProgress(this.elem, parseInt(this.elem.style.width), perc * 100);
			},
			"callback": () => {
				this.loadFinished = true;
				this.elem.style.opacity = '0';

				sleep(400).then(() => { this.elem.remove(); });
			}
		});
	}
}

//  AJAX window class
class AjaxWindow {
	constructor(type, caption, code, message) {
		this.wndShowTime = 4000;

		this.elemWrapper = document.createElement('div');
		this.elemCaption = document.createElement('div');
		this.elemCode    = document.createElement('div');
		this.elemMessage = document.createElement('div');

		this.elemWrapper.id = 'ajax-info';
		this.elemWrapper.classList.add(type, 'animate-fadein-css');
		this.elemCaption.id = 'ajax-info-caption';
		this.elemCaption.innerHTML = caption;
		this.elemCode.id = 'ajax-info-code';
		this.elemCode.innerHTML = code;
		this.elemMessage.id = 'ajax-info-message';
		this.elemMessage.innerHTML = message;

		this.elemWrapper.append(this.elemCaption, this.elemCode, this.elemMessage);
		elemMasterContainer.append(this.elemWrapper);
	}

	showWindow() {
		this.elemWrapper.style.display = 'block';

		sleep(this.wndShowTime).then(() => { this.closeWindow(); });
	}

	closeWindow() {
		animate({
			duration: 600,
			timing:   makeLinear,
			elem:     this.elemWrapper,
			draw:     function(perc) { drawOpacity(this.elem, 1 - perc); },
			callback: function() {
				this.elem.style.display = 'none';
				this.elem.style.opacity = '1';
			}
		});
	}

	terminate() { this.elemWrapper.remove(); }
}

//  AJAX function
export async function newAjax(url, params = {}, type = 'json') {
	if (ajaxController) ajaxController.abort();

	ajaxController   = new AbortController();
	const ajaxSignal = ajaxController.signal;

	//params['r'] = Math.floor(Math.random() * (1000 - 1) + 1);

	try {
		let req = await fetch(url + queryStringify(params), { ajaxSignal });

		if (req.ok) {
			ajaxController = null;

			return (type == 'json' ? await req.json() : await req.text());
		} else {
			const serverResponse = await req.text();

			return { error: { error_code: req.status, error_desc: req.statusText, error_type: 'server' }, response: serverResponse };
		}
	} catch (err) {
		if (err.name !== 'AbortError') {
			return { error: { error_code: 500, error_desc: err.message, error_type: 'client' }, response: null };
		} else {
			return { error: { error_code: 409, error_desc: 'request aborted', error_type: 'aborted' }, response: null };
		}
	}
}

//  AJAX onsuccess
export function ajaxDone(message, subject = '&nbsp;') {
	if (window.AjaxWindow) { window.AjaxWindow.terminate(); window.AjaxWindow = null; }
	window.AjaxWindow = new AjaxWindow('success', subject, 'Выполнено', message);
	window.AjaxWindow.showWindow();

	return message;
}

//  AJAX onerror
export function ajaxErr(status, message, subject = '&nbsp;') {
	if (window.AjaxWindow) { window.AjaxWindow.terminate(); window.AjaxWindow = null; }
	window.AjaxWindow = new AjaxWindow('error', subject, 'Ошибка ' + status, message);
	window.AjaxWindow.showWindow();

	return message;
}

//  Exports
export default {
	ProgressBar, newAjax, ajaxDone, ajaxErr
}
