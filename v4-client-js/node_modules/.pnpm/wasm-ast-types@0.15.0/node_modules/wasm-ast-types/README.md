# wasm-ast-types

## working with ASTs

### 1 edit the fixture

edit `./scripts/fixture.ts`, for example:

```js
// ./scripts/fixture.ts
export interface InstantiateMsg {
    admin?: string | null;
    members: Member[];
}
```

### 2 run AST generator

```
yarn test:ast
```

### 3 look at the JSON produced

```
code ./scripts/test-output.json
```

We use the npm module `ast-stringify` to strip out unneccesary props, and generate a JSON for reference.

You will see a `File` and `Program`... only concern yourself with the `body[]`:

```json
{
  "type": "File",
  "errors": [],
  "program": {
    "type": "Program",
    "sourceType": "module",
    "interpreter": null,
    "body": [
      {
        "type": "ExportNamedDeclaration",
        "exportKind": "type",
        "specifiers": [],
        "source": null,
        "declaration": {
          "type": "TSInterfaceDeclaration",
          "id": {
            "type": "Identifier",
            "name": "InstantiateMsg"
          },
          "body": {
            "type": "TSInterfaceBody",
            "body": [
              {
                "type": "TSPropertySignature",
                "key": {
                  "type": "Identifier",
                  "name": "admin"
                },
                "computed": false,
                "optional": true,
                "typeAnnotation": {
                  "type": "TSTypeAnnotation",
                  "typeAnnotation": {
                    "type": "TSUnionType",
                    "types": [
                      {
                        "type": "TSStringKeyword"
                      },
                      {
                        "type": "TSNullKeyword"
                      }
                    ]
                  }
                }
              },
              {
                "type": "TSPropertySignature",
                "key": {
                  "type": "Identifier",
                  "name": "members"
                },
                "computed": false,
                "typeAnnotation": {
                  "type": "TSTypeAnnotation",
                  "typeAnnotation": {
                    "type": "TSArrayType",
                    "elementType": {
                      "type": "TSTypeReference",
                      "typeName": {
                        "type": "Identifier",
                        "name": "Member"
                      }
                    }
                  }
                }
              }
            ]
          }
        }
      }
    ],
    "directives": []
  },
  "comments": []
}
```

### 4 code with `@babel/types` using the JSON as a reference

NOTE: 4 continued ideally you should be writing a test with your generator!

```js
import * as t from '@babel/types';

export const createNewGenerator = () => {
    return t.exportNamedDeclaration(
        t.tsInterfaceDeclaration(
            t.identifier('InstantiateMsg'),
            null,
            [],
            t.tsInterfaceBody([
                // ... more code ...
            ])
        )
    );
};
```
