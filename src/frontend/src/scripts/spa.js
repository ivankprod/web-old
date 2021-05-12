/* 
    SPA SCRIPT

    Author : IvanK Production
*/

import { newAjax, ajaxErr } from './ajax.js';
import { sleep, fadeOut, queryParse, queryStringify, getMeta, setMeta, rewriteMetas } from './utils.js';

//  Hostname var
const strServerHost = String('https://' + (process.env.SERVER_HOST != '' ? process.env.SERVER_HOST : 'ivankprod.ru'));

//  HistoryAPI: state
const intHrefStart  = strServerHost.length;
let loc     = window.location.href;
let locHref = loc.substring(intHrefStart + 1, (loc.indexOf('?') !== -1 ? loc.indexOf('?') : loc.length));
let hState  = {
	href:   '/' + locHref,
	params: queryParse(loc.substring(intHrefStart + 1).replace(locHref, '').substring(1)),
	title:  document.title,
	url:    loc.substring(intHrefStart)
};

//  Loads ajax page
export async function loadPage(strHref, params = {}, changeAddress = false, callback = null) {
	const res = await newAjax(strHref, params, 'text');

	if (res.error) {
		ajaxErr(res.error.error_code, res.error.error_desc);
	} else {
		const elemActiveNavItem = document.querySelector('ul.mnav li a.nav-item-active');
		if (elemActiveNavItem) elemActiveNavItem.classList.remove('nav-item-active');

		const oParser    = new DOMParser();
		const oDoc       = oParser.parseFromString(res, 'text/html');
		const newContent = oDoc.getElementById('content');
		let   oldContent = document.getElementById('content');
		document.title   = oDoc.title;

		rewriteMetas({
			docSource: oDoc,
			docDest:   document,
			metas: [
				'robots',
				'og:title', 'og:description', 'og:type', 'og:image', 'og:url', 'og:site_name', 'og:locale',
				'twitter:card', 'twitter:title', 'twitter:description', 'twitter:image'
			],
			withCanonical: true
		});

		const elemDataExtras = oDoc.getElementById('data-extras');
		if (elemDataExtras) dataExtras = JSON.parse(elemDataExtras.textContent);

		hState.href   = strHref;
		hState.params = params;
		hState.title  = oDoc.title;
		hState.url    = strHref + queryStringify(params);

		if (changeAddress) window.history.pushState(hState, hState.title, hState.url);

		console.log(hState);

		fadeOut(oldContent).then(() => {
			sleep(110).then(() => {
				oldContent.parentNode.replaceChild(newContent, oldContent); //window.onPageLoaded();
			});
		});

		const scope = getMeta(oDoc, 'app:scope');
		setMeta(document, 'app:scope', scope);

		let itemActive = document.querySelector('ul.mnav li a[data-scope="' + scope + '"]');
		if (itemActive) itemActive.classList.add('nav-item-active');

		if (callback) callback();
	}
}

//  HistoryAPI: replace state on load
export function init() {
	window.history.replaceState(hState, hState.title, hState.url);
};

//  HistoryAPI: when back or forward
export function popstate(oEvent) {
	loadPage(oEvent.state.href, oEvent.state.params);
};

//  Exports
export default {
	init, popstate, loadPage
}
