/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (function() { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./scripts/ajax.js":
/*!*************************!*\
  !*** ./scripts/ajax.js ***!
  \*************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"newAjax\": function() { return /* binding */ newAjax; },\n/* harmony export */   \"ajaxDone\": function() { return /* binding */ ajaxDone; },\n/* harmony export */   \"ajaxErr\": function() { return /* binding */ ajaxErr; }\n/* harmony export */ });\n/* harmony import */ var _utils_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./utils.js */ \"./scripts/utils.js\");\n/* \r\n    AJAX SCRIPT\r\n\r\n    Author : IvanK Production\r\n*/\n\nlet elemMasterContainer = document.getElementById('master-container');\nlet ajaxController; //  AJAX window class\n\nclass AjaxWindow {\n  constructor(type, caption, code, message) {\n    this.wndShowTime = 4000;\n    this.elemWrapper = document.createElement('div');\n    this.elemCaption = document.createElement('div');\n    this.elemCode = document.createElement('div');\n    this.elemMessage = document.createElement('div');\n    this.elemWrapper.id = 'ajax-info';\n    this.elemWrapper.classList.add(type, 'animate-fadein-css');\n    this.elemCaption.id = 'ajax-info-caption';\n    this.elemCaption.innerHTML = caption;\n    this.elemCode.id = 'ajax-info-code';\n    this.elemCode.innerHTML = code;\n    this.elemMessage.id = 'ajax-info-message';\n    this.elemMessage.innerHTML = message;\n    this.elemWrapper.append(this.elemCaption, this.elemCode, this.elemMessage);\n    elemMasterContainer.append(this.elemWrapper);\n  }\n\n  showWindow() {\n    this.elemWrapper.style.display = 'block';\n    (0,_utils_js__WEBPACK_IMPORTED_MODULE_0__.sleep)(this.wndShowTime).then(() => {\n      this.closeWindow();\n    });\n  }\n\n  closeWindow() {\n    (0,_utils_js__WEBPACK_IMPORTED_MODULE_0__.animate)({\n      duration: 600,\n      timing: _utils_js__WEBPACK_IMPORTED_MODULE_0__.makeLinear,\n      elem: this.elemWrapper,\n      draw: function (perc) {\n        (0,_utils_js__WEBPACK_IMPORTED_MODULE_0__.drawOpacity)(this.elem, 1 - perc);\n      },\n      callback: function () {\n        this.elem.style.display = 'none';\n        this.elem.style.opacity = '1';\n      }\n    });\n  }\n\n  terminate() {\n    this.elemWrapper.remove();\n  }\n\n} //  AJAX function\n\n\nasync function newAjax(url, params = {}, type = 'json') {\n  if (ajaxController) ajaxController.abort();\n  ajaxController = new AbortController();\n  const ajaxSignal = ajaxController.signal; //params['r'] = Math.floor(Math.random() * (1000 - 1) + 1);\n\n  try {\n    let req = await fetch(url + (0,_utils_js__WEBPACK_IMPORTED_MODULE_0__.queryStringify)(params), {\n      ajaxSignal\n    });\n\n    if (req.ok) {\n      ajaxController = null;\n      return type == 'json' ? await req.json() : await req.text();\n    } else {\n      const serverResponse = await req.text();\n      return {\n        error: {\n          error_code: req.status,\n          error_desc: req.statusText,\n          error_type: 'server'\n        },\n        response: serverResponse\n      };\n    }\n  } catch (err) {\n    if (err.name !== 'AbortError') {\n      return {\n        error: {\n          error_code: 500,\n          error_desc: err.message,\n          error_type: 'client'\n        },\n        response: null\n      };\n    }\n  }\n} //  AJAX onsuccess\n\nfunction ajaxDone(message, subject = '&nbsp;') {\n  if (window.AjaxWindow) {\n    window.AjaxWindow.terminate();\n    window.AjaxWindow = null;\n  }\n\n  window.AjaxWindow = new AjaxWindow('success', subject, 'Выполнено', message);\n  window.AjaxWindow.showWindow();\n  return message;\n} //  AJAX onerror\n\nfunction ajaxErr(status, message, subject = '&nbsp;') {\n  if (window.AjaxWindow) {\n    window.AjaxWindow.terminate();\n    window.AjaxWindow = null;\n  }\n\n  window.AjaxWindow = new AjaxWindow('error', subject, 'Ошибка ' + status, message);\n  window.AjaxWindow.showWindow();\n  return message;\n} //  Exports\n\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  newAjax,\n  ajaxDone,\n  ajaxErr\n});\n\n//# sourceURL=webpack:///./scripts/ajax.js?");

