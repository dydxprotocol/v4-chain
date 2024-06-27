import * as t from '@babel/types';
import { getFieldOptionality, getFieldOptionalityForDefaults, getOneOfs } from '..';
import { identifier, objectMethod } from '../../../utils';
import { ProtoParseContext } from '../../context';
import { ProtoField, ProtoType } from '@osmonauts/types';
import { arrayTypes, toJSON } from './utils';

const needsImplementation = (name: string, field: ProtoField) => {
    throw new Error(`need to implement toJSON (${field.type} rules[${field.rule}] name[${name}])`);
}

export interface ToJSONMethod {
    context: ProtoParseContext;
    field: ProtoField;
    isOneOf: boolean;
    isOptional: boolean;
}

export const toJSONMethodFields = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const oneOfs = getOneOfs(proto);
    const fields = Object.keys(proto.fields ?? {}).reduce((m, fieldName) => {
        const field = {
            name: fieldName,
            ...proto.fields[fieldName]
        };

        const isOneOf = oneOfs.includes(fieldName);
        const isOptional = getFieldOptionalityForDefaults(context, field, isOneOf);

        const args: ToJSONMethod = {
            context,
            field,
            isOneOf,
            isOptional
        };

        if (field.rule === 'repeated') {
            switch (field.type) {
                case 'string':
                    return [...m, toJSON.array(args, arrayTypes.string())];
                case 'bytes':
                    return [...m, toJSON.array(args, arrayTypes.bytes(args))];
                case 'bool':
                    return [...m, toJSON.array(args, arrayTypes.bool())];
                case 'double':
                    return [...m, toJSON.array(args, arrayTypes.double())];
                case 'float':
                    return [...m, toJSON.array(args, arrayTypes.float())];
                case 'int32':
                    return [...m, toJSON.array(args, arrayTypes.int32())];
                case 'sint32':
                    return [...m, toJSON.array(args, arrayTypes.sint32())];
                case 'uint32':
                    return [...m, toJSON.array(args, arrayTypes.uint32())];
                case 'fixed32':
                    return [...m, toJSON.array(args, arrayTypes.fixed32())];
                case 'sfixed32':
                    return [...m, toJSON.array(args, arrayTypes.sfixed32())];
                case 'int64':
                    return [...m, toJSON.array(args, arrayTypes.int64(args))];
                case 'sint64':
                    return [...m, toJSON.array(args, arrayTypes.sint64(args))];
                case 'uint64':
                    return [...m, toJSON.array(args, arrayTypes.uint64(args))];
                case 'fixed64':
                    return [...m, toJSON.array(args, arrayTypes.fixed64(args))];
                case 'sfixed64':
                    return [...m, toJSON.array(args, arrayTypes.sfixed64(args))];
                default:
                    switch (field.parsedType.type) {
                        case 'Enum':
                            return [...m, toJSON.array(args, arrayTypes.enum(args))];
                        case 'Type':
                            return [...m, toJSON.array(args, arrayTypes.type(args))];
                    }
                    return needsImplementation(fieldName, field);
            }

        }

        if (field.keyType) {
            switch (field.keyType) {
                case 'string':
                case 'int32':
                case 'sint32':
                case 'uint32':
                case 'fixed32':
                case 'sfixed32':
                case 'int64':
                case 'sint64':
                case 'uint64':
                case 'fixed64':
                case 'sfixed64':
                    return [...m, ...toJSON.keyHash(args)];
                default:
                    return needsImplementation(fieldName, field);
            }
        }

        switch (field.type) {
            case 'string':
                return [...m, toJSON.string(args)];
            case 'double':
                return [...m, toJSON.double(args)];
            case 'float':
                return [...m, toJSON.float(args)];
            case 'bytes':
                return [...m, toJSON.bytes(args)];
            case 'bool':
                return [...m, toJSON.bool(args)];
            case 'int32':
                return [...m, toJSON.int32(args)];
            case 'sint32':
                return [...m, toJSON.sint32(args)];
            case 'uint32':
                return [...m, toJSON.uint32(args)];
            case 'fixed32':
                return [...m, toJSON.fixed32(args)];
            case 'sfixed32':
                return [...m, toJSON.sfixed32(args)];
            case 'int64':
                return [...m, toJSON.int64(args)];
            case 'sint64':
                return [...m, toJSON.sint64(args)];
            case 'uint64':
                return [...m, toJSON.uint64(args)];
            case 'fixed64':
                return [...m, toJSON.fixed64(args)];
            case 'sfixed64':
                return [...m, toJSON.sfixed64(args)];
            case 'google.protobuf.Duration':
            case 'Duration':
                return [...m, toJSON.duration(args)];
            case 'google.protobuf.Timestamp':
            case 'Timestamp':
                return [...m, toJSON.timestamp(args)];
            default:
                switch (field.parsedType.type) {
                    case 'Enum':
                        return [...m, toJSON.enum(args)];
                    case 'Type':
                        return [...m, toJSON.type(args)];
                }
                return needsImplementation(fieldName, field);
        }
    }, []);
    return fields;
};

export const toJSONMethod = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const fields = toJSONMethodFields(context, name, proto);
    let varName = 'message';
    if (!fields.length) {
        varName = '_';
    }
    return objectMethod('method',
        t.identifier('toJSON'),
        [
            identifier(
                varName,
                t.tsTypeAnnotation(
                    t.tsTypeReference(
                        t.identifier(name)
                    )
                )
            )
        ],
        t.blockStatement([
            t.variableDeclaration(
                'const',
                [
                    t.variableDeclarator(
                        identifier('obj', t.tsTypeAnnotation(t.tsAnyKeyword())),
                        t.objectExpression([])
                    )
                ]
            ),

            ...fields,

            // RETURN 
            t.returnStatement(t.identifier('obj'))

        ]),
        false,
        false,
        false,
        t.tsTypeAnnotation(
            t.tsUnknownKeyword()
        )
    );
};
