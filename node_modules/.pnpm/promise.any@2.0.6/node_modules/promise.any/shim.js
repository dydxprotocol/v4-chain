'use strict';

var requirePromise = require('./requirePromise');

var getPolyfill = require('./polyfill');
var define = require('define-properties');

module.exports = function shimAny() {
	requirePromise();

	var polyfill = getPolyfill();
	define(Promise, { any: polyfill }, {
		any: function testAny() {
			return Promise.any !== polyfill;
		}
	});
	return polyfill;
};
