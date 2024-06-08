import * as t from '@babel/types';
import { GenericParseContext } from '../../../encoding';
import { objectPattern } from '../../../utils';

export const rpcFuncArguments = (): t.ObjectPattern[] => {
    return [
        objectPattern(
            [
                t.objectProperty(
                    t.identifier('rpc'),
                    t.identifier('rpc'),
                    false,
                    true
                )
            ],
            t.tsTypeAnnotation(
                t.tsTypeLiteral(
                    [
                        t.tsPropertySignature(
                            t.identifier('rpc'),
                            t.tsTypeAnnotation(
                                t.tsTypeReference(
                                    t.identifier('Rpc')
                                )
                            )
                        )
                    ]
                )
            )
        )
    ];
};

export const rpcClassArguments = (): t.ObjectExpression[] => {
    return [
        t.objectExpression(
            [
                t.objectProperty(
                    t.identifier('rpc'),
                    t.identifier('rpc'),
                    false,
                    true
                )
            ]
        )
    ];
};

export const rpcNewAwaitImport = (
    path: string,
    className: string
) => {
    return t.newExpression(
        t.memberExpression(
            t.awaitExpression(
                t.callExpression(
                    t.import(),
                    [
                        t.stringLiteral(
                            path
                        )
                    ]
                )
            ),
            t.identifier(className),
            false
        ),
        [
            t.identifier('rpc')
        ]
    );
}

export const rpcNewTmAwaitImport = (
    path: string,
    className: string
) => {
    return t.callExpression(
        t.memberExpression(
            t.awaitExpression(
                t.callExpression(
                    t.import(),
                    [
                        t.stringLiteral(
                            path
                        )
                    ]
                )
            ),
            t.identifier(className),
            false
        ),
        [
            t.identifier('client')
        ]
    )

}

export const rpcRecursiveObjectProps = (
    names: string[],
    leaf?: any
) => {
    const [name, ...rest] = names;

    let baseComponent;
    if (names.length === 1) {
        baseComponent = leaf ? leaf : t.identifier(name)
    } else {
        baseComponent = rpcRecursiveObjectProps(rest, leaf)
    }

    return t.objectExpression([
        t.objectProperty(
            t.identifier(name),
            baseComponent
        )
    ])
};

export const rpcNestedImportObject = (
    obj: object,
    className: string
) => {

    if (typeof obj === 'string') {
        return rpcNewAwaitImport(obj, className);
    }

    const keys = Object.keys(obj);

    return t.objectExpression(keys.map(name => {
        return t.objectProperty(
            t.identifier(name),
            rpcNestedImportObject(obj[name], className)
        )
    }))
};

export const rpcTmNestedImportObject = (
    obj: object,
    className: string
) => {

    if (typeof obj === 'string') {
        return rpcNewTmAwaitImport(obj, className);
    }

    const keys = Object.keys(obj);

    return t.objectExpression(keys.map(name => {
        return t.objectProperty(
            t.identifier(name),
            rpcTmNestedImportObject(obj[name], className)
        )
    }))
};

export const createScopedRpcFactory = (
    obj: object,
    identifier: string,
    className: string
) => {
    return t.exportNamedDeclaration(
        t.variableDeclaration(
            'const',
            [
                t.variableDeclarator(
                    t.identifier(identifier),
                    t.arrowFunctionExpression(
                        rpcFuncArguments(),
                        //
                        rpcNestedImportObject(
                            obj,
                            className
                        ),
                        true
                    )
                )
            ]
        )
    )
}

export const createScopedRpcTmFactory = (
    context: GenericParseContext,
    obj: object,
    identifier: string
) => {

    context.addUtil('Tendermint34Client');
    context.addUtil('HttpEndpoint');
    context.addUtil('QueryClient');

    return t.exportNamedDeclaration(
        t.variableDeclaration(
            'const',
            [
                t.variableDeclarator(
                    // createRPCQueryClient
                    t.identifier(identifier),
                    t.arrowFunctionExpression(
                        [
                            objectPattern([
                                t.objectProperty(
                                    t.identifier('rpcEndpoint'),
                                    t.identifier('rpcEndpoint'),
                                    false,
                                    true
                                )
                            ], t.tsTypeAnnotation(
                                t.tsTypeLiteral([
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
                                    )
                                ])
                            ))
                        ],
                        t.blockStatement([

                            t.variableDeclaration('const', [
                                t.variableDeclarator(
                                    t.identifier('tmClient'),
                                    t.awaitExpression(
                                        t.callExpression(
                                            t.memberExpression(
                                                t.identifier('Tendermint34Client'),
                                                t.identifier('connect')
                                            ),
                                            [
                                                t.identifier('rpcEndpoint')
                                            ]
                                        )
                                    )
                                )
                            ]),
                            /////
                            t.variableDeclaration('const', [
                                t.variableDeclarator(
                                    t.identifier('client'),
                                    t.newExpression(
                                        t.identifier('QueryClient'),
                                        [
                                            t.identifier('tmClient')
                                        ]
                                    )
                                )
                            ]),

                            t.returnStatement(
                                rpcTmNestedImportObject(
                                    obj,
                                    'createRpcQueryExtension'
                                )
                            )

                        ]),
                        true
                    )
                )
            ]
        )
    )
}
