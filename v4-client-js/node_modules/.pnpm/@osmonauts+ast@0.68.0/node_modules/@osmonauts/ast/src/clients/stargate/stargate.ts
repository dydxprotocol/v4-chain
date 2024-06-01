import * as t from '@babel/types';
import { GenericParseContext } from '../../encoding';
import { arrowFunctionExpression, identifier, memberExpressionOrIdentifier, objectPattern } from '../../utils';

interface CreateStargateClient {
  name: string;
  options: string;
  context: GenericParseContext;
}

interface CreateStargateClientProtoRegistry {
  registries: string[];
  protoTypeRegistry: string;
  context: GenericParseContext;
}

interface CreateStargateClientOptions {
  name: string;
  aminoConverters: string;
  protoTypeRegistry: string;
  context: GenericParseContext;
}
interface CreateStargateClientAminoConverters {
  aminos: string[];
  aminoConverters: string;
  context: GenericParseContext;
}


export const createStargateClient = ({ name, options, context }: CreateStargateClient) => {

  const includeDefaults = context.pluginValue('stargateClients.includeCosmosDefaultTypes');
  let optsFuncArgs = [];

  if (includeDefaults) {
    context.addUtil('defaultRegistryTypes');
    optsFuncArgs = [
      t.objectExpression([
        t.objectProperty(
          t.identifier('defaultTypes'),
          t.identifier('defaultTypes'),
          false,
          true
        )
      ])
    ];
  }

  context.addUtil('GeneratedType');
  context.addUtil('OfflineSigner');
  context.addUtil('Registry');
  context.addUtil('AminoTypes');
  context.addUtil('SigningStargateClient');
  context.addUtil('HttpEndpoint');

  const prop = t.tsPropertySignature(
    t.identifier('defaultTypes'),
    t.tsTypeAnnotation(
      t.tsTypeReference(t.identifier('ReadonlyArray'), t.tsTypeParameterInstantiation(
        [
          t.tsTupleType([
            t.tsStringKeyword(),
            t.tsTypeReference(
              t.identifier('GeneratedType')
            )]
          )
        ]
      ))
    )
  );
  prop.optional = true;

  return t.exportNamedDeclaration(
    t.variableDeclaration(
      'const',
      [
        t.variableDeclarator(
          t.identifier(name),
          t.arrowFunctionExpression(
            [
              objectPattern(
                [
                  t.objectProperty(
                    t.identifier('rpcEndpoint'),
                    t.identifier('rpcEndpoint'),
                    false,
                    true
                  ),
                  t.objectProperty(
                    t.identifier('signer'),
                    t.identifier('signer'),
                    false,
                    true
                  ),
                  includeDefaults && t.objectProperty(
                    t.identifier('defaultTypes'),
                    t.assignmentPattern(
                      t.identifier('defaultTypes'),
                      t.identifier('defaultRegistryTypes')
                    ),
                    false,
                    true
                  )
                ].filter(Boolean),
                t.tsTypeAnnotation(
                  t.tsTypeLiteral(
                    [
                      t.tsPropertySignature(
                        t.identifier('rpcEndpoint'),
                        t.tsTypeAnnotation(
                          t.tsUnionType([
                            t.tsStringKeyword(),
                            t.tsTypeReference(
                              t.identifier('HttpEndpoint')
                            )
                          ])
                        )
                      ),
                      t.tsPropertySignature(
                        t.identifier('signer'),
                        t.tsTypeAnnotation(t.tsTypeReference(
                          t.identifier('OfflineSigner')
                        ))
                      ),
                      includeDefaults && prop
                    ].filter(Boolean)
                  )
                )
              )
            ],
            t.blockStatement(
              [
                // props
                t.variableDeclaration(
                  'const',
                  [

                    t.variableDeclarator(
                      t.objectPattern(
                        [
                          t.objectProperty(
                            t.identifier('registry'),
                            t.identifier('registry'),
                            false,
                            true
                          ),
                          t.objectProperty(
                            t.identifier('aminoTypes'),
                            t.identifier('aminoTypes'),
                            false,
                            true
                          ),
                        ]
                      ),
                      t.callExpression(
                        t.identifier(options),
                        optsFuncArgs
                      )
                    )
                  ],
                ),
                // client
                t.variableDeclaration(
                  'const',
                  [
                    t.variableDeclarator(
                      t.identifier('client'),
                      t.awaitExpression(t.callExpression(
                        t.memberExpression(
                          t.identifier('SigningStargateClient'),
                          t.identifier('connectWithSigner')
                        ),
                        [
                          t.identifier('rpcEndpoint'),
                          t.identifier('signer'),
                          t.objectExpression([
                            t.objectProperty(
                              t.identifier('registry'),
                              t.identifier('registry'),
                              false,
                              true
                            ),
                            t.objectProperty(
                              t.identifier('aminoTypes'),
                              t.identifier('aminoTypes'),
                              false,
                              true
                            ),
                          ])

                        ]
                      ))
                    )
                  ]
                ),

                // return 
                t.returnStatement(t.identifier('client'))
              ]
            ),
            true
          )
        )
      ]
    )
  )
};

