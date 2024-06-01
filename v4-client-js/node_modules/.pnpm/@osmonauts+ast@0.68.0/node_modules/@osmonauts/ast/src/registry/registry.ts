import * as t from '@babel/types';
import { AminoParseContext } from '../encoding';
import { identifier } from '../utils';

export interface ServiceMethod {
  methodName: string;
  typeUrl: string;
  TypeName: string;
}

export const createTypeRegistryObject = (mutation: ServiceMethod) => {
  return t.objectProperty(
    t.stringLiteral(mutation.typeUrl),
    t.identifier(mutation.TypeName)
  );
};

export const createTypeRegistry = (context: AminoParseContext, mutations: ServiceMethod[]) => {
  context.addUtil('GeneratedType');

  return t.exportNamedDeclaration(
    t.variableDeclaration('const', [
      t.variableDeclarator(identifier(
        'registry',
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
      ), t.arrayExpression(
        [
          ...mutations.map(mutation => t.arrayExpression(
            [
              t.stringLiteral(mutation.typeUrl),
              t.identifier(mutation.TypeName)
            ]
          ))
        ]
      ))
    ]));

};

export const createRegistryLoader = (context: AminoParseContext) => {
  context.addUtil('Registry');

  return t.exportNamedDeclaration(t.variableDeclaration(
    'const',
    [
      t.variableDeclarator(
        t.identifier('load'),
        t.arrowFunctionExpression(
          [
            identifier('protoRegistry', t.tsTypeAnnotation(
              t.tsTypeReference(
                t.identifier('Registry')
              )
            ))
          ],
          t.blockStatement(
            [
              t.expressionStatement(
                t.callExpression(
                  t.memberExpression(
                    t.identifier('registry'),
                    t.identifier('forEach')
                  ),
                  [
                    t.arrowFunctionExpression(
                      [
                        t.arrayPattern([
                          t.identifier('typeUrl'),
                          t.identifier('mod')
                        ])
                      ],
                      t.blockStatement(
                        [
                          t.expressionStatement(
                            t.callExpression(
                              t.memberExpression(
                                t.identifier(
                                  'protoRegistry'
                                ),
                                t.identifier('register')
                              ),
                              [
                                t.identifier('typeUrl'),
                                t.identifier('mod')
                              ]
                            )
                          )
                        ]
                      )
                    )
                  ]
                )
              )
            ]
          )
        )
      )
    ]
  ))
};