import * as t from '@babel/types';
import { FromJSONMethod } from './index';
import { callExpression, identifier } from '../../../utils';
import { getDefaultTSTypeFromProtoType } from '../../types';
import { ProtoField } from '@osmonauts/types';

const getPropNames = (field: ProtoField) => {
    const messageProp = field.name;
    const objProp = field.options?.json_name ?? field.name;
    return {
        messageProp,
        objProp
    }
};

export const fromJSON = {

    // sender: isSet(object.sender) ? String(object.sender) : ""
    string(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        args.context.addUtil('isSet');

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.identifier('String'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                getDefaultTSTypeFromProtoType(args.context, args.field, args.isOneOf)
            )
        )
    },

    number(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        args.context.addUtil('isSet');

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.identifier('Number'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                getDefaultTSTypeFromProtoType(args.context, args.field, args.isOneOf)
            )
        )
    },

    double(args: FromJSONMethod) {
        return fromJSON.number(args);
    },
    float(args: FromJSONMethod) {
        return fromJSON.number(args);
    },
    int32(args: FromJSONMethod) {
        return fromJSON.number(args);
    },
    sint32(args: FromJSONMethod) {
        return fromJSON.number(args);
    },
    uint32(args: FromJSONMethod) {
        return fromJSON.number(args);
    },
    fixed32(args: FromJSONMethod) {
        return fromJSON.number(args);
    },
    sfixed32(args: FromJSONMethod) {
        return fromJSON.number(args);
    },

    // disableMacros: isSet(object.disableMacros) ? Boolean(object.disableMacros) : false
    bool(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        args.context.addUtil('isSet');

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.identifier('Boolean'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                getDefaultTSTypeFromProtoType(args.context, args.field, args.isOneOf)
            )
        )
    },

    // int64Value: isSet(object.int64Value) ? Long.fromValue(object.int64Value) : Long.UZERO,
    long(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        args.context.addUtil('isSet');
        args.context.addUtil('Long');

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.memberExpression(
                        t.identifier('Long'),
                        t.identifier('fromValue')
                    ),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                getDefaultTSTypeFromProtoType(args.context, args.field, args.isOneOf)
            )
        );
    },

    int64(args: FromJSONMethod) {
        return fromJSON.long(args);
    },

    // uint64Value: isSet(object.uint64Value) ? Long.fromString(object.uint64Value) : Long.ZERO,
    uint64(args: FromJSONMethod) {
        return fromJSON.long(args);
    },

    sint64(args: FromJSONMethod) {
        return fromJSON.long(args);
    },
    fixed64(args: FromJSONMethod) {
        return fromJSON.long(args);
    },
    sfixed64(args: FromJSONMethod) {
        return fromJSON.long(args);
    },

    // signDoc: isSet(object.signDoc) ? SignDocDirectAux.fromJSON(object.signDoc) : undefined,
    type(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        const name = args.context.getTypeName(args.field);
        args.context.addUtil('isSet');

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.memberExpression(
                        t.identifier(name),
                        t.identifier('fromJSON')
                    ),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.identifier('undefined')
            )
        );
    },

    // mode: isSet(object.mode) ? signModeFromJSON(object.mode) : 0,
    enum(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        args.context.addUtil('isSet');
        const fromJSONFuncName = args.context.getFromEnum(args.field);

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.identifier(fromJSONFuncName),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                args.isOptional ? t.identifier('undefined') : t.numericLiteral(0)
            )
        );
    },

    // queryData: isSet(object.queryData) ? bytesFromBase64(object.queryData) : new Uint8Array()
    bytes(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        args.context.addUtil('isSet');
        args.context.addUtil('bytesFromBase64');

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.identifier('bytesFromBase64'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                getDefaultTSTypeFromProtoType(args.context, args.field, args.isOneOf)
            )
        );
    },


    // period: isSet(object.period) ? Duration.fromJSON(object.period) : undefined,

    duration(args: FromJSONMethod) {
        const durationFormat = args.context.pluginValue('prototypes.typingsFormat.duration');
        switch (durationFormat) {
            case 'string':
                return fromJSON.durationString(args);
            case 'duration':
            default:
                return fromJSON.type(args);
        }
    },

    // period: isSet(object.period) ? String(object.period) : undefined,

    durationString(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        args.context.addUtil('isSet');
        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.identifier('String'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.identifier('undefined')
            )
        );
    },

    // periodReset: isSet(object.periodReset) ? fromJsonTimestamp(object.periodReset) : undefined

    timestamp(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        args.context.addUtil('isSet');
        args.context.addUtil('fromJsonTimestamp');

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isSet'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.callExpression(
                    t.identifier('fromJsonTimestamp'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                t.identifier('undefined')
            )
        );
    },

    //  labels: isObject(object.labels) ? Object.entries(object.labels).reduce<{
    //     [key: string]: string;
    //   }>((acc, [key, value]) => {
    //     acc[key] = String(value);
    //     return acc;
    //   }, {}) : {},

    //   referenceMap: isObject(object.referenceMap) ? Object.entries(object.referenceMap).reduce<{
    //     [key: Long]: Reference;
    //   }>((acc, [key, value]) => {
    //     acc[Number(key)] = Reference.fromJSON(value);
    //     return acc;
    //   }, {}) : {},


    keyHash(args: FromJSONMethod) {
        const { messageProp, objProp } = getPropNames(args.field);
        const keyType = args.field.keyType;
        const valueType = args.field.parsedType.name;

        args.context.addUtil('isObject');

        let fromJSON = null;
        // valueTypeType: string for identifier
        let valueTypeType = valueType;
        switch (valueType) {
            case 'string':
                fromJSON = t.callExpression(
                    t.identifier('String'),
                    [
                        t.identifier('value')
                    ]
                )

                break;
            case 'int32':
            case 'uint32':
                valueTypeType = 'number';
                fromJSON = t.callExpression(
                    t.identifier('Number'),
                    [
                        t.identifier('value')
                    ]
                );

                break;
            case 'int64':
            case 'uint64':
                valueTypeType = 'Long';
                fromJSON = t.callExpression(
                    t.memberExpression(
                        t.identifier('Long'),
                        t.identifier('fromValue')
                    ),
                    [
                        t.tsAsExpression(
                            t.identifier('value'),
                            t.tsUnionType(
                                [
                                    t.tsTypeReference(
                                        t.identifier('Long')
                                    ),
                                    t.tsStringKeyword()
                                ]
                            )
                        )
                    ]
                )
                break;
            default:
                fromJSON = t.callExpression(
                    t.memberExpression(
                        t.identifier(valueType),
                        t.identifier('fromJSON')
                    ),
                    [
                        t.identifier('value')
                    ]
                );
        }

        let wrapKey = null;
        let keyTypeType = null;
        switch (keyType) {
            case 'string':
                wrapKey = (a) => a;
                keyTypeType = t.tsStringKeyword();
                break;
            case 'int64':
            case 'uint64':
                wrapKey = (a) => t.callExpression(
                    t.identifier('Number'),
                    [
                        a
                    ]
                );
                keyTypeType = t.tsTypeReference(t.identifier('Long'));
                break;
            case 'uint32':
            case 'int32':
                wrapKey = (a) => t.callExpression(
                    t.identifier('Number'),
                    [
                        a
                    ]
                );
                keyTypeType = t.tsTypeReference(t.identifier('number'));
                break;
            default:
                throw new Error('keyHash requires new type. Ask maintainers.');
        }

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.identifier('isObject'),
                    [
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        )
                    ]
                ),
                callExpression(
                    t.memberExpression(
                        t.callExpression(
                            t.memberExpression(
                                t.identifier('Object'),
                                t.identifier('entries')
                            ),
                            [
                                t.memberExpression(
                                    t.identifier('object'),
                                    t.identifier(objProp)
                                )
                            ]
                        ),
                        t.identifier('reduce')
                    ),
                    [
                        t.arrowFunctionExpression(
                            [
                                t.identifier('acc'),
                                t.arrayPattern(
                                    [
                                        t.identifier('key'),
                                        t.identifier('value')
                                    ]
                                )
                            ],
                            t.blockStatement([
                                t.expressionStatement(
                                    t.assignmentExpression(
                                        '=',
                                        t.memberExpression(
                                            t.identifier('acc'),
                                            wrapKey(t.identifier('key')),
                                            true
                                        ),
                                        fromJSON
                                    )
                                ),
                                t.returnStatement(
                                    t.identifier('acc')
                                )
                            ])
                        ),
                        t.objectExpression(
                            []
                        )
                    ],
                    t.tsTypeParameterInstantiation(
                        [
                            t.tsTypeLiteral(
                                [
                                    t.tsIndexSignature(
                                        [
                                            identifier('key', t.tsTypeAnnotation(
                                                keyTypeType
                                            ))
                                        ],
                                        t.tsTypeAnnotation(
                                            t.tsTypeReference(
                                                t.identifier(valueTypeType)
                                            )
                                        )
                                    )
                                ]
                            )
                        ]
                    )
                ),
                t.objectExpression([])
            )
        )
    },

    // codeIds: Array.isArray(object?.codeIds) ? object.codeIds.map((e: any) => Long.fromString(e)) : [],
    array(args: FromJSONMethod, expr: t.Expression) {
        const { messageProp, objProp } = getPropNames(args.field);

        return t.objectProperty(
            t.identifier(messageProp),
            t.conditionalExpression(
                t.callExpression(
                    t.memberExpression(
                        t.identifier('Array'),
                        t.identifier('isArray')
                    ),
                    [
                        t.optionalMemberExpression(
                            t.identifier('object'),
                            t.identifier(objProp),
                            false,
                            true
                        )
                    ]
                ),
                t.callExpression(
                    t.memberExpression(
                        t.memberExpression(
                            t.identifier('object'),
                            t.identifier(objProp)
                        ),
                        t.identifier('map')
                    ),
                    [
                        t.arrowFunctionExpression(
                            [
                                identifier('e', t.tsTypeAnnotation(
                                    t.tsAnyKeyword()
                                ))
                            ],
                            expr,
                            false
                        )
                    ]
                ),
                t.arrayExpression([])
            )
        )
    }
};

