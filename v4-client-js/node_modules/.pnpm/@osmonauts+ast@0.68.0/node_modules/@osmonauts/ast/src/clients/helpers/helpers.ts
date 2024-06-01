import * as t from '@babel/types';
import { identifier } from '../../utils';
import { Mutation } from '../../types';
import { camel } from '@osmonauts/utils';
import { GenericParseContext } from '../../encoding';

export const addMsgMethod = ({ methodName, typeUrl, TypeName, methodToCall }) => {
    return t.objectMethod('method', t.identifier(methodName), [
        identifier('value', t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(TypeName)
            )
        ))
    ], t.blockStatement([
        t.returnStatement(
            t.objectExpression([
                t.objectProperty(
                    t.identifier('typeUrl'),
                    t.stringLiteral(typeUrl)
                ),
                t.objectProperty(
                    t.identifier('value'),
                    t.callExpression(
                        t.memberExpression(
                            t.identifier(TypeName),
                            t.identifier(methodToCall)
                        ),
                        [t.identifier('value')]
                    )
                )
            ])
        )]));
};

export const addFromJSONMethod = ({ methodName, typeUrl, TypeName }) => {
    return t.objectMethod('method', t.identifier(methodName), [
        identifier('value', t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier('any')
            )
        ))
    ], t.blockStatement([
        t.returnStatement(
            t.objectExpression([
                t.objectProperty(
                    t.identifier('typeUrl'),
                    t.stringLiteral(typeUrl)
                ),
                t.objectProperty(
                    t.identifier('value'),
                    t.callExpression(
                        t.memberExpression(
                            t.identifier(TypeName),
                            t.identifier('fromJSON')
                        ),
                        [t.identifier('value')]
                    )
                )
            ])
        )]));
};

export const addFromPartialMethod = ({ methodName, typeUrl, TypeName }) => {
    return addMsgMethod({ methodName, typeUrl, TypeName, methodToCall: 'fromPartial' });
};

export const addToJSONMethod = ({ methodName, typeUrl, TypeName }) => {
    return addMsgMethod({ methodName, typeUrl, TypeName, methodToCall: 'toJSON' });
};


export const addJsonMethod = ({ methodName, typeUrl, TypeName }) => {
    return t.objectMethod('method', t.identifier(methodName), [
        identifier('value', t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(TypeName)
            )
        ))
    ], t.blockStatement([
        t.returnStatement(
            t.objectExpression([
                t.objectProperty(
                    t.identifier('typeUrl'),
                    t.stringLiteral(typeUrl)
                ),
                t.objectProperty(
                    t.identifier('value'),
                    t.identifier('value'),
                    false,
                    true
                )
            ])
        )]));
};

export const addEncodedMethod = ({ methodName, typeUrl, TypeName }) => {
    return t.objectMethod('method', t.identifier(methodName), [
        identifier('value', t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(TypeName)
            )
        ))
    ], t.blockStatement([
        t.returnStatement(
            t.objectExpression([

                t.objectProperty(
                    t.identifier('typeUrl'),
                    t.stringLiteral(typeUrl)
                ),

                t.objectProperty(
                    t.identifier('value'),
                    t.callExpression(
                        t.memberExpression(
                            t.callExpression(
                                t.memberExpression(
                                    t.identifier(TypeName),
                                    t.identifier('encode')
                                ),
                                [
                                    t.identifier('value')
                                ]
                            ),
                            t.identifier('finish')
                        ),
                        []
                    )
                )
            ])
        )]));
};

interface HelperObject {
    context: GenericParseContext;
    name: string;
    mutations: Mutation[];
}

export const createHelperObject = ({
    context, // for future plugins
    name,
    mutations
}: HelperObject) => {

    mutations = mutations.map(mutation => {
        return {
            ...mutation,
            methodName: camel(mutation.methodName)
        }
    });

    const methods = [
        context.pluginValue('prototypes.methods.encode') && (t.objectProperty(
            t.identifier('encoded'),
            t.objectExpression(
                mutations.map(mutation => addEncodedMethod(mutation))
            )
        )),
        t.objectProperty(
            t.identifier('withTypeUrl'),
            t.objectExpression(
                mutations.map(mutation => addJsonMethod(mutation))
            )
        ),
        context.pluginValue('prototypes.methods.toJSON') && (t.objectProperty(
            t.identifier('toJSON'),
            t.objectExpression(
                mutations.map(mutation => addToJSONMethod(mutation))
            )
        )),
        context.pluginValue('prototypes.methods.fromJSON') && (t.objectProperty(
            t.identifier('fromJSON'),
            t.objectExpression(
                mutations.map(mutation => addFromJSONMethod(mutation))
            )
        )),
        context.pluginValue('prototypes.methods.fromPartial') && (t.objectProperty(
            t.identifier('fromPartial'),
            t.objectExpression(
                mutations.map(mutation => addFromPartialMethod(mutation))
            )
        ))
    ].filter(Boolean);

    return t.exportNamedDeclaration(
        t.variableDeclaration(
            'const',
            [
                t.variableDeclarator(
                    t.identifier(name),
                    t.objectExpression(methods)
                )
            ]
        )
    )
}