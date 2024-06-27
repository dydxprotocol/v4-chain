// There's a bug where a newly created channel is immediately garbage collected
// @see https://github.com/nodejs/node/pull/47520
const PHONY_SUBSCRIBE = function AVOID_GARBAGE_COLLECTION() {};

const {
  ObjectDefineProperty,
  ObjectGetOwnPropertyDescriptor
} = require('./primordials.js');

module.exports = function(unpatched) {
  const dc_channel = unpatched.channel;
  const channels = new WeakSet();

  const dc = { ...unpatched };

  dc.channel = function() {
    const ch = dc_channel.apply(this, arguments);

    if (channels.has(ch)) return ch;

    dc_channel(arguments[0]).subscribe(PHONY_SUBSCRIBE);

    channels.add(ch);

    if (!ObjectGetOwnPropertyDescriptor(ch, 'hasSubscribers')) {
      ObjectDefineProperty(ch, 'hasSubscribers', {
        get: function() {
          const subscribers = ch._subscribers;
          if (subscribers.length > 1) return true;
          if (subscribers.length < 1) return false;
          if (subscribers[0] === PHONY_SUBSCRIBE) return false;
          return true;
        },
      });
    }

    return ch;
  };

  return dc;
};