export const arrayTypes = {
    string() {
        return t.callExpression(
            t.identifier('String'),
            [
                t.identifier('e')
            ]
        );
    },

    bool() {
        return t.callExpression(
            t.identifier('Boolean'),
            [
                t.identifier('e')
            ]
        );
    },

    // myBytesArray: Array.isArray(object?.myBytesArray) ? object.myBytesArray.map((e: any) => bytesFromBase64(e)) : [],
    bytes(args: FromJSONMethod) {
        args.context.addUtil('bytesFromBase64');
        return t.callExpression(
            t.identifier('bytesFromBase64'),
            [
                t.identifier('e')
            ]
        );
    },
    // codeIds: Array.isArray(object?.codeIds) ? object.codeIds.map((e: any) => Long.fromValue(e)) : [],
    long() {
        return t.callExpression(
            t.memberExpression(
                t.identifier('Long'),
                t.identifier('fromValue')
            ),
            [
                t.identifier('e')
            ]
        );
    },
    uint64() {
        return arrayTypes.long();
    },
    int64() {
        return arrayTypes.long();
    },
    sint64() {
        return arrayTypes.long();
    },
    fixed64() {
        return arrayTypes.long();
    },
    sfixed64() {
        return arrayTypes.long();
    },
    // myUint32Array: Array.isArray(object?.myUint32Array) ? object.myUint32Array.map((e: any) => Number(e)) : [],
    number() {
        return t.callExpression(
            t.identifier('Number'),
            [
                t.identifier('e')
            ]
        );
    },

    // myDoubleArray: Array.isArray(object?.myDoubleArray) ? object.myDoubleArray.map((e: any) => Number(e)) : [],
    uint32() {
        return arrayTypes.number();
    },
    int32() {
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
    double() {
        return arrayTypes.number();
    },
    float() {
        return arrayTypes.number();
    },

    // arrayField: Array.isArray(object?.arrayField) ? object.arrayField.map((e: any) => scalarTypeFromJSON(e)) : []
    enum(args: FromJSONMethod) {
        const fromJSONFuncName = args.context.getFromEnum(args.field);
        return t.callExpression(
            t.identifier(fromJSONFuncName),
            [
                t.identifier('e')
            ]
        );
    },

    // tokenInMaxs: Array.isArray(object?.tokenInMaxs) ? object.tokenInMaxs.map((e: any) => Coin.fromJSON(e)) : []
    type(args: FromJSONMethod) {
        const name = args.context.getTypeName(args.field);
        return t.callExpression(
            t.memberExpression(
                t.identifier(name),
                t.identifier('fromJSON')
            ),
            [
                t.identifier('e')
            ]
        );
    }
};