/***/ }),

/***/ "./scripts/app.js":
/*!************************!*\
  !*** ./scripts/app.js ***!
  \************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _styles_bundle_css__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../styles/bundle.css */ \"./styles/bundle.css\");\n/* harmony import */ var _utils_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./utils.js */ \"./scripts/utils.js\");\n/* harmony import */ var _spa_js__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./spa.js */ \"./scripts/spa.js\");\n/* \r\n    MAIN SCRIPT\r\n\r\n    Author : IvanK Production\r\n*/\n\n\n //  Elements\n\nlet elemSubnavWrappers = document.querySelectorAll('.subnav-container'); //  DOMContentLoaded\n\ndocument.addEventListener(\"DOMContentLoaded\", function () {\n  _spa_js__WEBPACK_IMPORTED_MODULE_2__.default.init(); //  Submenus\n\n  document.querySelectorAll('a.subnav').forEach(function (elem, i) {\n    elem.addEventListener('mouseover', function () {\n      elemSubnavWrappers[i].classList.add('showed');\n      _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.drawOpacity(elemSubnavWrappers[i], 1);\n    });\n    elem.addEventListener('mouseleave', function (event) {\n      if (event.relatedTarget != elemSubnavWrappers[i]) {\n        elem.classList.remove('hovered');\n        _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.animate({\n          \"duration\": 100,\n          \"timing\": _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.makeLinearEaseInOut,\n          \"draw\": perc => {\n            _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.drawOpacity(elemSubnavWrappers[i], 1 - perc);\n          },\n          \"callback\": () => {\n            elemSubnavWrappers[i].classList.remove('showed');\n          }\n        });\n      }\n\n      elemSubnavWrappers[i].addEventListener('mouseover', function () {\n        this.classList.add('showed');\n        _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.drawOpacity(this, 1);\n        elem.classList.add('hovered');\n      });\n      elemSubnavWrappers[i].addEventListener('mouseleave', function (_event) {\n        if (_event.relatedTarget != event.target) {\n          elem.classList.remove('hovered');\n          _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.animate({\n            \"duration\": 100,\n            \"elemw\": this,\n            \"timing\": _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.makeLinearEaseInOut,\n            \"draw\": function (perc) {\n              _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.drawOpacity(this.elemw, 1 - perc);\n            },\n            \"callback\": function () {\n              this.elemw.classList.remove('showed');\n            }\n          });\n        }\n      });\n    });\n  });\n}); //  Buttons onClick\n\nfunction fillButtonsOnClick() {\n  document.querySelectorAll('button').forEach(elem => {\n    if (elem.classList.contains('spa')) {\n      elem.onclick = function (e) {\n        let dest = e.target.dataset.href;\n        let params = {};\n\n        if (dest.indexOf('?') !== -1) {\n          const arr = dest.split('?');\n          dest = arr[0];\n          params = _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.queryParse(arr[1]);\n        }\n\n        _spa_js__WEBPACK_IMPORTED_MODULE_2__.default.loadPage(dest, params, true);\n      };\n    } else {\n      elem.onclick = function (e) {\n        window.location.href = e.target.dataset.href;\n      };\n    }\n  });\n} //  Links onClick\n\n\nfunction fillLinksOnClick() {\n  document.querySelectorAll('a.spa').forEach(elem => {\n    elem.onclick = function (e) {\n      e.preventDefault();\n      let dest = elem.getAttribute('href');\n      let params = {};\n\n      if (dest.indexOf('?') !== -1) {\n        const arr = dest.split('?');\n        dest = arr[0];\n        params = _utils_js__WEBPACK_IMPORTED_MODULE_1__.default.queryParse(arr[1]);\n      }\n\n      _spa_js__WEBPACK_IMPORTED_MODULE_2__.default.loadPage(dest, params, true);\n      return false;\n    };\n  });\n} //  onResize\n\n\nwindow.onresize = function () {//\n}; //  onLoad\n\n\nwindow.onload = function () {\n  document.body.classList.remove('preload');\n  clearTimeout(window.tLoader);\n  document.getElementById('loader').style.display = 'none';\n  document.getElementById('master-container').style.opacity = '1';\n  fillButtonsOnClick();\n  fillLinksOnClick();\n}; //  onSPAPageLoaded\n\n\nwindow.onPageLoaded = dataExtras => {\n  fillButtonsOnClick();\n  fillLinksOnClick();\n}; //  onPopState\n\n\nwindow.onpopstate = _spa_js__WEBPACK_IMPORTED_MODULE_2__.default.popstate;\n\n//# sourceURL=webpack:///./scripts/app.js?");

