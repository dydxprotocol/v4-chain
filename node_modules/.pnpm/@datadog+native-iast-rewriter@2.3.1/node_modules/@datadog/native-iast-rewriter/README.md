# dd-native-iast-rewriter-js

Nodejs native AST rewriter heavily based on [Speedy Web Compiler o SWC compiler](https://github.com/swc-project/swc) used to instrument Javascript source files.

## Workflow

1. Parse Javascript code to obtain the AST
2. Replace certain AST expressions -> Currently it is focused in certain operators like `+` and `+=`, template literals and some String methods
3. Generate the new Javascript code as from the modified AST -> In addition to the Javascript code, the corresponding source map is returned chaining it with the original source map if necessary.

## Usage

```javascript
const Rewriter = require('@datadog/native-iast-rewriter')

const rewriter = new Rewriter(rewriterConfig)
const result = rewriter.rewrite(code, filename)
```

## Configuration options

```javascript

RewriterConfig {
  // enable/disable sourceMap chaining - false by default
  chainSourceMap?: boolean

  // enable/disable comments printing - false by default
  comments?: boolean

  // establishes the prefix for the injected local variables - 6 random characters by default
  localVarPrefix?: string

  // sets the list of methods or operators to be rewritten
  csiMethods?: Array<CsiMethod>

  // extracts hardcoded string literals - true by default
  literals?: boolean
}

CsiMethod {
  // name of the String method to rewrite
  src: string

  // optional name of the replacement method. If not specified a convention shall be used
  dst?: string

  // indicates if it is an operator like +
  operator?: boolean
}
```

## Example

```javascript
const Rewriter = require('@datadog/native-iast-rewriter')

const rewriterConfig = {
  csiMethods: [{ src: 'substring' }],
  localVarPrefix: 'test',
}

const rewriter = new Rewriter(rewriterConfig)

const code = `function sub(a) {
  return a.substring(1)
}`
const result = rewriter.rewrite(code, filename)

console.log(result.content)
/*
function sub(a) {
  let __datadog_test_0, __datadog_test_1;
  return (__datadog_test_0 = a, __datadog_test_1 = __datadog_test_0.substring, _ddiast.stringSubstring(__datadog_test_1.call(__datadog_test_0, 1), __datadog_test_1, __datadog_test_0, 1));
}
*/
```

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

It will compile WASM binaries by default
and then it will be possible to run the tests with

```
$ npm t
```
