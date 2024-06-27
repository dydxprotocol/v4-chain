'use strict';

var AggregateError = require('es-aggregate-error/polyfill')();

var assertArray = function (t, value, length, assertType) {
	t.ok(Array.isArray(value), 'value is an array');
	t.equal(value.length, length, 'length is ' + length);
	if (typeof assertType === 'function') {
		for (var i = 0; i < value.length; i += 1) {
			assertType(value[i]);
		}
	}
};

module.exports = function (any, t) {
	if (typeof Promise !== 'function') {
		return t.skip('No global Promise detected');
	}

	var a = {};
	var b = {};
	var c = {};

	t.test('empty iterable', function (st) {
		st.plan(2);
		any([]).then(
			function () { st.fail(); },
			function (error) {
				st.equal(error instanceof AggregateError, true, 'is an AggregateError');
				st.deepEqual(error.errors, []);
			}
		);
	});

	t.test('no promise values', function (st) {
		st.plan(1);
		any([a, b, c]).then(function (result) {
			st.deepEqual(result, a);
		});
	});

	t.test('all fulfilled', function (st) {
		st.plan(1);
		any([
			Promise.resolve(a),
			Promise.resolve(b),
			Promise.resolve(c)
		]).then(function (result) {
			st.deepEqual(result, a);
		});
	});

	t.test('all rejected', function (st) {
		st.plan(2);
		any([
			Promise.reject(a),
			Promise.reject(b),
			Promise.reject(c)
		]).then(
			function () { st.fail(); },
			function (error) {
				st.equal(error instanceof AggregateError, true, 'is an AggregateError');
				st.deepEqual(error.errors, [a, b, c]);
			}
		);
	});

	t.test('mixed', function (st) {
		st.plan(2);
		any([
			a,
			Promise.resolve(b),
			Promise.reject(c)
		]).then(function (result) {
			st.deepEqual(result, a);
		});

		any([
			Promise.reject(a),
			Promise.resolve(b),
			Promise.reject(c)
		]).then(function (result) {
			st.deepEqual(result, b);
		});
	});

	t.test('poisoned .then', function (st) {
		st.plan(1);
		var poison = new EvalError();
		var promise = new Promise(function () {});
		promise.then = function poisionedThen() { throw poison; };
		any([promise]).then(function () {
			st.fail('should not reach here');
		}, function (error) {
			st.equal(error, poison, 'error is whatever the poisoned then throws');
		});
	});

	var Subclass = (function () {
		try {
			// eslint-disable-next-line no-new-func
			return Function('class Subclass extends Promise { constructor(...args) { super(...args); this.thenArgs = []; } then(...args) { Subclass.thenArgs.push(args); this.thenArgs.push(args); return super.then(...args); } } Subclass.thenArgs = []; return Subclass;')();
		} catch (e) { /**/ }

		return false;
	}());
	t.test('inheritance', { skip: !Subclass }, function (st) {
		st.test('preserves correct subclass', function (s2t) {
			var promise = any.call(Subclass, [1]);
			s2t.ok(promise instanceof Subclass, 'promise is instanceof Subclass');
			s2t.equal(promise.constructor, Subclass, 'promise.constructor is Subclass');

			s2t.end();
		});

		st.test('invokes the subclass’ then', function (s2t) {
			Subclass.thenArgs.length = 0;

			var original = Subclass.resolve();
			assertArray(s2t, Subclass.thenArgs, 0);
			assertArray(s2t, original.thenArgs, 0);

			any.call(Subclass, [original]);

			assertArray(s2t, original.thenArgs, 1);
			assertArray(s2t, original.thenArgs[0], 2);

			s2t.test('proper subclass then invocation count', { todo: true }, function (s3t) {
				// native implementations report 1, this implementation reports 2
				assertArray(s3t, Subclass.thenArgs, 1);

				s3t.end();
			});
			s2t.ok(Array.isArray(Subclass.thenArgs), 'value is an array');
			s2t.match(String(Subclass.thenArgs.length), /^[12]$/, 'length is 1 or 2');

			assertArray(s2t, Subclass.thenArgs[0], 2);

			s2t.end();
		});
	});

	return t.comment('tests completed');
};
