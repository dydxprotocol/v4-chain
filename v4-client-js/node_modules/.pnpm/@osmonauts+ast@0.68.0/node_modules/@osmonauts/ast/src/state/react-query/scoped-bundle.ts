import * as t from '@babel/types';
import { GenericParseContext } from '../../encoding';
import { objectPattern } from '../../utils';
import { variableSlug } from '@osmonauts/utils';

export const rpcHookFuncArguments = (): t.ObjectPattern[] => {
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

export const rpcHookClassArguments = (): t.ObjectExpression[] => {
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

export const rpcHookNewTmRequire = (
    imports: HookImport[],
    path: string,
    methodName: string
) => {

    imports.push({
        as: variableSlug(path),
        path
    });

    return t.callExpression(
        t.memberExpression(
            t.identifier(variableSlug(path)),
            t.identifier(methodName)
        ),
        [
            t.identifier('rpc')
        ]
    )

}

export const rpcHookRecursiveObjectProps = (
    names: string[],
    leaf?: any
) => {
    const [name, ...rest] = names;

    let baseComponent;
    if (names.length === 1) {
        baseComponent = leaf ? leaf : t.identifier(name)
    } else {
        baseComponent = rpcHookRecursiveObjectProps(rest, leaf)
    }

    return t.objectExpression([
        t.objectProperty(
            t.identifier(name),
            baseComponent
        )
    ])
};

export const rpcHookTmNestedImportObject = (
    imports: HookImport[],
    obj: object,
    methodName: string
) => {

    if (typeof obj === 'string') {
        return rpcHookNewTmRequire(imports, obj, methodName);
    }

    const keys = Object.keys(obj);

    return t.objectExpression(keys.map(name => {
        return t.objectProperty(
            t.identifier(name),
            rpcHookTmNestedImportObject(imports, obj[name], methodName)
        )
    }))
};

interface HookImport {
    as: string;
    path: string;
}

export const createScopedRpcHookFactory = (
    context: GenericParseContext,
    obj: object,
    identifier: string
) => {

    context.addUtil('ProtobufRpcClient');

    const hookImports: HookImport[] = [];

    const ast = t.exportNamedDeclaration(
        t.variableDeclaration(
            'const',
            [
                t.variableDeclarator(
                    // createRPCQueryHooks
                    t.identifier(identifier),
                    t.arrowFunctionExpression(
                        [
                            objectPattern([
                                t.objectProperty(
                                    t.identifier('rpc'),
                                    t.identifier('rpc'),
                                    false,
                                    true
                                )
                            ], t.tsTypeAnnotation(
                                t.tsTypeLiteral([
                                    t.tsPropertySignature(
                                        t.identifier('rpc'),
                                        t.tsTypeAnnotation(
                                            t.tsUnionType([
                                                t.tsTypeReference(
                                                    t.identifier('ProtobufRpcClient')
                                                ),
                                                t.tsUndefinedKeyword()
                                            ])
                                        )
                                    )
                                ])
                            ))
                        ],
                        t.blockStatement([

                            t.returnStatement(
                                rpcHookTmNestedImportObject(
                                    hookImports,
                                    obj,
                                    'createRpcQueryHooks'
                                )
                            )

                        ]),
                        false
                    )
                )
            ]
        )
    );

    const imports = hookImports.map(hookport => {
        return {
            "type": "ImportDeclaration",
            "importKind": "value",
            "specifiers": [
                {
                    "type": "ImportNamespaceSpecifier",
                    "local": {
                        "type": "Identifier",
                        "name": hookport.as
                    }
                }
            ],
            "source": {
                "type": "StringLiteral",
                "value": hookport.path
            }
        };
    });

    return [...imports, ast];
};
