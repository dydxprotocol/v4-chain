# `dc-polyfill`: Diagnostics Channel Polyfill

This package provides a polyfill (or ponyfill) for the `diagnostics_channel` core Node.js module (including `TracingChannel`) for use with older versions of Node.js. It aims to remain simple, with zero dependencies, and only takes up a few kilobytes of space.

> If your module or application uses `diagnostics_channel` and needs to run on multiple versions of Node.js then it is recommended that you use `require('dc-polyfill')` instead of `require('diagnostics_channel')`.

**dc-polyfill** backports features and bugfixes that are added to Node.js core. If a feature hasn't been backported then please open a Pull Request or create an issue.

Since this package recreates a Node.js API, read the [Node.js `diagnostics_channel` documentation](https://nodejs.org/dist/latest-v20.x/docs/api/diagnostics_channel.html) to understand what it does.

|                                  | Version |
|----------------------------------|---------|
| Oldest Supported Node.js Version | 12.17.0 |
| Target Node.js DC API Version    | 20.6.0  |

> Note that `dc-polyfill` currently has the `TracingChannel#hasSubscribers` getter backported from Node.js v22 however it doesn't yet support the tracing channel early exit feature. Once that's been added we'll delete this clause and update the above table.

Whenever the currently running version of Node.js ships with `diagnostics_channel` (i.e. v16+, v15.14+, v14.17+), **dc-polyfill** will make sure to use the global registry of channels provided by the core module. For older versions of Node.js **dc-polyfill** instead uses its own global collection of channels. This global collection remains in the same location and is shared across all instances of **dc-polyfill**. This avoids the issue wherein multiple versions of an npm library installed in a module dependency hierarchy would otherwise provide different singleton instances.

Ideally, this package will forever remain backwards compatible, there will never be a v2.x release, and there will never be an additional global channel collection.


## Usage

Install the module in your project:

```sh
npm install dc-polyfill
```

Replace any existing `require('diagnostics_channel')` calls:

```javascript
const diagnostics_channel = require('dc-polyfill');
```


## Contributing

When a Pull Request is created the test suite runs against dozens of versions of Node.js. Notably, versions right before a change and versions right after a change, the first version of a release line, and the last version of a release line. To test locally it's recommended to use a node version management tool, such as `nvm`, to test changes with.


## License / Copyright

See [LICENSE.txt](LICENSE.txt) for full details.

MIT License - Copyright (c) 2023 Datadog, Inc.
