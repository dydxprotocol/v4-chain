import * as t from '@babel/types';
import { getFieldOptionality, getOneOfs } from '..';
import { identifier, objectMethod } from '../../../utils';
import { ProtoParseContext } from '../../context';
import { getBaseCreateTypeFuncName } from '../types';
import { ProtoType, ProtoField } from '@osmonauts/types';
import { baseTypes, decode } from './utils';

const needsImplementation = (name: string, field: ProtoField) => {
    throw new Error(`need to implement decode (${field.type} rules[${field.rule}] name[${name}])`);
}
export interface DecodeMethod {
    typeName: string;
    context: ProtoParseContext;
    field: ProtoField;
    isOptional: boolean;
}

export const decodeMethodFields = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const oneOfs = getOneOfs(proto);
    return Object.keys(proto.fields ?? {}).map(fieldName => {
        const field = {
            name: fieldName,
            ...proto.fields[fieldName]
        };

        const isOneOf = oneOfs.includes(fieldName);
        const isOptional = getFieldOptionality(context, field, isOneOf);

        const args: DecodeMethod = {
            typeName: name,
            context,
            field,
            isOptional
        };

        if (field.rule === 'repeated') {
            switch (field.type) {
                case 'string':
                    return decode.array(args, baseTypes.string(args));
                case 'bytes':
                    return decode.array(args, baseTypes.bytes(args));
                case 'double':
                    return decode.tagDelimArray(args, baseTypes.double(args));
                case 'bool':
                    return decode.tagDelimArray(args, baseTypes.bool(args));
                case 'float':
                    return decode.tagDelimArray(args, baseTypes.float(args));
                case 'int32':
                    return decode.tagDelimArray(args, baseTypes.int32(args));
                case 'sint32':
                    return decode.tagDelimArray(args, baseTypes.sint32(args));
                case 'uint32':
                    return decode.tagDelimArray(args, baseTypes.uint32(args));
                case 'fixed32':
                    return decode.tagDelimArray(args, baseTypes.fixed32(args));
                case 'sfixed32':
                    return decode.tagDelimArray(args, baseTypes.sfixed32(args));
                case 'int64':
                    return decode.tagDelimArray(args, baseTypes.int64(args));
                case 'sint64':
                    return decode.tagDelimArray(args, baseTypes.sint64(args));
                case 'uint64':
                    return decode.tagDelimArray(args, baseTypes.uint64(args));
                case 'fixed64':
                    return decode.tagDelimArray(args, baseTypes.fixed64(args));
                case 'sfixed64':
                    return decode.tagDelimArray(args, baseTypes.sfixed64(args));
                default:
                    switch (field.parsedType.type) {
                        case 'Enum':
                            return decode.tagDelimArray(args, baseTypes.enum(args));
                        case 'Type':
                            return decode.typeArray(args);
                    }
                    return needsImplementation(fieldName, field);
            }

        }


        if (field.keyType) {
            // currently they all look the same for decode()
            return decode.keyHash(args);
        }

        switch (field.type) {
            case 'string':
                return decode.string(args);
            case 'int32':
                return decode.int32(args);
            case 'sint32':
                return decode.sint32(args);
            case 'uint32':
                return decode.uint32(args);
            case 'fixed32':
                return decode.fixed32(args);
            case 'sfixed32':
                return decode.sfixed32(args);
            case 'int64':
                return decode.int64(args);
            case 'sint64':
                return decode.sint64(args);
            case 'uint64':
                return decode.uint64(args);
            case 'fixed64':
                return decode.fixed64(args);
            case 'sfixed64':
                return decode.sfixed64(args);
            case 'double':
                return decode.double(args);
            case 'float':
                return decode.float(args);
            case 'bytes':
                return decode.bytes(args);
            case 'bool':
                return decode.bool(args);
            case 'google.protobuf.Duration':
            case 'Duration':
                return decode.duration(args);
            case 'google.protobuf.Timestamp':
            case 'Timestamp':
                return decode.timestamp(args);
            default:
                switch (field.parsedType.type) {
                    case 'Enum':
                        return decode.enum(args);
                    case 'Type':
                        return decode.type(args);
                }
                return needsImplementation(fieldName, field);
        }
    });
};

