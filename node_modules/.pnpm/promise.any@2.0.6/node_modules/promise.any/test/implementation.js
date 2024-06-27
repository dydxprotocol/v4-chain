'use strict';

var test = require('tape');
var callBind = require('call-bind');

var any = require('../implementation');
var runTests = require('./tests');

var bound = callBind(any);

// eslint-disable-next-line no-shadow
var rebindable = function any(iterable) {
	return bound(typeof this === 'undefined' ? Promise : this, iterable);
};

test('as a function', function (t) {
	t.test('bad Promise/this value', function (st) {
		// eslint-disable-next-line no-useless-call
		st['throws'](function () { any.call(undefined, []); }, TypeError, 'undefined is not an object');

		// eslint-disable-next-line no-useless-call
		st['throws'](function () { any.call(null, []); }, TypeError, 'null is not an object');
		st.end();
	});

	runTests(rebindable, t);

	t.end();
});