export const createStargateClientAminoRegistry = ({
  aminos,
  aminoConverters,
  context
}: CreateStargateClientAminoConverters) => {
  return t.exportNamedDeclaration(
    t.variableDeclaration(
      'const',
      [
        t.variableDeclarator(
          t.identifier(aminoConverters),
          t.objectExpression([
            ...aminos.map(pkg =>
              t.spreadElement(
                memberExpressionOrIdentifier(
                  `${pkg}.AminoConverter`.split('.').reverse()
                )
              )
            )
          ])
        )
      ]
    )
  );
};

export const createStargateClientProtoRegistry = ({ registries, protoTypeRegistry, context }: CreateStargateClientProtoRegistry) => {

  context.addUtil('GeneratedType');

  return t.exportNamedDeclaration(
    t.variableDeclaration(
      'const',
      [
        t.variableDeclarator(
          identifier(protoTypeRegistry,
            t.tsTypeAnnotation(
              t.tsTypeReference(
                t.identifier('ReadonlyArray'),
                t.tsTypeParameterInstantiation(
                  [
                    t.tsTupleType([
                      t.tsStringKeyword(),
                      t.tsTypeReference(
                        t.identifier('GeneratedType')
                      )
                    ])
                  ]
                )
              )
            )
          ),
          t.arrayExpression([
            ...registries.map(pkg =>
              t.spreadElement(
                memberExpressionOrIdentifier(
                  `${pkg}.registry`.split('.').reverse()
                )
              )
            )
          ])
        )
      ]
    )
  );
};

export const createStargateClientOptions = ({
  name,
  aminoConverters,
  protoTypeRegistry,
  context
}: CreateStargateClientOptions) => {

  const includeDefaults = context.pluginValue('stargateClients.includeCosmosDefaultTypes');

  if (includeDefaults) {
    context.addUtil('defaultRegistryTypes');
  }

  context.addUtil('GeneratedType');
  context.addUtil('Registry');
  context.addUtil('AminoTypes');
  context.addUtil('SigningStargateClient');

  const prop = t.tsPropertySignature(
    t.identifier('defaultTypes'),
    t.tsTypeAnnotation(
      t.tsTypeReference(t.identifier('ReadonlyArray'), t.tsTypeParameterInstantiation(
        [
          t.tsTupleType([
            t.tsStringKeyword(),
            t.tsTypeReference(
              t.identifier('GeneratedType')
            )]
          )
        ]
      ))
    )
  );
  prop.optional = true;

  let funcArgs = [

  ];
  if (includeDefaults) {
    const args = [
      t.objectProperty(
        t.identifier('defaultTypes'),
        t.assignmentPattern(
          t.identifier('defaultTypes'),
          t.identifier('defaultRegistryTypes')
        ),
        false,
        true
      )
    ];
    //
    const funcTypes = t.tsTypeAnnotation(
      t.tsTypeLiteral(
        [
          prop
        ].filter(Boolean)
      )
    );

    funcArgs = [
      t.assignmentPattern(
        objectPattern(
          args,
          funcTypes
        ),
        t.objectExpression([])
      )
    ];
  }

  return t.exportNamedDeclaration(
    t.variableDeclaration(
      'const',
      [
        t.variableDeclarator(
          t.identifier(name),
          arrowFunctionExpression(
            funcArgs,
            t.blockStatement(
              [
                t.variableDeclaration(
                  'const',
                  [
                    t.variableDeclarator(
                      t.identifier('registry'),
                      t.newExpression(
                        t.identifier('Registry'),
                        [
                          t.arrayExpression(
                            [
                              includeDefaults && t.spreadElement(
                                t.identifier('defaultTypes')
                              ),
                              t.spreadElement(
                                t.identifier(protoTypeRegistry)
                              )
                            ].filter(Boolean)
                          )
                        ]
                      )
                    )
                  ]
                ),
                // amino
                t.variableDeclaration(
                  'const',
                  [
                    t.variableDeclarator(
                      t.identifier('aminoTypes'),
                      t.newExpression(
                        t.identifier('AminoTypes'),
                        [
                          t.objectExpression(
                            [
                              t.spreadElement(
                                t.identifier(aminoConverters)
                              )
                            ]
                          )
                        ]
                      )
                    )
                  ]
                ),

                // NEW CODE
                // return 
                t.returnStatement(t.objectExpression([
                  t.objectProperty(
                    t.identifier('registry'),
                    t.identifier('registry'),
                    false,
                    true
                  ),
                  t.objectProperty(
                    t.identifier('aminoTypes'),
                    t.identifier('aminoTypes'),
                    false,
                    true
                  )

                ])),
              ]
            ),
            t.tsTypeAnnotation(
              t.tsTypeLiteral([
                t.tsPropertySignature(
                  t.identifier('registry'),
                  t.tsTypeAnnotation(
                    t.tsTypeReference(
                      t.identifier('Registry')
                    )
                  )
                ),
                t.tsPropertySignature(
                  t.identifier('aminoTypes'),
                  t.tsTypeAnnotation(
                    t.tsTypeReference(
                      t.identifier('AminoTypes')
                    )
                  )
                )
              ])
            ),
            false
          )
        )
      ]
    )
  )
};
