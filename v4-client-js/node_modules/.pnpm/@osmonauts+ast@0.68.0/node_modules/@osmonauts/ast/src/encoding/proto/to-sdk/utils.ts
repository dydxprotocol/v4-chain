import * as t from '@babel/types';
import { getFieldNames } from '../../types';
import { ToSDKMethod } from './index';

const notUndefinedSetValue = (sdkName: string, msgName: string, expr: t.Expression) => {
    return t.expressionStatement(
        t.logicalExpression(
            '&&',
            t.binaryExpression(
                '!==',
                t.memberExpression(
                    t.identifier('message'),
                    t.identifier(msgName)
                ),
                t.identifier('undefined')
            ),
            t.assignmentExpression(
                '=',
                t.memberExpression(
                    t.identifier('obj'),
                    t.identifier(sdkName)
                ),
                expr
            )
        )
    );
}

export const toSDK = {

    scalar(args: ToSDKMethod) {
        const { propName, origName } = getFieldNames(args.field);

        return notUndefinedSetValue(origName, propName, t.memberExpression(
            t.identifier('message'),
            t.identifier(propName)
        ));
    },

    string(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    double(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    float(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    bool(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },

    number(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },

    int32(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },

    uint32(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },

    sint32(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    fixed32(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    sfixed32(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    long(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    int64(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    uint64(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    sint64(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    fixed64(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },
    sfixed64(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },

    type(args: ToSDKMethod) {
        const { propName, origName } = getFieldNames(args.field);
        const name = args.context.getTypeName(args.field);

        // TODO isn't the nested conditional a waste? (using ts-proto as reference)
        // maybe null is OK?
        return notUndefinedSetValue(origName, propName, t.conditionalExpression(
            t.memberExpression(
                t.identifier('message'),
                t.identifier(propName)
            ),
            t.callExpression(
                t.memberExpression(
                    t.identifier(name),
                    t.identifier('toSDK')
                ),
                [
                    t.memberExpression(
                        t.identifier('message'),
                        t.identifier(propName)
                    )
                ]
            ),
            t.identifier('undefined')
        ));
    },

    enum(args: ToSDKMethod) {
        const { propName, origName } = getFieldNames(args.field);

        const enumFuncName = args.context.getToEnum(args.field);
        return notUndefinedSetValue(origName, propName, t.callExpression(
            t.identifier(enumFuncName),
            [
                t.memberExpression(
                    t.identifier('message'),
                    t.identifier(propName)
                )
            ]
        ));
    },

    bytes(args: ToSDKMethod) {
        return toSDK.scalar(args);
    },

    duration(args: ToSDKMethod) {
        return toSDK.type(args);
    },

    timestamp(args: ToSDKMethod) {
        return toSDK.type(args);
    },

    keyHash(args: ToSDKMethod) {

        const { propName, origName } = getFieldNames(args.field);
        const keyType = args.field.keyType;
        const valueType = args.field.parsedType.name;

        let toSDK = null;
        switch (valueType) {
            case 'string':
                toSDK = t.identifier('v')
                break;
            case 'uint32':
            case 'int32':
                toSDK = t.callExpression(
                    t.memberExpression(
                        t.identifier('Math'),
                        t.identifier('round')
                    ),
                    [
                        t.identifier('v')
                    ]
                )
                break;
            case 'int64':
            case 'uint64':
                toSDK = t.callExpression(
                    t.memberExpression(
                        t.identifier('v'),
                        t.identifier('toString')
                    ),
                    []
                )
                break;
            default:
                toSDK = t.callExpression(
                    t.memberExpression(
                        t.identifier(valueType),
                        t.identifier('toSDK')
                    ),
                    [
                        t.identifier('v')
                    ]
                )
        }


        return [
            t.expressionStatement(
                t.assignmentExpression(
                    '=',
                    t.memberExpression(
                        t.identifier('obj'),
                        t.identifier(origName)
                    ),
                    t.objectExpression([])
                )
            ),
            //
            t.ifStatement(
                t.memberExpression(
                    t.identifier('message'),
                    t.identifier(propName)
                ),
                t.blockStatement([
                    t.expressionStatement(
                        t.callExpression(
                            t.memberExpression(
                                t.callExpression(
                                    t.memberExpression(
                                        t.identifier('Object'),
                                        t.identifier('entries')
                                    ),
                                    [
                                        t.memberExpression(
                                            t.identifier('message'),
                                            t.identifier(propName)
                                        )
                                    ]
                                ),
                                t.identifier('forEach')
                            ),
                            [
                                t.arrowFunctionExpression(
                                    [
                                        t.arrayPattern(
                                            [
                                                t.identifier('k'),
                                                t.identifier('v')
                                            ]
                                        )
                                    ],
                                    t.blockStatement([
                                        t.expressionStatement(
                                            t.assignmentExpression(
                                                '=',
                                                t.memberExpression(
                                                    t.memberExpression(
                                                        t.identifier('obj'),
                                                        t.identifier(origName)
                                                    ),
                                                    t.identifier('k'),
                                                    true
                                                ),
                                                toSDK
                                            )
                                        )
                                    ])
                                )
                            ]
                        )
                    )
                ])
            )
        ]
    },

    array(args: ToSDKMethod, expr: t.Expression) {
        const { propName, origName } = getFieldNames(args.field);

        return t.ifStatement(
            t.memberExpression(
                t.identifier('message'),
                t.identifier(propName)
            ),
            t.blockStatement([
                t.expressionStatement(
                    t.assignmentExpression(
                        '=',
                        t.memberExpression(
                            t.identifier('obj'),
                            t.identifier(origName)
                        ),
                        t.callExpression(
                            t.memberExpression(
                                t.memberExpression(
                                    t.identifier('message'),
                                    t.identifier(propName)
                                ),
                                t.identifier('map')
                            ),
                            [
                                t.arrowFunctionExpression(
                                    [
                                        t.identifier('e')
                                    ],
                                    expr
                                )
                            ]
                        )
                    )
                )
            ]),
            t.blockStatement([
                t.expressionStatement(
                    t.assignmentExpression(
                        '=',
                        t.memberExpression(
                            t.identifier('obj'),
                            t.identifier(origName)
                        ),
                        t.arrayExpression([])
                    )
                )
            ])
        );
    }

};

export const arrayTypes = {
    scalar() {
        return t.identifier('e');
    },
    string() {
        return arrayTypes.scalar();
    },
    double() {
        return arrayTypes.scalar();
    },
    float() {
        return arrayTypes.scalar();
    },
    bool() {
        return arrayTypes.scalar();
    },
    number() {
        return arrayTypes.scalar();
    },
    int32() {
        return arrayTypes.number();
    },
    uint32() {
        return arrayTypes.number();
    },
    sint32() {
        return arrayTypes.number();
    },
    fixed32() {
        return arrayTypes.number();
    },
    sfixed32() {
        return arrayTypes.number();
    },
    long(args: ToSDKMethod) {
        return arrayTypes.scalar();
    },
    int64(args: ToSDKMethod) {
        return arrayTypes.long(args);
    },
    uint64(args: ToSDKMethod) {
        return arrayTypes.long(args);
    },
    sint64(args: ToSDKMethod) {
        return arrayTypes.long(args);
    },
    fixed64(args: ToSDKMethod) {
        return arrayTypes.long(args);
    },
    sfixed64(args: ToSDKMethod) {
        return arrayTypes.long(args);
    },
    bytes(args: ToSDKMethod) {
        return arrayTypes.scalar();
    },
    enum(args: ToSDKMethod) {
        const enumFuncName = args.context.getToEnum(args.field);
        return t.callExpression(
            t.identifier(enumFuncName),
            [
                t.identifier('e')
            ]
        );
    },
    type(args: ToSDKMethod) {
        const name = args.context.getTypeName(args.field);
        return t.conditionalExpression(
            t.identifier('e'),
            t.callExpression(
                t.memberExpression(
                    t.identifier(name),
                    t.identifier('toSDK')
                ),
                [
                    t.identifier('e')
                ]
            ),
            t.identifier('undefined')
        );
    }
}

