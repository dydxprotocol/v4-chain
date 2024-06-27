const {
  ObjectDefineProperty,
  ObjectGetPrototypeOf,
} = require('./primordials.js');

module.exports = function (unpatched) {
  const dc = { ...unpatched };

  {
    const fauxTrCh = dc.tracingChannel('dc-polyfill-faux');

    const protoTrCh = ObjectGetPrototypeOf(fauxTrCh);

    ObjectDefineProperty(protoTrCh, 'hasSubscribers', {
      get: function () {
        return this.start.hasSubscribers
          || this.end.hasSubscribers
          || this.asyncStart.hasSubscribers
          || this.asyncEnd.hasSubscribers
          || this.error.hasSubscribers;
      },
      configurable: true
    });
  }

  return dc;
};
