/* 
    SLIDER COMPONENT SCRIPT

    Author : IvanK Production
*/

import { sleep } from "../utils.js";

export default class Slider {
	constructor(sliderContainerID, sliderConfig) {
		this.elemContainer = document.getElementById(sliderContainerID);
		if (!this.elemContainer) return null;

		this.slidesList = this.elemContainer.querySelectorAll('.slider-slide');
		if (!this.slidesList || this.slidesList.length == 0) return null;

		this.maxSlideHeight = Math.max.apply(null, [...this.slidesList].map(elem => {
			return elem.clientHeight;
		}));
		[...this.slidesList].filter(elem => { elem.style.height = this.maxSlideHeight + 'px'; });

		this.elemContainer.children[0].style.height = this.maxSlideHeight + 'px';

		this.activeIdx         = 0;
		this.events            = {};
		this.paginationBullets = [];

		this.speed = sliderConfig.speed ? sliderConfig.speed : 600;

		this.paginationElement = sliderConfig.pagination 
			&& sliderConfig.pagination.element 
			? this.elemContainer.querySelector(sliderConfig.pagination.element) : null;
		this.paginationClickable = sliderConfig.pagination 
			&& sliderConfig.pagination.clickable && this.paginationElement 
			? sliderConfig.pagination.clickable : false;

		if (this.paginationClickable) this.initPagination();

		this.autoplayEnabled = sliderConfig.autoplay ? true : false;
		this.autoplayDelay = sliderConfig.autoplay 
			&& sliderConfig.autoplay.delay 
			? sliderConfig.autoplay.delay : 0;

		this.fireCurr(0);

		return this;
	}

	get activeIndex() { return this.activeIdx - 1; }
	get slides()      { return this.slidesList;    }

	on(eventName, callback) { if (this.events) this.events[eventName] = callback; }

	fireNext() { if (this.activeIdx != undefined) this.fire(this.activeIdx - 1, this.activeIdx += 1); }
	firePrev() { if (this.activeIdx != undefined) this.fire(this.activeIdx - 1, this.activeIdx -= 1); }
	fireCurr(index) {
		if (!this.elemContainer) return false;

		if (this._timer) { window.clearInterval(this._timer); this._timer = null; }

		this.fire(this.activeIdx - 1, this.activeIdx = index + 1);
	}

	fire(prev, index) {
		if (!this.elemContainer) return false;

		const showSlide = () => {
			if (index > this.slidesList.length) this.activeIdx = 1;
			if (index < 1) this.activeIdx = this.slidesList.length;
	
			if (this.events['transitionStart']) this.events['transitionStart'](this);

			if (this.paginationBullets) this.paginationBullets[this.activeIdx - 1].classList.add('slider-pagination-bullet-active');
	
			this.slidesList[this.activeIdx - 1].style.zIndex  = '11';
			this.slidesList[this.activeIdx - 1].style.opacity = '1';
	
			sleep(this.speed + 10).then(() => {
				if (this.events['transitionEnd']) this.events['transitionEnd'](this);
	
				if (!this._timer && this.autoplayEnabled) this.startAutoplay();
			});
		};
		
		const activeSlide = this.slidesList[prev];
		if (activeSlide && activeSlide.style.opacity == '1') {
			activeSlide.style.opacity = '0';

			sleep(this.speed + 10).then(() => {
				activeSlide.style.zIndex  = '10';
				this.paginationBullets[prev].classList.remove('slider-pagination-bullet-active');

				showSlide();
			});
		} else { showSlide(); }
	}

	startAutoplay() {
		if (this.autoplayEnabled) { this._timer = window.setInterval(() => { this.fireNext(); }, this.autoplayDelay + this.speed * 2); }
		else return false;
	}

	initPagination() {
		if (!this.paginationElement) return false;

		const wrapper = this.paginationElement;

		this.slidesList.forEach((_, i) => {
			let bullet = document.createElement('div');

			bullet.classList.add('slider-pagination-bullet');
			if (this.paginationClickable) bullet.onclick = () => { this.fireCurr(i); };

			this.paginationBullets[i] = bullet;
			wrapper.append(bullet);
		});
	}

	destroy() {
		if (this._timer) { window.clearInterval(this._timer); this._timer = null; }

		this.events        = {};
		this.elemContainer = null;
		this.slidesList    = null;
	}
}