/***/ }),

/***/ "./scripts/spa.js":
/*!************************!*\
  !*** ./scripts/spa.js ***!
  \************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"loadPage\": function() { return /* binding */ loadPage; },\n/* harmony export */   \"init\": function() { return /* binding */ init; },\n/* harmony export */   \"popstate\": function() { return /* binding */ popstate; }\n/* harmony export */ });\n/* harmony import */ var _ajax_js__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./ajax.js */ \"./scripts/ajax.js\");\n/* harmony import */ var _utils_js__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./utils.js */ \"./scripts/utils.js\");\n/* \r\n    SPA SCRIPT\r\n\r\n    Author : IvanK Production\r\n*/\n\n //  Hostname var\n\nconst strServerHost = String('https://' + ( true ? \"ivankprod.local\" : 0)); //  Extras Data\n\nlet dataExtras = null; //  HistoryAPI: state\n\nconst intHrefStart = strServerHost.length;\nlet loc = window.location.href;\nlet locHref = loc.substring(intHrefStart + 1, loc.indexOf('?') !== -1 ? loc.indexOf('?') : loc.length);\nlet hState = {\n  href: '/' + locHref,\n  params: (0,_utils_js__WEBPACK_IMPORTED_MODULE_1__.queryParse)(loc.substring(intHrefStart + 1).replace(locHref, '').substring(1)),\n  title: document.title,\n  url: loc.substring(intHrefStart)\n}; //  Loads ajax page\n\nasync function loadPage(strHref, params = {}, changeAddress = false, callback = null) {\n  let res = await (0,_ajax_js__WEBPACK_IMPORTED_MODULE_0__.newAjax)(strHref, params, 'text');\n\n  if (res.error && res.error.error_type == 'client') {\n    (0,_ajax_js__WEBPACK_IMPORTED_MODULE_0__.ajaxErr)(res.error.error_code, res.error.error_desc);\n    return;\n  }\n\n  if (res.error && res.error.error_type == 'server') {\n    res = res.response;\n  }\n\n  const elemActiveNavItem = document.querySelector('ul.mnav li a.nav-item-active');\n  if (elemActiveNavItem) elemActiveNavItem.classList.remove('nav-item-active');\n  const oParser = new DOMParser();\n  const oDoc = oParser.parseFromString(res, 'text/html');\n  const newContent = oDoc.getElementById('content');\n  let oldContent = document.getElementById('content');\n  document.title = oDoc.title;\n  (0,_utils_js__WEBPACK_IMPORTED_MODULE_1__.rewriteMetas)({\n    docSource: oDoc,\n    docDest: document,\n    metas: ['robots', 'og:title', 'og:description', 'og:type', 'og:image', 'og:url', 'og:site_name', 'og:locale', 'twitter:card', 'twitter:title', 'twitter:description', 'twitter:image'],\n    withCanonical: true\n  });\n  const elemDataExtras = oDoc.getElementById('data-extras');\n  if (elemDataExtras) dataExtras = JSON.parse(elemDataExtras.textContent);\n  hState.href = strHref;\n  hState.params = params;\n  hState.title = oDoc.title;\n  hState.url = strHref + (0,_utils_js__WEBPACK_IMPORTED_MODULE_1__.queryStringify)(params);\n  if (changeAddress) window.history.pushState(hState, hState.title, hState.url);\n  (0,_utils_js__WEBPACK_IMPORTED_MODULE_1__.fadeOut)(oldContent).then(() => {\n    (0,_utils_js__WEBPACK_IMPORTED_MODULE_1__.sleep)(110).then(() => {\n      oldContent.parentNode.replaceChild(newContent, oldContent);\n      window.onPageLoaded(dataExtras);\n    });\n  });\n  const scope = (0,_utils_js__WEBPACK_IMPORTED_MODULE_1__.getMeta)(oDoc, 'app:scope');\n  (0,_utils_js__WEBPACK_IMPORTED_MODULE_1__.setMeta)(document, 'app:scope', scope);\n  let itemActive = document.querySelector('ul.mnav li a[data-scope=\"' + scope + '\"]');\n  if (itemActive) itemActive.classList.add('nav-item-active');\n  if (callback) callback();\n} //  HistoryAPI: replace state on load\n\nfunction init() {\n  window.history.replaceState(hState, hState.title, hState.url);\n}\n; //  HistoryAPI: when back or forward\n\nfunction popstate(oEvent) {\n  loadPage(oEvent.state.href, oEvent.state.params);\n}\n; //  Exports\n\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  init,\n  popstate,\n  loadPage\n});\n\n//# sourceURL=webpack:///./scripts/spa.js?");

