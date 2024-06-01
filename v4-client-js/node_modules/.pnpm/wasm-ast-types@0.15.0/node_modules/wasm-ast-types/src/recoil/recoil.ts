import * as t from '@babel/types';
import { camel, pascal } from 'case';
import {
  callExpression,
  getMessageProperties,
  getResponseType
} from '../utils';
import { QueryMsg } from '../types';
import { RenderContext } from '../context';

export const createRecoilSelector = (
  context: RenderContext,
  keyPrefix: string,
  QueryClient: string,
  methodName: string,
  responseType: string
) => {

  context.addUtil('selectorFamily');

  const selectorName = camel(`${methodName}Selector`);
  const getterKey = camel(`${keyPrefix}${pascal(methodName)}`);

  return t.exportNamedDeclaration(
    t.variableDeclaration(
      'const',
      [t.variableDeclarator(
        t.identifier(selectorName),
        callExpression(
          t.identifier('selectorFamily'),
          [
            t.objectExpression(
              [
                t.objectProperty(
                  t.identifier('key'),
                  t.stringLiteral(getterKey)
                ),
                t.objectProperty(
                  t.identifier('get'),
                  t.arrowFunctionExpression(
                    [
                      t.objectPattern(
                        [
                          t.objectProperty(
                            t.identifier('params'),
                            t.identifier('params'),
                            false,
                            true
                          ),
                          t.restElement(
                            t.identifier('queryClientParams')
                          )
                        ]
                      )
                    ],
                    t.arrowFunctionExpression(
                      [
                        t.objectPattern(
                          [
                            t.objectProperty(
                              t.identifier('get'),
                              t.identifier('get'),
                              false,
                              true
                            )
                          ]
                        )
                      ],
                      t.blockStatement(
                        [
                          t.variableDeclaration('const',
                            [
                              t.variableDeclarator(
                                t.identifier('client'),
                                t.callExpression(
                                  t.identifier('get'),
                                  [
                                    t.callExpression(
                                      t.identifier('queryClient'),
                                      [
                                        t.identifier('queryClientParams')
                                      ]
                                    )
                                  ]
                                )
                              )
                            ]),
                          t.returnStatement(
                            t.awaitExpression(
                              t.callExpression(
                                t.memberExpression(
                                  t.identifier('client'),
                                  t.identifier(methodName)
                                ),
                                [
                                  t.spreadElement(
                                    t.identifier('params')
                                  )
                                ]
                              )
                            )
                          )
                        ]
                      ),
                      true
                    )
                  )
                )
              ]
            )
          ],
          t.tsTypeParameterInstantiation(
            [
              t.tsTypeReference(
                t.identifier(responseType)
              ),
              t.tsIntersectionType(
                [
                  t.tsTypeReference(
                    t.identifier('QueryClientParams')
                  ),
                  t.tsTypeLiteral(
                    [
                      t.tsPropertySignature(
                        t.identifier('params'),
                        t.tsTypeAnnotation(
                          t.tsTypeReference(
                            t.identifier('Parameters'),
                            t.tsTypeParameterInstantiation(
                              [
                                t.tsIndexedAccessType(
                                  t.tsTypeReference(
                                    t.identifier(QueryClient)
                                  ),
                                  t.tsLiteralType(
                                    t.stringLiteral(methodName)
                                  )
                                )
                              ]
                            )
                          )
                        )
                      )
                    ]
                  )
                ]
              )
            ]
          )
        )
      )]
    )
  )

};

export const createRecoilSelectors = (
  context: RenderContext,
  keyPrefix: string,
  QueryClient: string,
  queryMsg: QueryMsg
) => {
  return getMessageProperties(queryMsg)
    .map(schema => {

      const underscoreName = Object.keys(schema.properties)[0];
      const methodName = camel(underscoreName);
      const responseType = getResponseType(context, underscoreName);

      return createRecoilSelector(
        context,
        keyPrefix,
        QueryClient,
        methodName,
        responseType
      );

    });
};

export const createRecoilQueryClientType = () => ({
  "type": "TSTypeAliasDeclaration",
  "id": {
    "type": "Identifier",
    "name": "QueryClientParams"
  },
  "typeAnnotation": {
    "type": "TSTypeLiteral",
    "members": [
      {
        "type": "TSPropertySignature",
        "key": {
          "type": "Identifier",
          "name": "contractAddress"
        },
        "computed": false,
        "typeAnnotation": {
          "type": "TSTypeAnnotation",
          "typeAnnotation": {
            "type": "TSStringKeyword"
          }
        }
      }
    ]
  }
});

export const createRecoilQueryClient = (
  context: RenderContext,
  keyPrefix: string,
  QueryClient: string
) => {

  context.addUtil('selectorFamily');

  const getterKey = camel(`${keyPrefix}${'QueryClient'}`);

  return t.exportNamedDeclaration(
    t.variableDeclaration(
      'const',
      [t.variableDeclarator(
        t.identifier('queryClient'),
        callExpression(
          t.identifier('selectorFamily'),
          [
            t.objectExpression(
              [
                t.objectProperty(
                  t.identifier('key'),
                  t.stringLiteral(getterKey)
                ),
                t.objectProperty(
                  t.identifier('get'),
                  t.arrowFunctionExpression(
                    [
                      t.objectPattern(
                        [
                          t.objectProperty(
                            t.identifier('contractAddress'),
                            t.identifier('contractAddress'),
                            false,
                            true
                          )
                        ]
                      )
                    ],
                    t.arrowFunctionExpression(
                      [
                        t.objectPattern(
                          [
                            t.objectProperty(
                              t.identifier('get'),
                              t.identifier('get'),
                              false,
                              true
                            )
                          ]
                        )
                      ],
                      t.blockStatement(
                        [
                          t.variableDeclaration('const',
                            [
                              t.variableDeclarator(
                                t.identifier('client'),
                                t.callExpression(
                                  t.identifier('get'),
                                  [
                                    t.identifier('cosmWasmClient')
                                  ]
                                )
                              )
                            ]),
                          t.returnStatement(
                            t.newExpression(
                              t.identifier(QueryClient),
                              [
                                t.identifier('client'),
                                t.identifier('contractAddress')
                              ]
                            )
                          )
                        ]
                      ),
                      false
                    )
                  )
                )
              ]
            )
          ],
          t.tsTypeParameterInstantiation(
            [
              t.tsTypeReference(
                t.identifier(QueryClient)
              ),
              t.tsTypeReference(
                t.identifier('QueryClientParams')
              )
            ]
          )
        )
      )]
    )
  )

};
