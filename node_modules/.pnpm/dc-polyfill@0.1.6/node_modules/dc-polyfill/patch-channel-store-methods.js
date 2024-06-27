const { ReflectApply } = require('./primordials.js');

module.exports = function (unpatched) {
  const channels = new WeakSet();

  const dc_channel = unpatched.channel;

  const dc = { ...unpatched };

  dc.channel = function() {
    const ch = dc_channel.apply(this, arguments);

    if (channels.has(ch)) return ch;

    ch._stores = new Map();

    ch.bindStore = function(store, transform) {
      // const replacing = this._stores.has(store);
      // if (!replacing) channels.incRef(this.name);
      this._stores.set(store, transform);
    };

    ch.unbindStore = function(store) {
      if (!this._stores.has(store)) {
        return false;
      }

      this._stores.delete(store);

      // channels.decRef(this.name);
      // maybeMarkInactive(this);

      return true;
    };

    ch.runStores = function(data, fn, thisArg, ...args) {
      let run = () => {
        this.publish(data);
        return ReflectApply(fn, thisArg, args);
      };

      for (const entry of this._stores.entries()) {
        const store = entry[0];
        const transform = entry[1];
        run = wrapStoreRun(store, data, run, transform);
      }

      return run();
    };

    return ch;
  };

  return dc;
};

function wrapStoreRun(store, data, next, transform = defaultTransform) {
  return () => {
    let context;
    try {
      context = transform(data);
    } catch (err) {
      process.nextTick(() => {
        // triggerUncaughtException(err, false);
        throw err;
      });
      return next();
    }

    return store.run(context, next);
  };
}

function defaultTransform(data) {
  return data;
}
