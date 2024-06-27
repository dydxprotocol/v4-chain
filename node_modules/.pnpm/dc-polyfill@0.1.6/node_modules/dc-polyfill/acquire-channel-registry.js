/**
 * This file allows for a global shared channels registry
 * across multiple installed versions of this package.
 *
 * Without such a registry, if someone were to have versions
 * 1.0.0 and 1.0.1 both installed in their node_modules directory,
 * then there would be two unrelated channel collections. This
 * would thus not allow the two collections to communicate
 * across each others channels.
 *
 * Of course, when using the built-in diagnostics_channel, this
 * isn't a problem, as the different versions of this package all
 * share the same internal channels registry.
 *
 * We're attaching this to process instead of global as there are
 * test suites which make sure globals don't get polluted between
 * test runs.
 */

const { ObjectDefineProperty, SymbolFor } = require('./primordials.js');

const REGISTRY_SYMBOL = SymbolFor('dc-polyfill-v1');

if (!process[REGISTRY_SYMBOL]) {
  ObjectDefineProperty(process, REGISTRY_SYMBOL, {
    configurable: false,
    enumerable: false,
    writable: false,
    value: {},
  });
}

module.exports = process[REGISTRY_SYMBOL];
