/* 
    AJAX SCRIPT

    Author : IvanK Production
*/

import { sleep, animate, makeLinear, drawOpacity, queryStringify } from './utils.js';

let elemMasterContainer = document.getElementById('master-container');
let ajaxController;

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
		} else { return { error: { error_code: req.status, error_desc: req.statusText } }; }
	} catch (err) {
		if (err.name !== 'AbortError') {
			return { error: { error_code: 500, error_desc: err.message } };
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
	newAjax, ajaxDone, ajaxErr
}