export const decodeMethod = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    context.addUtil('_m0');

    let returnType = name;
    // decode can be coupled to API requests
    if (context.store.responses[name]) {
        // returnType = name + 'SDKType';
        returnType = name;
    }

    return objectMethod(
        'method',
        t.identifier('decode'),
        [
            identifier('input',
                t.tsTypeAnnotation(
                    t.tsUnionType(
                        [
                            t.tsTypeReference(
                                t.tsQualifiedName(
                                    t.identifier('_m0'),
                                    t.identifier('Reader')
                                ),
                                null
                            ),
                            t.tsTypeReference(
                                t.identifier('Uint8Array')
                            )
                        ]
                    )
                ),
                false
            ),
            identifier('length', t.tsTypeAnnotation(
                t.tsNumberKeyword()
            ), true)
        ],
        t.blockStatement([

            /*
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
            */
            t.variableDeclaration(
                'const',
                [
                    t.variableDeclarator(
                        t.identifier('reader'),
                        t.conditionalExpression(
                            t.binaryExpression(
                                'instanceof',
                                t.identifier('input'),
                                t.memberExpression(
                                    t.identifier('_m0'),
                                    t.identifier('Reader')
                                )
                            ),
                            t.identifier('input'),
                            t.newExpression(
                                t.memberExpression(
                                    t.identifier('_m0'),
                                    t.identifier('Reader')
                                ),
                                [
                                    t.identifier('input')
                                ]
                            )
                        )
                    )
                ]
            ),

            /*
    let end = length === undefined ? reader.len : reader.pos + length;
            */

            t.variableDeclaration(
                'let',
                [
                    t.variableDeclarator(
                        t.identifier('end'),
                        t.conditionalExpression(
                            t.binaryExpression('===',
                                t.identifier('length'),
                                t.identifier('undefined')
                            ),
                            t.memberExpression(
                                t.identifier('reader'),
                                t.identifier('len')
                            ),
                            t.binaryExpression(
                                '+',
                                t.memberExpression(
                                    t.identifier('reader'),
                                    t.identifier('pos')
                                ),
                                t.identifier('length')
                            )
                        )
                    )
                ]
            ),

            /*
            
    const message = createBaseMsgJoinPool();

            */

            t.variableDeclaration(
                'const',
                [
                    t.variableDeclarator(
                        t.identifier('message'),
                        t.callExpression(

                            // 
                            t.identifier(getBaseCreateTypeFuncName(name)),
                            []
                        )
                    )
                ]
            ),

            ///////////
            ///////////
            ///////////

            t.whileStatement(
                t.binaryExpression(
                    '<',
                    t.memberExpression(
                        t.identifier('reader'),
                        t.identifier('pos')
                    ),
                    t.identifier('end')
                ),
                t.blockStatement([

                    /// DECODE BODY
                    t.variableDeclaration(
                        'const',
                        [
                            t.variableDeclarator(
                                t.identifier('tag'),
                                t.callExpression(
                                    t.memberExpression(
                                        t.identifier('reader'),
                                        t.identifier('uint32')
                                    ),
                                    []
                                )
                            )
                        ]
                    ),


                    t.switchStatement(
                        t.binaryExpression(
                            '>>>',
                            t.identifier('tag'),
                            t.numericLiteral(3)
                        ),
                        [


                            ...decodeMethodFields(context, name, proto),

                            /*
                            default:
                                    reader.skipType(tag & 7);
                                    break;
                            */
                            t.switchCase(
                                null,
                                [
                                    t.expressionStatement(
                                        t.callExpression(
                                            t.memberExpression(
                                                t.identifier('reader'),
                                                t.identifier('skipType')
                                            ),
                                            [
                                                t.binaryExpression(
                                                    '&',
                                                    t.identifier('tag'),
                                                    t.numericLiteral(7)
                                                )
                                            ]
                                        )
                                    ),
                                    t.breakStatement()
                                ]
                            )
                        ]
                    )

                ])
            ),

            // RETURN STATEMENT
            t.returnStatement(
                t.identifier('message')
            )

        ]),
        false,
        false,
        false,
        t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(returnType)
            )
        )
    )
};

