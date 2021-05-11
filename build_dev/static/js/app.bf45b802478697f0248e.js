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

/***/ "./scripts/app.js":
/*!************************!*\
  !*** ./scripts/app.js ***!
  \************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _styles_bundle_css__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../styles/bundle.css */ \"./styles/bundle.css\");\n/* \r\n    MAIN SCRIPT\r\n\r\n    Author : IvanK Production\r\n*/\n //  Elements\n\nlet elemSubnavWrappers = document.querySelectorAll('.subnav-container');\n/*const elemFooterBlockFirst  = document.querySelector(\".footer-block.first\");\r\nconst elemFooterBlockSecond = document.querySelector(\".footer-block.second\");\r\nconst elemFooterBlockThird  = document.querySelector(\".footer-block.third\");\r\nconst elemFooterBlockHelper = document.querySelector(\".footer-block.helper\");\r\nlet   elemFooterContainer   = elemFooterBlockThird.parentNode;*/\n//  Animations: Animate function\n\nfunction animate(opts) {\n  var start = performance.now();\n  requestAnimationFrame(function animate(time) {\n    var timeFraction = (time - start) / opts.duration;\n    if (timeFraction > 1) timeFraction = 1;\n    var progress = opts.timing(timeFraction);\n    if (opts.draw) opts.draw(progress);\n    if (opts.move) opts.move(progress);\n\n    if (timeFraction < 1) {\n      requestAnimationFrame(animate);\n    } else {\n      opts.callback();\n    }\n  });\n} //  Animations: linear\n\n\nfunction makeLinear(timeFraction) {\n  return timeFraction;\n} //  Animations: EaseInOut\n\n\nfunction makeEaseInOut(timing) {\n  return function (timeFraction) {\n    if (timeFraction < 0.5) return timing(2 * timeFraction) / 2;else return (2 - timing(2 * (1 - timeFraction))) / 2;\n  };\n} //  Animations: Complete timing function\n\n\nvar makeLinearEaseInOut = makeEaseInOut(makeLinear); //  Animations: Opacity \n\nfunction drawOpacity(elem, value, grd = 100) {\n  elem.style['-webkit-opacity'] = value / 100 * grd;\n  elem.style['-khtml-opacity'] = value / 100 * grd;\n  elem.style['-moz-opacity'] = value / 100 * grd;\n  elem.style['opacity'] = value / 100 * grd;\n} //  swapping footer blocks\n\n/*function fixFooterBlocks() {\r\n\tif (window.matchMedia(\"(min-width: 1200px)\").matches) {\r\n\t\telemFooterContainer.append(elemFooterBlockFirst, elemFooterBlockSecond, elemFooterBlockThird, elemFooterBlockHelper);\r\n\t} else if (window.matchMedia(\"(min-width: 992px)\").matches) {\r\n\t\telemFooterContainer.append(elemFooterBlockFirst, elemFooterBlockSecond, elemFooterBlockThird, elemFooterBlockHelper);\r\n\t} else if (window.matchMedia(\"(min-width: 768px)\").matches) {\r\n\t\telemFooterContainer.append(elemFooterBlockSecond, elemFooterBlockThird, elemFooterBlockFirst, elemFooterBlockHelper);\r\n\t} else if (window.matchMedia(\"(min-width: 10px)\").matches) {\r\n\t\telemFooterContainer.append(elemFooterBlockSecond, elemFooterBlockThird, elemFooterBlockFirst, elemFooterBlockHelper);\r\n\t}\r\n}*/\n//\n//  MAIN PAGE EVENTS:\n//\n//  DOMContentLoaded\n\n\ndocument.addEventListener(\"DOMContentLoaded\", function () {\n  //  Submenus\n  document.querySelectorAll('a.subnav').forEach(function (elem, i) {\n    elem.addEventListener('mouseover', function () {\n      elemSubnavWrappers[i].classList.add('showed');\n      drawOpacity(elemSubnavWrappers[i], 1);\n    });\n    elem.addEventListener('mouseleave', function (event) {\n      if (event.relatedTarget != elemSubnavWrappers[i]) {\n        elem.classList.remove('hovered');\n        animate({\n          \"duration\": 100,\n          \"timing\": makeLinearEaseInOut,\n          \"draw\": function (perc) {\n            drawOpacity(elemSubnavWrappers[i], 1 - perc);\n          },\n          \"callback\": function () {\n            elemSubnavWrappers[i].classList.remove('showed');\n          }\n        });\n      }\n\n      elemSubnavWrappers[i].addEventListener('mouseover', function () {\n        this.classList.add('showed');\n        drawOpacity(this, 1);\n        elem.classList.add('hovered');\n      });\n      elemSubnavWrappers[i].addEventListener('mouseleave', function (_event) {\n        if (_event.relatedTarget != event.target) {\n          elem.classList.remove('hovered');\n          animate({\n            \"duration\": 100,\n            \"elemw\": this,\n            \"timing\": makeLinearEaseInOut,\n            \"draw\": function (perc) {\n              drawOpacity(this.elemw, 1 - perc);\n            },\n            \"callback\": function () {\n              this.elemw.classList.remove('showed');\n            }\n          });\n        }\n      });\n    });\n  }); //fixFooterBlocks();\n}); //  onResize\n\nwindow.onresize = function () {//fixFooterBlocks();\n}; //  onLoad\n\n\nwindow.onload = function () {\n  document.body.classList.remove('preload');\n  /* FOR TEST\r\n  setTimeout(function() {*/\n\n  clearTimeout(window.tLoader);\n  document.getElementById('loader').style.display = 'none';\n  document.getElementById('master-container').style.opacity = '1';\n  /*}, 2000);*/\n}; // Buttons onClick\n\n\n[...document.getElementsByTagName('button')].forEach(elem => {\n  elem.addEventListener('click', function (e) {\n    window.location.href = e.target.dataset.href;\n  });\n});\n\n//# sourceURL=webpack:///./scripts/app.js?");

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