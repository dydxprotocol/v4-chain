import * as t from '@babel/types';
import { ProtoField } from '@osmonauts/types';
import { getDefaultTSTypeFromProtoType } from '../../types';
import { ToJSONMethod } from './index';

const notUndefinedSetValue = (messageProp: string, objProp: string, expr: t.Expression) => {
    return t.expressionStatement(
        t.logicalExpression(
            '&&',
            t.binaryExpression(
                '!==',
                t.memberExpression(
                    t.identifier('message'),
                    t.identifier(messageProp)
                ),
                t.identifier('undefined')
            ),
            t.assignmentExpression(
                '=',
                t.memberExpression(
                    t.identifier('obj'),
                    t.identifier(objProp)
                ),
                expr
            )
        )
    );
}

const getPropNames = (field: ProtoField) => {
    const messageProp = field.name;
    const objProp = field.options?.json_name ?? field.name;
    return {
        messageProp,
        objProp
    }
}

export const toJSON = {

    //  message.sender !== undefined && (obj.sender = message.sender);
    identity(args: ToJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        return notUndefinedSetValue(messageProp, objProp, t.memberExpression(
            t.identifier('message'),
            t.identifier(messageProp)
        ));
    },

    string(args: ToJSONMethod) {
        return toJSON.identity(args);
    },
    double(args: ToJSONMethod) {
        return toJSON.identity(args);
    },
    float(args: ToJSONMethod) {
        return toJSON.identity(args);
    },
    bool(args: ToJSONMethod) {
        return toJSON.identity(args);
    },

    // message.maxDepth !== undefined && (obj.maxDepth = Math.round(message.maxDepth));
    number(args: ToJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        return notUndefinedSetValue(
            messageProp,
            objProp,
            t.callExpression(
                t.memberExpression(
                    t.identifier('Math'),
                    t.identifier('round')
                ),
                [
                    t.memberExpression(
                        t.identifier('message'),
                        t.identifier(messageProp)
                    )
                ]
            )
        );
    },
    // message.maxDepth !== undefined && (obj.maxDepth = Math.round(message.maxDepth));
    int32(args: ToJSONMethod) {
        return toJSON.number(args);
    },

    uint32(args: ToJSONMethod) {
        return toJSON.number(args);
    },

    sint32(args: ToJSONMethod) {
        return toJSON.number(args);
    },
    fixed32(args: ToJSONMethod) {
        return toJSON.number(args);
    },
    sfixed32(args: ToJSONMethod) {
        return toJSON.number(args);
    },

    // message.poolId !== undefined && (obj.poolId = (message.poolId || Long.UZERO).toString());
    // message.poolId !== undefined && (obj.poolId = (message.poolId || undefined).toString());
    long(args: ToJSONMethod) {
        args.context.addUtil('Long');
        const { messageProp, objProp } = getPropNames(args.field);
        return notUndefinedSetValue(
            messageProp,
            objProp,
            t.callExpression(
                t.memberExpression(
                    t.logicalExpression(
                        '||',
                        t.memberExpression(
                            t.identifier('message'),
                            t.identifier(messageProp)
                        ),
                        getDefaultTSTypeFromProtoType(args.context, args.field, args.isOneOf)
                    ),
                    t.identifier('toString')
                ),
                []
            )
        );
    },

    int64(args: ToJSONMethod) {
        return toJSON.long(args);
    },
    uint64(args: ToJSONMethod) {
        return toJSON.long(args);
    },
    sint64(args: ToJSONMethod) {
        return toJSON.long(args);
    },
    fixed64(args: ToJSONMethod) {
        return toJSON.long(args);
    },
    sfixed64(args: ToJSONMethod) {
        return toJSON.long(args);
    },

    // message.signDoc !== undefined && (obj.signDoc = message.signDoc ? SignDocDirectAux.toJSON(message.signDoc) : undefined);
    type(args: ToJSONMethod) {
        const name = args.context.getTypeName(args.field);
        const { messageProp, objProp } = getPropNames(args.field);
        // TODO isn't the nested conditional a waste? (using ts-proto as reference)
        // maybe null is OK?
        return notUndefinedSetValue(messageProp, objProp, t.conditionalExpression(
            t.memberExpression(
                t.identifier('message'),
                t.identifier(messageProp)
            ),
            t.callExpression(
                t.memberExpression(
                    t.identifier(name),
                    t.identifier('toJSON')
                ),
                [
                    t.memberExpression(
                        t.identifier('message'),
                        t.identifier(messageProp)
                    )
                ]
            ),
            t.identifier('undefined')
        ));
    },

    // message.mode !== undefined && (obj.mode = signModeToJSON(message.mode));
    enum(args: ToJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        const enumFuncName = args.context.getToEnum(args.field);
        return notUndefinedSetValue(messageProp, objProp, t.callExpression(
            t.identifier(enumFuncName),
            [
                t.memberExpression(
                    t.identifier('message'),
                    t.identifier(messageProp)
                )
            ]
        ));
    },

    // TODO again, another ts-proto reference that does not necessarily make sense
    // message.queryData !== undefined && (obj.queryData = base64FromBytes(message.queryData !== undefined ? message.queryData : new Uint8Array()));
    // message.queryData !== undefined && (obj.queryData = base64FromBytes(message.queryData !== undefined ? message.queryData : undefined));
    bytes(args: ToJSONMethod) {
        args.context.addUtil('base64FromBytes');
        const { messageProp, objProp } = getPropNames(args.field);

        let expr;
        if (args.isOptional) {
            // message.bytesValue !== undefined && (obj.bytesValue = message.bytesValue !== undefined ? base64FromBytes(message.bytesValue) : undefined);

            expr = t.conditionalExpression(
                t.binaryExpression(
                    '!==',
                    t.memberExpression(
                        t.identifier('message'),
                        t.identifier(messageProp)
                    ),
                    t.identifier('undefined')
                ),
                t.callExpression(
                    t.identifier('base64FromBytes'),
                    [
                        t.memberExpression(
                            t.identifier('message'),
                            t.identifier(messageProp)
                        )
                    ]
                ),
                t.identifier('undefined')
            );
        } else {
            // message.queryData !== undefined && (obj.queryData = base64FromBytes(message.queryData !== undefined ? message.queryData : new Uint8Array()));
            expr = t.callExpression(
                t.identifier('base64FromBytes'),
                [
                    t.conditionalExpression(
                        t.binaryExpression(
                            '!==',
                            t.memberExpression(
                                t.identifier('message'),
                                t.identifier(messageProp)
                            ),
                            t.identifier('undefined')
                        ),
                        t.memberExpression(
                            t.identifier('message'),
                            t.identifier(messageProp)
                        ),
                        getDefaultTSTypeFromProtoType(args.context, args.field, args.isOneOf)
                    )
                ]
            )
        }
        return notUndefinedSetValue(messageProp, objProp, expr);
    },

    // message.period !== undefined && (obj.period = message.period);

    duration(args: ToJSONMethod) {
        const durationFormat = args.context.pluginValue('prototypes.typingsFormat.duration');
        switch (durationFormat) {
            case 'string':
                return toJSON.durationString(args);
            case 'duration':
            default:
                return toJSON.type(args);
        }
    },

    durationString(args: ToJSONMethod) {
        return toJSON.identity(args);
    },

    timestamp(args: ToJSONMethod) {
        const timestampFormat = args.context.pluginValue('prototypes.typingsFormat.timestamp')
        switch (timestampFormat) {
            case 'timestamp':
                return toJSON.timestampTimestamp(args);
            case 'date':
            default:
                return toJSON.timestampDate(args);
        }
    },

    // message.periodReset !== undefined && (obj.periodReset = fromTimestamp(message.periodReset).toISOString());

    timestampTimestamp(args: ToJSONMethod) {
        args.context.addUtil('fromTimestamp');
        const { messageProp, objProp } = getPropNames(args.field);
        return notUndefinedSetValue(messageProp, objProp, t.callExpression(
            t.memberExpression(
                t.callExpression(
                    t.identifier('fromTimestamp'), [
                    t.memberExpression(
                        t.identifier('message'),
                        t.identifier(messageProp)
                    )
                ]),
                t.identifier('toISOString')
            ),
            []
        ));
    },

    // message.periodReset !== undefined && (obj.periodReset = message.periodReset.toISOString());

    timestampDate(args: ToJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        return notUndefinedSetValue(messageProp, objProp, t.callExpression(
            t.memberExpression(
                t.memberExpression(
                    t.identifier('message'),
                    t.identifier(messageProp)
                ),
                t.identifier('toISOString')
            ),
            []
        ));
    },

    // obj.labels = {};

    //   if (message.labels) {
    //     Object.entries(message.labels).forEach(([k, v]) => {
    //       obj.labels[k] = v;
    //     });
    //   }


    // obj.typeMap = {};

    // if (message.typeMap) {
    //   Object.entries(message.typeMap).forEach(([k, v]) => {
    //     obj.typeMap[k] = Type.toJSON(v);
    //   });
    // }

    keyHash(args: ToJSONMethod) {

        const { messageProp, objProp } = getPropNames(args.field);
        const keyType = args.field.keyType;
        const valueType = args.field.parsedType.name;

        let toJSON = null;
        switch (valueType) {
            case 'string':
                toJSON = t.identifier('v')
                break;
            case 'uint32':
            case 'int32':
                toJSON = t.callExpression(
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
                toJSON = t.callExpression(
                    t.memberExpression(
                        t.identifier('v'),
                        t.identifier('toString')
                    ),
                    []
                )
                break;
            default:
                toJSON = t.callExpression(
                    t.memberExpression(
                        t.identifier(valueType),
                        t.identifier('toJSON')
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
                        t.identifier(objProp)
                    ),
                    t.objectExpression([])
                )
            ),
            //
            t.ifStatement(
                t.memberExpression(
                    t.identifier('message'),
                    t.identifier(messageProp)
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
                                            t.identifier(messageProp)
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
                                                        t.identifier(objProp)
                                                    ),
                                                    t.identifier('k'),
                                                    true
                                                ),
                                                toJSON
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

    // if (message.codeIds) {
    //     obj.codeIds = message.codeIds.map(e => (e || Long.UZERO).toString());
    // } else {
    //     obj.codeIds = [];
    // }

    array(args: ToJSONMethod, expr: t.Expression) {
        const { messageProp, objProp } = getPropNames(args.field);
        return t.ifStatement(
            t.memberExpression(
                t.identifier('message'),
                t.identifier(messageProp)
            ),
            t.blockStatement([
                t.expressionStatement(
                    t.assignmentExpression(
                        '=',
                        t.memberExpression(
                            t.identifier('obj'),
                            t.identifier(objProp)
                        ),
                        t.callExpression(
                            t.memberExpression(
                                t.memberExpression(
                                    t.identifier('message'),
                                    t.identifier(messageProp)
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
                            t.identifier(objProp)
                        ),
                        t.arrayExpression([])
                    )
                )
            ])
        );
    }

};

export const arrayTypes = {

    identity() {
        return t.identifier('e');
    },


    // if (message.overloadId) {
    //     obj.overloadId = message.overloadId.map(e => e);
    // } else {
    //     obj.overloadId = [];
    // }

    string() {
        return arrayTypes.identity();
    },
    double() {
        return arrayTypes.identity();
    },
    float() {
        return arrayTypes.identity();
    },
    bool() {
        return arrayTypes.identity();
    },

    //   if (message.lineOffsets) {
    //     obj.lineOffsets = message.lineOffsets.map(e => Math.round(e));
    //   } else {
    //     obj.lineOffsets = [];
    //   }

    number() {
        return t.callExpression(
            t.memberExpression(
                t.identifier('Math'),
                t.identifier('round')
            ),
            [
                t.identifier('e')
            ]
        )
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


    // if (message.codeIds) {
    //     obj.codeIds = message.codeIds.map(e => (e || Long.UZERO).toString());
    // } else {
    //     obj.codeIds = [];
    // }

    long(args: ToJSONMethod) {
        return t.callExpression(
            t.memberExpression(
                t.logicalExpression(
                    '||',
                    t.identifier('e'),
                    getDefaultTSTypeFromProtoType(args.context, {
                        ...args.field,
                        rule: undefined, // so it's treated as type not an array...
                    }, args.isOneOf)
                ),
                t.identifier('toString')
            ),
            []
        )
    },

    int64(args: ToJSONMethod) {
        return arrayTypes.long(args);
    },
    uint64(args: ToJSONMethod) {
        return arrayTypes.long(args);
    },
    sint64(args: ToJSONMethod) {
        return arrayTypes.long(args);
    },
    fixed64(args: ToJSONMethod) {
        return arrayTypes.long(args);
    },
    sfixed64(args: ToJSONMethod) {
        return arrayTypes.long(args);
    },

    //   if (message.myBytesArray) {
    //     obj.myBytesArray = message.myBytesArray.map(e => base64FromBytes(e !== undefined ? e : new Uint8Array()));
    //   } else {
    //     obj.myBytesArray = [];
    //   }

    bytes(args: ToJSONMethod) {
        args.context.addUtil('base64FromBytes');
        return t.callExpression(
            t.identifier('base64FromBytes'),
            [
                t.conditionalExpression(
                    t.binaryExpression(
                        '!==',
                        t.identifier('e'),
                        t.identifier('undefined')
                    ),
                    t.identifier('e'),
                    getDefaultTSTypeFromProtoType(args.context, {
                        ...args.field,
                        rule: undefined, // so it's treated as type not an array...
                    }, args.isOneOf)
                )
            ]
        );
    },

    enum(args: ToJSONMethod) {
        const enumFuncName = args.context.getToEnum(args.field);
        return t.callExpression(
            t.identifier(enumFuncName),
            [
                t.identifier('e')
            ]
        );
    },

    // if (message.tokenInMaxs) {
    //     obj.tokenInMaxs = message.tokenInMaxs.map(e => e ? Coin.toJSON(e) : undefined);
    // } else {
    //     obj.tokenInMaxs = [];
    // }

    type(args: ToJSONMethod) {
        const name = args.context.getTypeName(args.field);
        return t.conditionalExpression(
            t.identifier('e'),
            t.callExpression(
                t.memberExpression(
                    t.identifier(name),
                    t.identifier('toJSON')
                ),
                [
                    t.identifier('e')
                ]
            ),
            t.identifier('undefined')
        );
    }
}

