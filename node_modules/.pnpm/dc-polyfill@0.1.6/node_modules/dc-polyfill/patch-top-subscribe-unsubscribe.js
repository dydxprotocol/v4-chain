module.exports = function (unpatched) {
  const dc = { ...unpatched };

  dc.subscribe = (channel, cb) => {
    return dc.channel(channel).subscribe(cb);
  };

  dc.unsubscribe = (channel, cb) => {
    return dc.channel(channel).unsubscribe(cb);
  };

  return dc;
};
