const { ObjectGetPrototypeOf } = require('./primordials.js');

// The ch.unsubscribe() method doesn't return a value
// Recent versions return if an unsubscribe succeeded
// @see https://github.com/nodejs/node/pull/40433
module.exports = function (unpatched) {
  const channels = new WeakSet();

  const dc_channel = unpatched.channel;

  const dc = { ...unpatched };

  dc.channel = function() {
    const ch = dc_channel.apply(this, arguments);

    if (channels.has(ch)) return ch;

    ch.unsubscribe = function () {
      delete ch.unsubscribe;

      const oldSubscriberCount = this._subscribers.length;

      ObjectGetPrototypeOf(ch).unsubscribe.apply(this, arguments);

      return this._subscribers.length < oldSubscriberCount;
    };

    channels.add(ch);

    return ch;
  };

  return dc;
};
