var test = require('tape');
var isError = require('.');

test('isError(value)', function (t) {
  t.ok(new Error('Error'));
  t.ok(new RangeError('Error'));
  t.notOk(isError({}));
  t.end();
});
