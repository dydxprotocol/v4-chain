import * as t from '@babel/types';
import { GenericParseContext } from '../../../encoding';
import { objectPattern } from '../../../utils';

export const lcdArguments = (): t.ObjectProperty[] => {
    return [
        t.objectProperty(
            t.identifier('requestClient'),
            t.identifier('requestClient'),
            false,
            true
        )
    ];
};

export const lcdFuncArguments = (): t.ObjectPattern[] => {
    return [
        objectPattern(
            [
                t.objectProperty(
                    t.identifier('restEndpoint'),
                    t.identifier('restEndpoint'),
                    false,
                    true
                )
            ],
            t.tsTypeAnnotation(
                t.tsTypeLiteral([
                    t.tsPropertySignature(
                        t.identifier('restEndpoint'),
                        t.tsTypeAnnotation(
                            t.tsStringKeyword()
                        )
                    )
                ])
            )
        )
    ];
};

export const lcdClassArguments = (): t.ObjectExpression[] => {
    return [
        t.objectExpression(
            lcdArguments()
        )
    ];
};

export const lcdNewAwaitImport = (
    path: string,
    className: string,
    _arguments: t.ObjectExpression[]
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
        _arguments
    );
}

export const lcdRecursiveObjectProps = (
    names: string[],
    leaf?: any
) => {
    const [name, ...rest] = names;

    let baseComponent;
    if (names.length === 1) {
        baseComponent = leaf ? leaf : t.identifier(name)
    } else {
        baseComponent = lcdRecursiveObjectProps(rest, leaf)
    }

    return t.objectExpression([
        t.objectProperty(
            t.identifier(name),
            baseComponent
        )
    ])
};

export const lcdNestedImportObject = (
    obj: object,
    className: string,
    _arguments: t.ObjectExpression[]
) => {

    if (typeof obj === 'string') {
        return lcdNewAwaitImport(obj, className, _arguments);
    }

    const keys = Object.keys(obj);

    return t.objectExpression(keys.map(name => {
        return t.objectProperty(
            t.identifier(name),
            lcdNestedImportObject(obj[name], className, _arguments)
        )
    }))
};

export const createScopedLCDFactory = (
    context: GenericParseContext,
    obj: object,
    identifier: string,
    className: string
) => {

    context.addUtil('LCDClient');

    return t.exportNamedDeclaration(
        t.variableDeclaration(
            'const',
            [
                t.variableDeclarator(
                    t.identifier(identifier),
                    t.arrowFunctionExpression(
                        lcdFuncArguments(),
                        //

                        t.blockStatement([
                            t.variableDeclaration(
                                'const',
                                [
                                    t.variableDeclarator(
                                        t.identifier('requestClient'),
                                        t.newExpression(
                                            t.identifier('LCDClient'),
                                            [
                                                t.objectExpression(
                                                    [
                                                        t.objectProperty(
                                                            t.identifier('restEndpoint'),
                                                            t.identifier('restEndpoint'),
                                                            false,
                                                            true
                                                        )
                                                    ]
                                                )
                                            ]
                                        )
                                    )
                                ]
                            ),
                            ////
                            t.returnStatement(
                                lcdNestedImportObject(
                                    obj,
                                    className,
                                    lcdClassArguments()
                                )
                            ),
                        ]),
                        // lcdNestedImportObject(
                        //     obj,
                        //     className,
                        //     lcdClassArguments()
                        // ),
                        true
                    )
                )
            ]
        )
    )
};