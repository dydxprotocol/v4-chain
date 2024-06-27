const checks = require('./checks.js');

require('./primordials.js');

let dc = checks.hasDiagnosticsChannel()
  ? require('diagnostics_channel')
  : require('./reimplementation.js');

if (checks.hasGarbageCollectionBug()) {
  dc = require('./patch-garbage-collection-bug.js')(dc);
}

if (!checks.hasTopSubscribeUnsubscribe()) {
  dc = require('./patch-top-subscribe-unsubscribe.js')(dc);
}

if (!checks.hasChUnsubscribeReturn()) {
  dc = require('./patch-channel-unsubscribe-return.js')(dc);
}

if (!checks.hasChannelStoreMethods()) {
  dc = require('./patch-channel-store-methods.js')(dc);
}

if (!checks.hasTracingChannel()) {
  dc = require('./patch-tracing-channel.js')(dc);
}

if (checks.hasSyncUnsubscribeBug()) {
  dc = require('./patch-sync-unsubscribe-bug.js')(dc);
}

if (!checks.hasTracingChannelHasSubscribers()) {
  dc = require('./patch-tracing-channel-has-subscribers.js')(dc);
}

module.exports = dc;

