// https://github.com/nodejs/node/pull/48933

module.exports = function(unpatched) {
  const channels = new WeakSet();

  const dc_channel = unpatched.channel;

  const dc = { ...unpatched };

  dc.channel = function() {
    const ch = dc_channel.apply(this, arguments);

    if (channels.has(ch)) return ch;

    const publish = ch.publish;

    ch.publish = function() {
      if (!ch._subscribers) {
        ch._subscribers = [];
      }

      return publish.apply(ch, arguments);
    };

    return ch;
  };

  return dc;
};
