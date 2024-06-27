import * as t from '@babel/types';
import { ProtoType, ProtoField } from '@osmonauts/types';
import { getFieldOptionality, getFieldOptionalityForDefaults, getOneOfs } from '..';
import { identifier, objectMethod } from '../../../utils';
import { ProtoParseContext } from '../../context';
import { encode, arrayTypes } from './utils';

const needsImplementation = (name: string, field: ProtoField) => {
    throw new Error(`need to implement encode (${field.type} rules[${field.rule}] name[${name}])`);
}
export interface EncodeMethod {
    typeName: string;
    context: ProtoParseContext;
    field: ProtoField;
    isOptional: boolean;
    isOneOf: boolean;
}

export const encodeMethodFields = (context: ProtoParseContext, name: string, proto: ProtoType) => {

    const oneOfs = getOneOfs(proto);
    return Object.keys(proto.fields ?? {}).reduce((m, fieldName) => {


        const field = {
            name: fieldName,
            ...proto.fields[fieldName]
        };

        const isOneOf = oneOfs.includes(fieldName);
        const isOptional = getFieldOptionalityForDefaults(context, field, isOneOf);

        const args: EncodeMethod = {
            typeName: name,
            context,
            field,
            isOneOf,
            isOptional
        };

        if (field.rule === 'repeated') {
            switch (field.type) {
                case 'string':
                    return [...m, ...encode.array(args, arrayTypes.string(args))];
                case 'bytes':
                    return [...m, ...encode.array(args, arrayTypes.bytes(args))];
                case 'bool':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.bool())];
                case 'double':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.double())];
                case 'float':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.float())];
                case 'int32':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.int32())];
                case 'sint32':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.sint32())];
                case 'uint32':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.uint32())];
                case 'fixed32':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.fixed32())];
                case 'sfixed32':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.sfixed32())];
                case 'int64':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.int64())];
                case 'sint64':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.sint64())];
                case 'uint64':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.uint64())];
                case 'fixed64':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.fixed64())];
                case 'sfixed64':
                    return [...m, ...encode.forkDelimArray(args, arrayTypes.sfixed64())];
                default:
                    switch (field.parsedType.type) {
                        case 'Enum':
                            return [...m, ...encode.forkDelimArray(args, arrayTypes.enum())];
                        case 'Type':
                            return [...m, ...encode.typeArray(args)];
                    }
                    return needsImplementation(fieldName, field);
            }
        }

        if (field.keyType) {
            // currently they all look the same for encode()
            return [...m, encode.keyHash(args)];
        }

        switch (field.type) {
            case 'string':
                return [...m, encode.string(args)];
            case 'int32':
                return [...m, encode.int32(args)];
            case 'sint32':
                return [...m, encode.sint32(args)];
            case 'uint32':
                return [...m, encode.uint32(args)];
            case 'fixed32':
                return [...m, encode.fixed32(args)];
            case 'sfixed32':
                return [...m, encode.sfixed32(args)];
            case 'int64':
                return [...m, encode.int64(args)];
            case 'sint64':
                return [...m, encode.sint64(args)];
            case 'uint64':
                return [...m, encode.uint64(args)];
            case 'fixed64':
                return [...m, encode.fixed64(args)];
            case 'sfixed64':
                return [...m, encode.sfixed64(args)];
            case 'double':
                return [...m, encode.double(args)];
            case 'float':
                return [...m, encode.float(args)];
            case 'bool':
                return [...m, encode.bool(args)];
            case 'bytes':
                return [...m, encode.bytes(args)];
            case 'Duration':
            case 'google.protobuf.Duration':
                return [...m, encode.duration(args)];
            case 'Timestamp':
            case 'google.protobuf.Timestamp':
                return [...m, encode.timestamp(args)];
            default:
                switch (field.parsedType.type) {
                    case 'Enum':
                        return [...m, encode.enum(args)];
                    case 'Type':
                        return [...m, encode.type(args)];
                }
                return needsImplementation(fieldName, field);
        }
    }, []);
};

export const encodeMethod = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    context.addUtil('_m0');

    const fields = encodeMethodFields(context, name, proto);
    let varName = 'message';
    if (!fields.length) {
        varName = '_';
    }
    const body = [
        ...fields,
        /* RETURN writer */
        t.returnStatement(
            t.identifier('writer')
        )
    ];

    try {
        t.blockStatement(body);
    } catch (e) {
        console.log(body);
        throw e;
    }

    return objectMethod(
        'method',
        t.identifier('encode'),
        [
            // args

            identifier(varName, t.tsTypeAnnotation(
                t.tsTypeReference(
                    t.identifier(name)
                )
            ), false),

            t.assignmentPattern(
                identifier('writer', t.tsTypeAnnotation(
                    t.tsTypeReference(
                        t.tsQualifiedName(
                            t.identifier('_m0'),
                            t.identifier('Writer')
                        )
                    )
                )),
                t.callExpression(
                    t.memberExpression(
                        t.memberExpression(
                            t.identifier('_m0'),
                            t.identifier('Writer')
                        ),
                        t.identifier('create')
                    ),
                    []
                )
            )
        ],

        // body 
        t.blockStatement(body),
        false,
        false,
        false,
        // return type
        t.tsTypeAnnotation(
            t.tsTypeReference(
                t.tsQualifiedName(
                    t.identifier('_m0'),
                    t.identifier('Writer')
                )
            )
        )
    )
};
