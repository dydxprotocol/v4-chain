import * as t from '@babel/types';
import { ProtoType, ProtoField } from '@osmonauts/types';
import { pascal } from 'case';
import { getFieldOptionality, getFieldOptionalityForDefaults, getOneOfs } from '..';
import { identifier, objectMethod } from '../../../utils';
import { ProtoParseContext } from '../../context';
import { fromSDK, arrayTypes } from './utils';

const needsImplementation = (name: string, field: ProtoField) => {
    throw new Error(`need to implement fromSDK (${field.type} rules[${field.rule}] name[${name}])`);
}
export interface FromSDKMethod {
    context: ProtoParseContext;
    field: ProtoField;
    isOptional: boolean;
}

export const fromSDKMethodFields = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const oneOfs = getOneOfs(proto);
    const fields = Object.keys(proto.fields ?? {}).map(fieldName => {
        const field = {
            name: fieldName,
            ...proto.fields[fieldName]
        };

        const isOneOf = oneOfs.includes(fieldName);
        const isOptional = getFieldOptionalityForDefaults(context, field, isOneOf);

        const args: FromSDKMethod = {
            context,
            field,
            isOptional
        };

        if (field.rule === 'repeated') {
            switch (field.type) {
                case 'string':
                    return fromSDK.array(args, arrayTypes.string());
                case 'bytes':
                    return fromSDK.array(args, arrayTypes.bytes(args));
                case 'bool':
                    return fromSDK.array(args, arrayTypes.bool());
                case 'float':
                    return fromSDK.array(args, arrayTypes.float());
                case 'double':
                    return fromSDK.array(args, arrayTypes.double());
                case 'int32':
                    return fromSDK.array(args, arrayTypes.int32());
                case 'sint32':
                    return fromSDK.array(args, arrayTypes.sint32());
                case 'uint32':
                    return fromSDK.array(args, arrayTypes.uint32());
                case 'fixed32':
                    return fromSDK.array(args, arrayTypes.fixed32());
                case 'sfixed32':
                    return fromSDK.array(args, arrayTypes.sfixed32());
                case 'int64':
                    return fromSDK.array(args, arrayTypes.int64());
                case 'sint64':
                    return fromSDK.array(args, arrayTypes.sint64());
                case 'uint64':
                    return fromSDK.array(args, arrayTypes.uint64());
                case 'fixed64':
                    return fromSDK.array(args, arrayTypes.fixed64());
                case 'sfixed64':
                    return fromSDK.array(args, arrayTypes.sfixed64());
                default:
                    switch (field.parsedType.type) {
                        case 'Enum':
                            return fromSDK.array(args, arrayTypes.enum(args));
                        case 'Type':
                            return fromSDK.array(args, arrayTypes.type(args));
                    }
                    return needsImplementation(fieldName, field);
            }
        }

        if (field.keyType) {
            switch (field.keyType) {
                case 'string':
                case 'int64':
                case 'sint64':
                case 'uint64':
                case 'fixed64':
                case 'sfixed64':
                case 'int32':
                case 'sint32':
                case 'uint32':
                case 'fixed32':
                case 'sfixed32':
                    return fromSDK.keyHash(args);
                default:
                    return needsImplementation(fieldName, field);
            }
        }

        switch (field.type) {
            case 'string':
                return fromSDK.string(args);
            case 'bytes':
                return fromSDK.bytes(args);
            case 'bool':
                return fromSDK.bool(args);
            case 'double':
                return fromSDK.double(args);
            case 'float':
                return fromSDK.float(args);
            case 'int32':
                return fromSDK.int32(args);
            case 'sint32':
                return fromSDK.sint32(args);
            case 'uint32':
                return fromSDK.uint32(args);
            case 'fixed32':
                return fromSDK.fixed32(args);
            case 'sfixed32':
                return fromSDK.sfixed32(args);
            case 'int64':
                return fromSDK.int64(args);
            case 'sint64':
                return fromSDK.sint64(args);
            case 'uint64':
                return fromSDK.uint64(args);
            case 'fixed64':
                return fromSDK.fixed64(args);
            case 'sfixed64':
                return fromSDK.sfixed64(args);
            case 'Duration':
            case 'google.protobuf.Duration':
                return fromSDK.duration(args);
            case 'Timestamp':
            case 'google.protobuf.Timestamp':
                return fromSDK.timestamp(args);
            default:
                switch (field.parsedType.type) {
                    case 'Enum':
                        return fromSDK.enum(args);
                    case 'Type':
                        return fromSDK.type(args);
                }
                return needsImplementation(fieldName, field);
        }
    });
    return fields;
};


export const fromSDKMethod = (context: ProtoParseContext, name: string, proto: ProtoType) => {
    const fields = fromSDKMethodFields(context, name, proto);
    let varName = 'object';
    if (!fields.length) {
        varName = '_';
    }

    const SDKTypeName =
        [name, 'SDKType']
            .filter(Boolean).join('');

    return objectMethod('method',
        t.identifier('fromSDK'),
        [
            identifier(varName,
                t.tsTypeAnnotation(
                    t.tsTypeReference(
                        t.identifier(SDKTypeName)
                    )
                ),
                false
            )

        ],
        t.blockStatement(
            [
                t.returnStatement(
                    t.objectExpression(fields)
                )
            ]
        ),
        false,
        false,
        false,
        t.tsTypeAnnotation(
            t.tsTypeReference(
                t.identifier(name)
            )
        )
    )
};
