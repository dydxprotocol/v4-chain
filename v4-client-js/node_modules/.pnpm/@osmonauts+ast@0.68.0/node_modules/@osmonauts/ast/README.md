# @osmonauts/ast

Cosmos Typescript ASTs
## example usage

```js
import generate from '@babel/generator';
import * as t from '@babel/types';
import { recursiveNamespace } from '@osmonauts/ast';

const myModule = recursiveNamespace(
    [
        'osmosis',
        'gamm',
        'v1beta',
        'pools'
    ].reverse(), [astBody]);

console.log(generate(t.program(myModule)).code)
```

produces:

```js
export namespace osmosis {
  export namespace gamm {
    export namespace v1beta {
      export namespace pools {

          // astBody here

      }
    }
  }
}
```