'use strict';

var requirePromise = require('./requirePromise');

requirePromise();

var AggregateError = require('es-aggregate-error/polyfill')();
var PromiseResolve = require('es-abstract/2023/PromiseResolve');
var Type = require('es-abstract/2023/Type');
var callBind = require('call-bind');
var GetIntrinsic = require('get-intrinsic');
var iterate = require('iterate-value');
var map = require('array.prototype.map');

var all = callBind(GetIntrinsic('%Promise.all%'));
var reject = callBind(GetIntrinsic('%Promise.reject%'));
var $then = callBind(GetIntrinsic('%Promise.prototype.then%'));

module.exports = function any(iterable) {
	var C = this;
	if (Type(C) !== 'Object') {
		throw new TypeError('`this` value must be an object');
	}
	var thrower = function (value) {
		return reject(C, value);
	};
	try {
		return $then(
			all(C, map(iterate(iterable), function (item) {
				var itemPromise = PromiseResolve(C, item);
				return itemPromise.then(thrower, function identity(x) {
					return x;
				});
			})),
			function (errors) {
				throw new AggregateError(errors, 'Every promise rejected');
			},
			function (x) { return x; }
		);
	} catch (e) {
		return reject(C, e);
	}
};