/***/ }),

/***/ "./scripts/utils.js":
/*!**************************!*\
  !*** ./scripts/utils.js ***!
  \**************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"sleep\": function() { return /* binding */ sleep; },\n/* harmony export */   \"animate\": function() { return /* binding */ animate; },\n/* harmony export */   \"makeLinear\": function() { return /* binding */ makeLinear; },\n/* harmony export */   \"makePow\": function() { return /* binding */ makePow; },\n/* harmony export */   \"makeCirc\": function() { return /* binding */ makeCirc; },\n/* harmony export */   \"makeEaseOut\": function() { return /* binding */ makeEaseOut; },\n/* harmony export */   \"makeEaseInOut\": function() { return /* binding */ makeEaseInOut; },\n/* harmony export */   \"makeLinearEaseInOut\": function() { return /* binding */ makeLinearEaseInOut; },\n/* harmony export */   \"makePowEaseOut\": function() { return /* binding */ makePowEaseOut; },\n/* harmony export */   \"makeCircEaseInOut\": function() { return /* binding */ makeCircEaseInOut; },\n/* harmony export */   \"drawOpacity\": function() { return /* binding */ drawOpacity; },\n/* harmony export */   \"fadeOut\": function() { return /* binding */ fadeOut; },\n/* harmony export */   \"drawHeight\": function() { return /* binding */ drawHeight; },\n/* harmony export */   \"movePX\": function() { return /* binding */ movePX; },\n/* harmony export */   \"queryStringify\": function() { return /* binding */ queryStringify; },\n/* harmony export */   \"queryParse\": function() { return /* binding */ queryParse; },\n/* harmony export */   \"getMeta\": function() { return /* binding */ getMeta; },\n/* harmony export */   \"setMeta\": function() { return /* binding */ setMeta; },\n/* harmony export */   \"rewriteMetas\": function() { return /* binding */ rewriteMetas; }\n/* harmony export */ });\n/* \r\n    UTILS SCRIPT\r\n\r\n    Author : IvanK Production\r\n*/\n////////////////////\n//  MAIN SECTION  //\n////////////////////\n//  sleep\nfunction sleep(ms) {\n  return new Promise(resolve => setTimeout(resolve, ms));\n} //////////////////////////\n//  ANIMATIONS SECTION  //\n//////////////////////////\n//  Animations: main function\n\nfunction animate(opts) {\n  let start = performance.now();\n  requestAnimationFrame(function animate(time) {\n    let timeFraction = (time - start) / opts.duration;\n    if (timeFraction > 1) timeFraction = 1;\n    if (timeFraction < 0) timeFraction = 0;\n    let progress = opts.timing(timeFraction);\n\n    if (opts.draw) {\n      opts.draw(progress);\n    }\n\n    if (opts.move) {\n      opts.move(progress);\n    }\n\n    if (timeFraction < 1) {\n      requestAnimationFrame(animate);\n    } else {\n      if (opts.callback) {\n        opts.callback();\n      }\n    }\n  });\n} //  Animations: linear\n\nfunction makeLinear(timeFraction) {\n  return timeFraction;\n} //  Animations: pow\n\nfunction makePow(timeFraction) {\n  return Math.pow(timeFraction, 5);\n} //  Animations: circ\n\nfunction makeCirc(timeFraction) {\n  return 1 - Math.sin(Math.acos(timeFraction));\n} //  Animations: EaseOut\n\nfunction makeEaseOut(timing) {\n  return function (timeFraction) {\n    return 1 - timing(1 - timeFraction);\n  };\n} //  Animations: EaseInOut\n\nfunction makeEaseInOut(timing) {\n  return function (timeFraction) {\n    if (timeFraction < 0.5) return timing(2 * timeFraction) / 2;else return (2 - timing(2 * (1 - timeFraction))) / 2;\n  };\n} //  Animations: complete timing functions\n\nvar makeLinearEaseInOut = makeEaseInOut(makeLinear);\nvar makePowEaseOut = makeEaseOut(makePow);\nvar makeCircEaseInOut = makeEaseInOut(makeCirc); //  Animations: opacity \n\nfunction drawOpacity(elem, value) {\n  elem.style.opacity = value;\n} //  Animations: async opacity (fadeout)\n\nasync function fadeOut(elem) {\n  elem.style.opacity = '0';\n}\n; //  Animations: height \n\nfunction drawHeight(elem, value) {\n  elem.style.height = value + 'px';\n} //  Animations: move by pixel parameter\n\nfunction movePX(elem, style, value) {\n  elem.style[style] = value + 'px';\n} /////////////////////\n//  QUERY SECTION  //\n/////////////////////\n//  Query object to string\n\nfunction queryStringify(obj) {\n  let params = [];\n  Object.keys(obj).forEach(key => {\n    if (obj[key] !== '') params.push(String(key + '=' + obj[key]).replace(/\\s/g, '_'));\n  });\n  return params.length ? '?' + params.join('&') : '';\n} //  Query string to object\n\nfunction queryParse(str) {\n  let result = {};\n  if (str == '') return result;\n  const obj = new URLSearchParams(str);\n\n  for (const [key, value] of obj.entries()) {\n    result[key] = value;\n  }\n\n  return result;\n}\n; ////////////////////\n//  META SECTION  //\n////////////////////\n//  Get canonical link\n\nfunction getCanonical(doc) {\n  let result = doc.querySelector('link[rel=\"canonical\"]');\n  if (!result) throw new Error(\"Canonical link not found!\");\n  return result.href;\n} //  Set canonical link\n\n\nfunction setCanonical(doc, value) {\n  let result = doc.querySelector('link[rel=\"canonical\"]');\n  if (!result) throw new Error(\"Canonical link not found!\");\n  result.href = value;\n} //  Get meta tag\n\n\nfunction getMeta(doc, metaName) {\n  let result = doc.querySelector('meta[name=\"' + metaName + '\"]');\n  if (!result) result = doc.querySelector('meta[property=\"' + metaName + '\"]');\n  if (!result) throw new Error(\"Meta not found!\");\n  return result.content;\n} //  Set meta tag\n\nfunction setMeta(doc, metaName, metaContent) {\n  let result = doc.querySelector('meta[name=\"' + metaName + '\"]');\n  if (!result) result = doc.querySelector('meta[property=\"' + metaName + '\"]');\n  if (!result) throw new Error(\"Meta not found!\");\n  result.content = metaContent;\n} //  Rewrite meta tags from ajax-loaded page to current page\n\nfunction rewriteMetas(opts) {\n  if (opts && opts.metas && opts.docSource && opts.docDest) {\n    opts.metas.forEach(elem => {\n      setMeta(opts.docDest, elem, getMeta(opts.docSource, elem));\n    });\n    if (opts.withCanonical) setCanonical(opts.docDest, getCanonical(opts.docSource));\n  } else {\n    throw new Error(\"Options not specified!\");\n  }\n} //////////////////////\n//  EXPORT SECTION  //\n//////////////////////\n\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  animate,\n  makeLinear,\n  makePow,\n  makeCirc,\n  makeEaseOut,\n  makeEaseInOut,\n  makeLinearEaseInOut,\n  makePowEaseOut,\n  makeCircEaseInOut,\n  drawOpacity,\n  drawHeight,\n  movePX,\n  fadeOut,\n  queryStringify,\n  queryParse,\n  rewriteMetas,\n  sleep\n});\n\n//# sourceURL=webpack:///./scripts/utils.js?");

/***/ }),

/***/ "./styles/bundle.css":
/*!***************************!*\
  !*** ./styles/bundle.css ***!
  \***************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n// extracted by mini-css-extract-plugin\n\n\n//# sourceURL=webpack:///./styles/bundle.css?");

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/define property getters */
/******/ 	!function() {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = function(exports, definition) {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	!function() {
/******/ 		__webpack_require__.o = function(obj, prop) { return Object.prototype.hasOwnProperty.call(obj, prop); }
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	!function() {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = function(exports) {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	}();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module can't be inlined because the eval devtool is used.
/******/ 	var __webpack_exports__ = __webpack_require__("./scripts/app.js");
/******/ 	
/******/ })()
;