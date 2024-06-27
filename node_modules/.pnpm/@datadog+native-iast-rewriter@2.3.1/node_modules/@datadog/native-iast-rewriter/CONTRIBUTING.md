# Contributing to dd-native-iast-rewriter-js

Please reach out before starting work on any major code changes.
This will ensure we avoid duplicating work, or that your code can't be merged due to a rapidly changing
base. If you would like support for a module that is not listed, [contact support][1] to share a request.

[1]: https://docs.datadoghq.com/help

## Local setup

To set up the project locally, you should install `cargo` and `wasm-pack`:

```
$ curl https://sh.rustup.rs -sSf | sh

$ cargo install wasm-pack
```

and project dependencies:

```
$ npm install
```

### Build

Build the project with

```
$ npm run build

```

and then it will be possible to run the tests with

```
$ npm t
```
