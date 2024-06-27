const {
  ReflectApply,
  PromiseReject,
  PromiseResolve,
  PromisePrototypeThen,
  ArrayPrototypeSplice,
  ArrayPrototypeAt,
} = require('./primordials.js');

const { ERR_INVALID_ARG_TYPE } = require('./errors.js');

const traceEvents = [
  'start',
  'end',
  'asyncStart',
  'asyncEnd',
  'error',
];

module.exports = function (unpatched) {
  const { channel } = unpatched;

  const dc = { ...unpatched };

  class TracingChannel {
    constructor(nameOrChannels) {
      if (typeof nameOrChannels === 'string') {
        this.start = channel(`tracing:${nameOrChannels}:start`);
        this.end = channel(`tracing:${nameOrChannels}:end`);
        this.asyncStart = channel(`tracing:${nameOrChannels}:asyncStart`);
        this.asyncEnd = channel(`tracing:${nameOrChannels}:asyncEnd`);
        this.error = channel(`tracing:${nameOrChannels}:error`);
      } else if (typeof nameOrChannels === 'object') {
        const { start, end, asyncStart, asyncEnd, error } = nameOrChannels;

        // assertChannel(start, 'nameOrChannels.start');
        // assertChannel(end, 'nameOrChannels.end');
        // assertChannel(asyncStart, 'nameOrChannels.asyncStart');
        // assertChannel(asyncEnd, 'nameOrChannels.asyncEnd');
        // assertChannel(error, 'nameOrChannels.error');

        this.start = start;
        this.end = end;
        this.asyncStart = asyncStart;
        this.asyncEnd = asyncEnd;
        this.error = error;
      } else {
        throw new ERR_INVALID_ARG_TYPE('nameOrChannels',
                                       ['string', 'object', 'Channel'],
                                       nameOrChannels);
      }
    }

    subscribe(handlers) {
      for (const name of traceEvents) {
        if (!handlers[name]) continue;

        if (this[name]) this[name].subscribe(handlers[name]);
      }
    }

    unsubscribe(handlers) {
      let done = true;

      for (const name of traceEvents) {
        if (!handlers[name]) continue;

        if (!(this[name] && this[name].unsubscribe(handlers[name]))) {
          done = false;
        }
      }

      return done;
    }

    traceSync(fn, context = {}, thisArg, ...args) {
      const { start, end, error } = this;

      return start.runStores(context, () => {
        try {
          const result = ReflectApply(fn, thisArg, args);
          context.result = result;
          return result;
        } catch (err) {
          context.error = err;
          error.publish(context);
          throw err;
        } finally {
          end.publish(context);
        }
      });
    }

    tracePromise(fn, context = {}, thisArg, ...args) {
      const { start, end, asyncStart, asyncEnd, error } = this;

      function reject(err) {
        context.error = err;
        error.publish(context);
        asyncStart.publish(context);

        asyncEnd.publish(context);
        return PromiseReject(err);
      }

      function resolve(result) {
        context.result = result;
        asyncStart.publish(context);

        asyncEnd.publish(context);
        return result;
      }

      return start.runStores(context, () => {
        try {
          let promise = ReflectApply(fn, thisArg, args);
          // Convert thenables to native promises
          if (!(promise instanceof Promise)) {
            promise = PromiseResolve(promise);
          }
          return PromisePrototypeThen(promise, resolve, reject);
        } catch (err) {
          context.error = err;
          error.publish(context);
          throw err;
        } finally {
          end.publish(context);
        }
      });
    }

    traceCallback(fn, position = -1, context = {}, thisArg, ...args) {
      const { start, end, asyncStart, asyncEnd, error } = this;

      function wrappedCallback(err, res) {
        if (err) {
          context.error = err;
          error.publish(context);
        } else {
          context.result = res;
        }

        // Using runStores here enables manual context failure recovery
        asyncStart.runStores(context, () => {
          try {
            if (callback) {
              return ReflectApply(callback, this, arguments);
            }
          } finally {
            asyncEnd.publish(context);
          }
        });
      }

      const callback = ArrayPrototypeAt(args, position);
      if (typeof callback !== 'function') {
        throw new ERR_INVALID_ARG_TYPE('callback', ['function'], callback);
      }
      ArrayPrototypeSplice(args, position, 1, wrappedCallback);

      return start.runStores(context, () => {
        try {
          return ReflectApply(fn, thisArg, args);
        } catch (err) {
          context.error = err;
          error.publish(context);
          throw err;
        } finally {
          end.publish(context);
        }
      });
    }
  }

  function tracingChannel(nameOrChannels) {
    return new TracingChannel(nameOrChannels);
  }

  dc.tracingChannel = tracingChannel;

  return dc;
};
