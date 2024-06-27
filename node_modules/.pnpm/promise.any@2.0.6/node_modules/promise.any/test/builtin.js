'use strict';

var defineProperties = require('define-properties');
var isEnumerable = Object.prototype.propertyIsEnumerable;
var functionsHaveNames = require('functions-have-names')();

var runTests = require('./tests');

module.exports = function (t) {
	t.equal(Promise.any.length, 1, 'Promise.any has a length of 1');
	t.test('Function name', { skip: !functionsHaveNames }, function (st) {
		st.equal(Promise.any.name, 'any', 'Promise.any has name "any"');
		st.end();
	});

	t.test('enumerability', { skip: !defineProperties.supportsDescriptors }, function (et) {
		et.equal(false, isEnumerable.call(Promise, 'any'), 'Promise.any is not enumerable');
		et.end();
	});

	var supportsStrictMode = (function () { return typeof this === 'undefined'; }());

	t.test('bad object value', { skip: !supportsStrictMode }, function (st) {
		st['throws'](function () { return Promise.any.call(undefined); }, TypeError, 'undefined is not an object');
		st['throws'](function () { return Promise.any.call(null); }, TypeError, 'null is not an object');
		st.end();
	});

	runTests(function any(iterable) { return Promise.any.call(typeof this === 'undefined' ? Promise : this, iterable); }, t);
};
